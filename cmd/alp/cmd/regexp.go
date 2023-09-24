package cmd

import (
	"os"

	"github.com/tkuchiki/alp/parsers"
	"github.com/tkuchiki/alp/profiler"

	"github.com/spf13/cobra"
)

func NewRegexpCmd(commandFlags *flags) *cobra.Command {
	var regexpCmd = &cobra.Command{
		Use:   "regexp",
		Short: "Profile the logs that match a regular expression",
		Long:  `Profile the logs that match a regular expression`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := commandFlags.createRegexpOptions(cmd)
			if err != nil {
				return err
			}

			prof := profiler.NewProfiler(os.Stdout, os.Stderr, opts)

			f, err := prof.Open(opts.File)
			if err != nil {
				return err
			}
			defer f.Close()

			names := parsers.NewSubexpNames(opts.Regexp.UriSubexp, opts.Regexp.MethodSubexp, opts.Regexp.TimeSubexp,
				opts.Regexp.ResponseTimeSubexp, opts.Regexp.RequestTimeSubexp, opts.Regexp.BodyBytesSubexp, opts.Regexp.StatusSubexp)
			parser, err := parsers.NewRegexpParser(f, opts.Regexp.Pattern, names, opts.QueryString, opts.QueryStringIgnoreValues)
			if err != nil {
				return err
			}

			err = prof.Run(commandFlags.sortOptions, parser)

			return err
		},
	}

	commandFlags.defineProfileOptions(regexpCmd)
	commandFlags.defineRegexpOptions(regexpCmd)

	regexpCmd.Flags().SortFlags = false
	regexpCmd.PersistentFlags().SortFlags = false
	regexpCmd.InheritedFlags().SortFlags = false

	return regexpCmd
}
