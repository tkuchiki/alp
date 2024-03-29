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

func newRegexpCmd(flags *flags) *cobra.Command {
	var regexpCmd = &cobra.Command{
		Use:   "regexp",
		Short: "Profile the log_reader that match a regular expression",
		Long:  `Profile the log_reader that match a regular expression`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := flags.createRegexpOptions(cmd)
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

			parser, err := newRegexpParser(opts, f)
			if err != nil {
				return err
			}

			err = prof.Run(flags.sortOptions, parser, nil)

			return err
		},
	}

	flags.defineProfileOptions(regexpCmd)
	flags.defineRegexpOptions(regexpCmd)

	regexpCmd.Flags().SortFlags = false
	regexpCmd.PersistentFlags().SortFlags = false
	regexpCmd.InheritedFlags().SortFlags = false

	return regexpCmd
}

func newRegexpParser(opts *options.Options, f *os.File) (parsers.Parser, error) {
	names := parsers.NewSubexpNames(opts.Regexp.UriSubexp, opts.Regexp.MethodSubexp, opts.Regexp.TimeSubexp,
		opts.Regexp.ResponseTimeSubexp, opts.Regexp.RequestTimeSubexp, opts.Regexp.BodyBytesSubexp, opts.Regexp.StatusSubexp)
	return parsers.NewRegexpParser(f, opts.Regexp.Pattern, names, opts.QueryString, opts.QueryStringIgnoreValues)
}

func newRegexpDiffCmd(flags *flags) *cobra.Command {
	regexpDiffCmd := newDiffSubCmd()
	regexpDiffCmd.RunE = func(cmd *cobra.Command, args []string) error {
		opts, err := flags.createRegexpDiffOptions(cmd)
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

		fromParser, err := newRegexpParser(opts, fromf)
		if err != nil {
			return err
		}

		toProf := profiler.NewProfiler(os.Stdout, os.Stderr, opts)
		toProf.DisableLoad()

		tof, err := toProf.Open(to)
		if err != nil {
			return err
		}
		defer tof.Close()

		toParser, err := newRegexpParser(opts, tof)
		if err != nil {
			return err
		}

		return runDiff(flags.sortOptions,
			fromProf, fromParser,
			toProf, toParser,
		)
	}

	flags.defineDiffSubCommandOptions(regexpDiffCmd)
	flags.defineRegexpOptions(regexpDiffCmd)

	regexpDiffCmd.Flags().SortFlags = false
	regexpDiffCmd.PersistentFlags().SortFlags = false
	regexpDiffCmd.InheritedFlags().SortFlags = false

	return regexpDiffCmd
}

func newRegexpTopNCmd(flags *flags) *cobra.Command {
	regexpTopNCmd := newTopNSubCmd()
	regexpTopNCmd.RunE = func(cmd *cobra.Command, args []string) error {
		opts, err := flags.createRegexpTopNOptions(cmd)
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

		parser, err := newRegexpParser(opts, f)
		if err != nil {
			return err
		}

		return runTopN(logReader, parser)
	}

	flags.defineTopNSubCommandOptions(regexpTopNCmd)
	flags.defineRegexpOptions(regexpTopNCmd)

	regexpTopNCmd.Flags().SortFlags = false
	regexpTopNCmd.PersistentFlags().SortFlags = false
	regexpTopNCmd.InheritedFlags().SortFlags = false

	return regexpTopNCmd
}

func newRegexpCountCmd(flags *flags) *cobra.Command {
	regexpCountCmd := newCountSubCmd()
	regexpCountCmd.RunE = func(cmd *cobra.Command, args []string) error {
		opts, err := flags.createRegexpCountOptions(cmd)
		if err != nil {
			return err
		}

		counter := counter.NewCounter(os.Stdout, os.Stderr, opts)

		f, err := counter.Open(opts.File)
		if err != nil {
			return err
		}
		defer f.Close()

		parser, err := newRegexpParser(opts, f)
		if err != nil {
			return err
		}

		return runCount(counter, parser, opts)
	}

	flags.defineCountSubCommandOptions(regexpCountCmd)
	flags.defineRegexpOptions(regexpCountCmd)

	regexpCountCmd.Flags().SortFlags = false
	regexpCountCmd.PersistentFlags().SortFlags = false
	regexpCountCmd.InheritedFlags().SortFlags = false

	return regexpCountCmd
}
