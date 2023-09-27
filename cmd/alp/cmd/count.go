package cmd

import (
	"os"

	"github.com/tkuchiki/alp/options"

	"github.com/spf13/cobra"
	"github.com/tkuchiki/alp/counter"
	"github.com/tkuchiki/alp/parsers"
)

func newCountCmd(flags *flags) *cobra.Command {
	var countCmd = &cobra.Command{
		Use:   "count",
		Short: "Count by log entries",
		Long:  `Count by log entries`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := flags.createCountOptions(cmd)
			if err != nil {
				return err
			}

			// TODO: start
			// Remove these after implementing `alp (json|ltsv|regex) count`.
			pattern, err := cmd.PersistentFlags().GetString("pattern")
			if err != nil {
				return err
			}

			format, err := cmd.PersistentFlags().GetString("format")
			if err != nil {
				return err
			}

			opts = options.SetOptions(opts,
				options.Pattern(pattern),
				options.Format(format),
			)
			// TODO: end

			cnter := counter.NewCounter(os.Stdout, os.Stderr, opts.Reverse)

			f, err := cnter.Open(opts.File)
			if err != nil {
				return err
			}
			defer f.Close()

			var parser parsers.Parser
			switch format {
			case "json":
				jsonKeys := parsers.NewJSONKeys(opts.JSON.UriKey, opts.JSON.MethodKey, opts.JSON.TimeKey,
					opts.JSON.ResponseTimeKey, opts.JSON.RequestTimeKey, opts.JSON.BodyBytesKey, opts.JSON.StatusKey)
				parser = parsers.NewJSONParser(f, jsonKeys, opts.QueryString, opts.QueryStringIgnoreValues)
			case "ltsv":
				label := parsers.NewLTSVLabel(opts.LTSV.UriLabel, opts.LTSV.MethodLabel, opts.LTSV.TimeLabel,
					opts.LTSV.ApptimeLabel, opts.LTSV.ReqtimeLabel, opts.LTSV.SizeLabel, opts.LTSV.StatusLabel,
				)
				parser = parsers.NewLTSVParser(f, label, opts.QueryString, opts.QueryStringIgnoreValues)
			case "regexp":
				names := parsers.NewSubexpNames(opts.Regexp.UriSubexp, opts.Regexp.MethodSubexp, opts.Regexp.TimeSubexp,
					opts.Regexp.ResponseTimeSubexp, opts.Regexp.RequestTimeSubexp, opts.Regexp.BodyBytesSubexp, opts.Regexp.StatusSubexp)
				parser, err = parsers.NewRegexpParser(f, opts.Regexp.Pattern, names, opts.QueryString, opts.QueryStringIgnoreValues)
				if err != nil {
					return err
				}
			}

			cnter.SetParser(parser)

			err = cnter.CountAndPrint(opts.Count.Keys)

			return err
		},
	}

	flags.defineCountOptions(countCmd)

	// TODO: Remove these after implementing `alp (json|ltsv|regex) count`.
	countCmd.PersistentFlags().StringP("pattern", "", options.DefaultPatternOption, "Regular expressions pattern matching the log")
	countCmd.PersistentFlags().StringP("format", "", "json", "Log format (json,ltsv,regexp)")

	countCmd.Flags().SortFlags = false
	countCmd.PersistentFlags().SortFlags = false
	countCmd.InheritedFlags().SortFlags = false

	return countCmd
}
