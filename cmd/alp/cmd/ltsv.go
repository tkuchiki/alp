package cmd

import (
	"os"

	"github.com/tkuchiki/alp/counter"

	"github.com/tkuchiki/alp/log_reader"

	"github.com/spf13/cobra"
	"github.com/tkuchiki/alp/options"
	"github.com/tkuchiki/alp/parsers"
	"github.com/tkuchiki/alp/profiler"
)

func newLTSVCmd(flags *flags) *cobra.Command {
	var ltsvCmd = &cobra.Command{
		Use:   "ltsv",
		Short: "Profile the log_reader for LTSV",
		Long:  `Profile the log_reader for LTSV`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := flags.createLTSVOptions(cmd)
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

			parser := newLTSVParser(opts, f)

			err = prof.Run(flags.sortOptions, parser, nil)

			return err
		},
	}

	flags.defineProfileOptions(ltsvCmd)
	flags.defineLTSVOptions(ltsvCmd)

	ltsvCmd.Flags().SortFlags = false
	ltsvCmd.PersistentFlags().SortFlags = false
	ltsvCmd.InheritedFlags().SortFlags = false

	return ltsvCmd
}

func newLTSVParser(opts *options.Options, f *os.File) parsers.Parser {
	label := parsers.NewLTSVLabel(opts.LTSV.UriLabel, opts.LTSV.MethodLabel, opts.LTSV.TimeLabel,
		opts.LTSV.ApptimeLabel, opts.LTSV.ReqtimeLabel, opts.LTSV.SizeLabel, opts.LTSV.StatusLabel,
	)

	return parsers.NewLTSVParser(f, label, opts.QueryString, opts.QueryStringIgnoreValues)
}

func newLTSVDiffCmd(flags *flags) *cobra.Command {
	ltsvDiffCmd := newDiffSubCmd()
	ltsvDiffCmd.RunE = func(cmd *cobra.Command, args []string) error {
		opts, err := flags.createLTSVDiffOptions(cmd)
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

		fromParser := newLTSVParser(opts, fromf)

		toProf := profiler.NewProfiler(os.Stdout, os.Stderr, opts)
		toProf.DisableLoad()

		tof, err := toProf.Open(to)
		if err != nil {
			return err
		}
		defer tof.Close()

		toParser := newLTSVParser(opts, tof)

		return runDiff(flags.sortOptions,
			fromProf, fromParser,
			toProf, toParser,
		)
	}

	flags.defineDiffSubCommandOptions(ltsvDiffCmd)
	flags.defineLTSVOptions(ltsvDiffCmd)

	ltsvDiffCmd.Flags().SortFlags = false
	ltsvDiffCmd.PersistentFlags().SortFlags = false
	ltsvDiffCmd.InheritedFlags().SortFlags = false

	return ltsvDiffCmd
}

func newLTSVTopNCmd(flags *flags) *cobra.Command {
	ltsvTopNCmd := newTopNSubCmd()
	ltsvTopNCmd.RunE = func(cmd *cobra.Command, args []string) error {
		opts, err := flags.createLTSVTopNOptions(cmd)
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

		parser := newLTSVParser(opts, f)

		return runTopN(logReader, parser)
	}

	flags.defineTopNSubCommandOptions(ltsvTopNCmd)
	flags.defineLTSVOptions(ltsvTopNCmd)

	ltsvTopNCmd.Flags().SortFlags = false
	ltsvTopNCmd.PersistentFlags().SortFlags = false
	ltsvTopNCmd.InheritedFlags().SortFlags = false

	return ltsvTopNCmd
}

func newLTSVCountCmd(flags *flags) *cobra.Command {
	ltsvCountCmd := newCountSubCmd()
	ltsvCountCmd.RunE = func(cmd *cobra.Command, args []string) error {
		opts, err := flags.createLTSVCountOptions(cmd)
		if err != nil {
			return err
		}

		counter := counter.NewCounter(os.Stdout, os.Stderr, opts)

		f, err := counter.Open(opts.File)
		if err != nil {
			return err
		}
		defer f.Close()

		parser := newLTSVParser(opts, f)

		return runCount(counter, parser, opts)
	}

	flags.defineCountSubCommandOptions(ltsvCountCmd)
	flags.defineLTSVOptions(ltsvCountCmd)

	ltsvCountCmd.Flags().SortFlags = false
	ltsvCountCmd.PersistentFlags().SortFlags = false
	ltsvCountCmd.InheritedFlags().SortFlags = false

	return ltsvCountCmd
}
