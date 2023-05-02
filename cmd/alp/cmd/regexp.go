package cmd

import (
	"os"

	"github.com/tkuchiki/alp/options"
	"github.com/tkuchiki/alp/parsers"
	"github.com/tkuchiki/alp/profiler"

	"github.com/spf13/cobra"
	"github.com/tkuchiki/alp/stats"
)

func NewRegexpCmd(rootCmd *cobra.Command) *cobra.Command {
	var regexpCmd = &cobra.Command{
		Use:   "regexp",
		Short: "Profile the logs that match a regular expression",
		Long:  `Profile the logs that match a regular expression`,
		RunE: func(cmd *cobra.Command, args []string) error {
			sortOptions := stats.NewSortOptions()
			opts, err := createOptions(cmd, sortOptions)
			if err != nil {
				return err
			}

			pattern, err := cmd.PersistentFlags().GetString("pattern")
			if err != nil {
				return err
			}

			uriSubexp, err := cmd.PersistentFlags().GetString("uri-subexp")
			if err != nil {
				return err
			}

			methodSubexp, err := cmd.PersistentFlags().GetString("method-subexp")
			if err != nil {
				return err
			}

			timeSubexp, err := cmd.PersistentFlags().GetString("time-subexp")
			if err != nil {
				return err
			}

			restimeSubexp, err := cmd.PersistentFlags().GetString("restime-subexp")
			if err != nil {
				return err
			}

			reqtimeSubexp, err := cmd.PersistentFlags().GetString("reqtime-subexp")
			if err != nil {
				return err
			}

			bodyBytesSubexp, err := cmd.PersistentFlags().GetString("body-bytes-subexp")
			if err != nil {
				return err
			}

			statusSubexp, err := cmd.PersistentFlags().GetString("status-subexp")
			if err != nil {
				return err
			}

			opts = options.SetOptions(opts,
				options.Pattern(pattern),
				options.UriSubexp(uriSubexp),
				options.MethodSubexp(methodSubexp),
				options.TimeSubexp(timeSubexp),
				options.ResponseTimeSubexp(restimeSubexp),
				options.RequestTimeSubexp(reqtimeSubexp),
				options.BodyBytesSubexp(bodyBytesSubexp),
				options.StatusSubexp(statusSubexp),
			)

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

			err = prof.Run(sortOptions, parser)

			return err
		},
	}

	defineOptions(regexpCmd)

	regexpCmd.PersistentFlags().StringP("pattern", "", options.DefaultPatternOption, "Regular expressions pattern matching the log")
	regexpCmd.PersistentFlags().StringP("uri-subexp", "", options.DefaultUriSubexpOption, "Change the uri sub expression")
	regexpCmd.PersistentFlags().StringP("method-subexp", "", options.DefaultMethodSubexpOption, "Change the method sub expression")
	regexpCmd.PersistentFlags().StringP("time-subexp", "", options.DefaultTimeSubexpOption, "Change the time sub expression")
	regexpCmd.PersistentFlags().StringP("restime-subexp", "", options.DefaultResponseTimeSubexpOption, "Change the response_time sub expression")
	regexpCmd.PersistentFlags().StringP("reqtime-subexp", "", options.DefaultRequestTimeSubexpOption, "Change the request_time sub expression")
	regexpCmd.PersistentFlags().StringP("body-bytes-subexp", "", options.DefaultBodyBytesSubexpOption, "Change the body_bytes sub expression")
	regexpCmd.PersistentFlags().StringP("status-subexp", "", options.DefaultStatusSubexpOption, "Change the status sub expression")

	return regexpCmd
}
