package cmd

import (
	"os"

	"github.com/tkuchiki/alp/options"
	"github.com/tkuchiki/alp/parsers"
	"github.com/tkuchiki/alp/profiler"

	"github.com/spf13/cobra"
	"github.com/tkuchiki/alp/stats"
)

func NewJSONCmd(rootCmd *cobra.Command) *cobra.Command {
	var jsonCmd = &cobra.Command{
		Use:   "json",
		Short: "Profile the logs for JSON",
		Long:  `Profile the logs for JSON`,
		RunE: func(cmd *cobra.Command, args []string) error {
			sortOptions := stats.NewSortOptions()
			opts, err := createOptions(cmd, sortOptions)
			if err != nil {
				return err
			}

			uriKey, err := cmd.PersistentFlags().GetString("uri-key")
			if err != nil {
				return err
			}

			methodKey, err := cmd.PersistentFlags().GetString("method-key")
			if err != nil {
				return err
			}

			timeKey, err := cmd.PersistentFlags().GetString("time-key")
			if err != nil {
				return err
			}

			responseTimeKey, err := cmd.PersistentFlags().GetString("restime-key")
			if err != nil {
				return err
			}

			requestTimeKey, err := cmd.PersistentFlags().GetString("reqtime-key")
			if err != nil {
				return err
			}

			bodyBytesKey, err := cmd.PersistentFlags().GetString("body-bytes-key")
			if err != nil {
				return err
			}

			statusKey, err := cmd.PersistentFlags().GetString("status-key")
			if err != nil {
				return err
			}

			opts = options.SetOptions(opts,
				options.UriKey(uriKey),
				options.MethodKey(methodKey),
				options.TimeKey(timeKey),
				options.ResponseTimeKey(responseTimeKey),
				options.RequestTimeKey(requestTimeKey),
				options.BodyBytesKey(bodyBytesKey),
				options.StatusKey(statusKey),
			)

			prof := profiler.NewProfiler(os.Stdout, os.Stderr, opts)

			f, err := prof.Open(opts.File)
			if err != nil {
				return err
			}
			defer f.Close()

			keys := parsers.NewJSONKeys(opts.JSON.UriKey, opts.JSON.MethodKey, opts.JSON.TimeKey,
				opts.JSON.ResponseTimeKey, opts.JSON.RequestTimeKey, opts.JSON.BodyBytesKey, opts.JSON.StatusKey)
			parser := parsers.NewJSONParser(f, keys, opts.QueryString, opts.QueryStringIgnoreValues)

			err = prof.Run(sortOptions, parser)

			return err
		},
	}

	defineOptions(jsonCmd)
	jsonCmd.PersistentFlags().StringP("uri-key", "", options.DefaultUriKeyOption, "Change the uri key")
	jsonCmd.PersistentFlags().StringP("method-key", "", options.DefaultMethodKeyOption, "Change the method key")
	jsonCmd.PersistentFlags().StringP("time-key", "", options.DefaultTimeKeyOption, "Change the time key")
	jsonCmd.PersistentFlags().StringP("restime-key", "", options.DefaultResponseTimeKeyOption, "Change the response_time key")
	jsonCmd.PersistentFlags().StringP("reqtime-key", "", options.DefaultRequestTimeKeyOption, "Change the request_time key")
	jsonCmd.PersistentFlags().StringP("body-bytes-key", "", options.DefaultBodyBytesKeyOption, "Change the body_bytes key")
	jsonCmd.PersistentFlags().StringP("status-key", "", options.DefaultStatusKeyOption, "Change the status key")

	return jsonCmd
}
