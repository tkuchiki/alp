package cmd

import (
	"os"

	"github.com/tkuchiki/alp/log_reader"

	"github.com/spf13/cobra"
	"github.com/tkuchiki/alp/options"
	"github.com/tkuchiki/alp/parsers"
	"github.com/tkuchiki/alp/profiler"
)

func newPcapCmd(flags *flags) *cobra.Command {
	var pcapCmd = &cobra.Command{
		Use:   "pcap",
		Short: "Profile the HTTP requests for captured packets",
		Long:  `Profile the HTTP requests for captured packets`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := flags.createPcapOptions(cmd)
			if err != nil {
				return err
			}

			prof := profiler.NewProfiler(os.Stdout, os.Stderr, opts)

			if err = prof.ValidatePrinter(); err != nil {
				return err
			}

			f, err := prof.Open(opts.File)
			if err != nil {
				return err
			}
			defer f.Close()

			parser, err := newPcapParser(opts, f)
			if err != nil {
				return err
			}

			err = prof.Run(flags.sortOptions, parser, nil)

			return err
		},
	}

	flags.defineProfileOptions(pcapCmd)
	flags.definePcapOptions(pcapCmd)

	pcapCmd.Flags().SortFlags = false
	pcapCmd.PersistentFlags().SortFlags = false
	pcapCmd.InheritedFlags().SortFlags = false

	return pcapCmd
}

func newPcapParser(opts *options.Options, f *os.File) (parsers.Parser, error) {
	return parsers.NewPcapParser(f, opts.Pcap.ServerIPs, opts.Pcap.ServerPort, opts.QueryString, opts.QueryStringIgnoreValues)
}

func newPcapDiffCmd(flags *flags) *cobra.Command {
	pcapDiffCmd := newDiffSubCmd()
	pcapDiffCmd.RunE = func(cmd *cobra.Command, args []string) error {
		opts, err := flags.createPcapDiffOptions(cmd)
		if err != nil {
			return err
		}

		from, to := getFromTo(opts.Load, args)

		fromProf := profiler.NewProfiler(os.Stdout, os.Stderr, opts)

		if err = fromProf.ValidatePrinter(); err != nil {
			return err
		}

		fromf, err := fromProf.Open(from)
		if err != nil {
			return err
		}
		defer fromf.Close()

		var fromParser parsers.Parser
		if opts.Load == "" {
			fromParser, err = newPcapParser(opts, fromf)
			if err != nil {
				return err
			}
		}

		toProf := profiler.NewProfiler(os.Stdout, os.Stderr, opts)
		toProf.DisableLoad()

		tof, err := toProf.Open(to)
		if err != nil {
			return err
		}
		defer tof.Close()

		toParser, err := newPcapParser(opts, tof)
		if err != nil {
			return err
		}

		return runDiff(flags.sortOptions,
			fromProf, fromParser,
			toProf, toParser,
		)
	}

	flags.defineDiffSubCommandOptions(pcapDiffCmd)
	flags.definePcapOptions(pcapDiffCmd)

	pcapDiffCmd.Flags().SortFlags = false
	pcapDiffCmd.PersistentFlags().SortFlags = false
	pcapDiffCmd.InheritedFlags().SortFlags = false

	return pcapDiffCmd
}

func newPcapTopNCmd(flags *flags) *cobra.Command {
	pcapTopNCmd := newTopNSubCmd()
	pcapTopNCmd.RunE = func(cmd *cobra.Command, args []string) error {
		opts, err := flags.createPcapTopNOptions(cmd)
		if err != nil {
			return err
		}

		n, err := getN(args)
		if err != nil {
			return err
		}

		logReader := log_reader.NewAccessLogReader(os.Stdout, os.Stderr, opts, n)

		f, err := logReader.Open(opts.File)
		if err != nil {
			return err
		}
		defer f.Close()

		parser, err := newPcapParser(opts, f)
		if err != nil {
			return err
		}

		return runTopN(logReader, parser)
	}

	flags.defineTopNSubCommandOptions(pcapTopNCmd)
	flags.definePcapOptions(pcapTopNCmd)

	pcapTopNCmd.Flags().SortFlags = false
	pcapTopNCmd.PersistentFlags().SortFlags = false
	pcapTopNCmd.InheritedFlags().SortFlags = false

	return pcapTopNCmd
}
