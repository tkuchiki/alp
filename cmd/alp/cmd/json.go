package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tkuchiki/alp/options"
	"github.com/tkuchiki/alp/parsers"
	"github.com/tkuchiki/alp/profiler"
)

func newJSONCmd(flags *flags) *cobra.Command {
	var jsonCmd = &cobra.Command{
		Use:   "json",
		Short: "Profile the logs for JSON",
		Long:  `Profile the logs for JSON`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := flags.createJSONOptions(cmd)
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

			parser := newJsonParser(opts, f)

			err = prof.Run(flags.sortOptions, parser, nil)

			return err
		},
	}

	flags.defineProfileOptions(jsonCmd)
	flags.defineJSONOptions(jsonCmd)

	jsonCmd.Flags().SortFlags = false
	jsonCmd.PersistentFlags().SortFlags = false
	jsonCmd.InheritedFlags().SortFlags = false

	return jsonCmd
}

func newJsonParser(opts *options.Options, f *os.File) parsers.Parser {
	keys := parsers.NewJSONKeys(opts.JSON.UriKey, opts.JSON.MethodKey, opts.JSON.TimeKey,
		opts.JSON.ResponseTimeKey, opts.JSON.RequestTimeKey, opts.JSON.BodyBytesKey, opts.JSON.StatusKey)

	return parsers.NewJSONParser(f, keys, opts.QueryString, opts.QueryStringIgnoreValues)
}

func newJsonDiffCmd(flags *flags) *cobra.Command {
	jsonDiffCmd := newDiffSubCmd()
	jsonDiffCmd.RunE = func(cmd *cobra.Command, args []string) error {
		opts, err := flags.createJSONDiffOptions(cmd)
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

		fromParser := newJsonParser(opts, fromf)

		toProf := profiler.NewProfiler(os.Stdout, os.Stderr, opts)
		toProf.DisableLoad()

		tof, err := toProf.Open(to)
		if err != nil {
			return err
		}
		defer tof.Close()

		toParser := newJsonParser(opts, tof)

		return runDiff(flags.sortOptions,
			fromProf, fromParser,
			toProf, toParser,
		)
	}

	flags.defineDiffSubCommandOptions(jsonDiffCmd)
	flags.defineJSONOptions(jsonDiffCmd)

	jsonDiffCmd.Flags().SortFlags = false
	jsonDiffCmd.PersistentFlags().SortFlags = false
	jsonDiffCmd.InheritedFlags().SortFlags = false

	return jsonDiffCmd
}
