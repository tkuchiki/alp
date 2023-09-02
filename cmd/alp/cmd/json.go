package cmd

import (
	"os"

	"github.com/tkuchiki/alp/options"
	"github.com/tkuchiki/alp/parsers"
	"github.com/tkuchiki/alp/profiler"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tkuchiki/alp/stats"
)

func NewJSONCmd(rootCmd *cobra.Command) *cobra.Command {
	var jsonCmd = &cobra.Command{
		Use:   "json",
		Short: "Profile the logs for JSON",
		Long:  `Profile the logs for JSON`,
		RunE: func(cmd *cobra.Command, args []string) error {
			sortOptions := stats.NewSortOptions()
			opts, err := createJSONOptions(cmd, sortOptions)
			if err != nil {
				return err
			}

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

func createJSONOptions(cmd *cobra.Command, sortOptions *stats.SortOptions) (*options.Options, error) {
	config, err := cmd.PersistentFlags().GetString("config")
	if err != nil {
		return nil, err
	}
	if config != "" {
		bindCommonFlags(cmd)
		bindJSONFlags(cmd)
		return createOptionsFromConfig(cmd, sortOptions, config)
	}

	opts, err := createCommonOptionsFromFlags(cmd, sortOptions)
	if err != nil {
		return nil, err
	}

	uriKey, err := cmd.PersistentFlags().GetString("uri-key")
	if err != nil {
		return nil, err
	}

	methodKey, err := cmd.PersistentFlags().GetString("method-key")
	if err != nil {
		return nil, err
	}

	timeKey, err := cmd.PersistentFlags().GetString("time-key")
	if err != nil {
		return nil, err
	}

	responseTimeKey, err := cmd.PersistentFlags().GetString("restime-key")
	if err != nil {
		return nil, err
	}

	requestTimeKey, err := cmd.PersistentFlags().GetString("reqtime-key")
	if err != nil {
		return nil, err
	}

	bodyBytesKey, err := cmd.PersistentFlags().GetString("body-bytes-key")
	if err != nil {
		return nil, err
	}

	statusKey, err := cmd.PersistentFlags().GetString("status-key")
	if err != nil {
		return nil, err
	}

	return options.SetOptions(opts,
		options.UriKey(uriKey),
		options.MethodKey(methodKey),
		options.TimeKey(timeKey),
		options.ResponseTimeKey(responseTimeKey),
		options.RequestTimeKey(requestTimeKey),
		options.BodyBytesKey(bodyBytesKey),
		options.StatusKey(statusKey),
	), nil
}

func bindJSONFlags(cmd *cobra.Command) {
	viper.BindPFlag("json.uri_key", cmd.PersistentFlags().Lookup("uri-key"))
	viper.BindPFlag("json.method_key", cmd.PersistentFlags().Lookup("method-key"))
	viper.BindPFlag("json.time_key", cmd.PersistentFlags().Lookup("time-key"))
	viper.BindPFlag("json.restime_key", cmd.PersistentFlags().Lookup("restime-key"))
	viper.BindPFlag("json.reqtime_key", cmd.PersistentFlags().Lookup("reqtime-key"))
	viper.BindPFlag("json.body_bytes_key", cmd.PersistentFlags().Lookup("body-bytes-key"))
	viper.BindPFlag("json.status_key", cmd.PersistentFlags().Lookup("status-key"))

}
