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

func NewLTSVCmd() *cobra.Command {
	var ltsvCmd = &cobra.Command{
		Use:   "ltsv",
		Short: "Profile the logs for LTSV",
		Long:  `Profile the logs for LTSV`,
		RunE: func(cmd *cobra.Command, args []string) error {
			sortOptions := stats.NewSortOptions()
			opts, err := createLTSVOptions(cmd, sortOptions)
			if err != nil {
				return err
			}

			prof := profiler.NewProfiler(os.Stdout, os.Stderr, opts)

			f, err := prof.Open(opts.File)
			if err != nil {
				return err
			}
			defer f.Close()

			label := parsers.NewLTSVLabel(opts.LTSV.UriLabel, opts.LTSV.MethodLabel, opts.LTSV.TimeLabel,
				opts.LTSV.ApptimeLabel, opts.LTSV.ReqtimeLabel, opts.LTSV.SizeLabel, opts.LTSV.StatusLabel,
			)
			parser := parsers.NewLTSVParser(f, label, opts.QueryString, opts.QueryStringIgnoreValues)

			err = prof.Run(sortOptions, parser)

			return err
		},
	}

	defineOptions(ltsvCmd)

	ltsvCmd.PersistentFlags().StringP("uri-label", "", options.DefaultUriLabelOption, "Change the uri label")
	ltsvCmd.PersistentFlags().StringP("method-label", "", options.DefaultMethodLabelOption, "Change the method label")
	ltsvCmd.PersistentFlags().StringP("time-label", "", options.DefaultTimeLabelOption, "Change the time label")
	ltsvCmd.PersistentFlags().StringP("apptime-label", "", options.DefaultApptimeLabelOption, "Change the apptime label")
	ltsvCmd.PersistentFlags().StringP("reqtime-label", "", options.DefaultReqtimeLabelOption, "Change the reqtime label")
	ltsvCmd.PersistentFlags().StringP("size-label", "", options.DefaultSizeLabelOption, "Change the size label")
	ltsvCmd.PersistentFlags().StringP("status-label", "", options.DefaultStatusLabelOption, "Change the status label")

	return ltsvCmd
}

func createLTSVOptions(cmd *cobra.Command, sortOptions *stats.SortOptions) (*options.Options, error) {
	config, err := cmd.PersistentFlags().GetString("config")
	if err != nil {
		return nil, err
	}
	if config != "" {
		bindCommonFlags(cmd)
		bindLTSVFlags(cmd)
		return createOptionsFromConfig(cmd, sortOptions, config)
	}

	opts, err := createCommonOptionsFromFlags(cmd, sortOptions)
	if err != nil {
		return nil, err
	}

	uriLabel, err := cmd.PersistentFlags().GetString("uri-label")
	if err != nil {
		return nil, err
	}

	methodLabel, err := cmd.PersistentFlags().GetString("method-label")
	if err != nil {
		return nil, err
	}

	timeLabel, err := cmd.PersistentFlags().GetString("time-label")
	if err != nil {
		return nil, err
	}

	appTimeLabel, err := cmd.PersistentFlags().GetString("apptime-label")
	if err != nil {
		return nil, err
	}

	reqTimeLabel, err := cmd.PersistentFlags().GetString("reqtime-label")
	if err != nil {
		return nil, err
	}

	sizeLabel, err := cmd.PersistentFlags().GetString("size-label")
	if err != nil {
		return nil, err
	}

	statusLabel, err := cmd.PersistentFlags().GetString("status-label")
	if err != nil {
		return nil, err
	}

	return options.SetOptions(opts,
		// ltsv
		options.UriLabel(uriLabel),
		options.MethodLabel(methodLabel),
		options.TimeLabel(timeLabel),
		options.ApptimeLabel(appTimeLabel),
		options.ReqtimeLabel(reqTimeLabel),
		options.SizeLabel(sizeLabel),
		options.StatusLabel(statusLabel),
	), nil
}

func bindLTSVFlags(cmd *cobra.Command) {
	viper.BindPFlag("ltsv.uri_label", cmd.PersistentFlags().Lookup("uri-label"))
	viper.BindPFlag("ltsv.method_label", cmd.PersistentFlags().Lookup("method-label"))
	viper.BindPFlag("ltsv.time_label", cmd.PersistentFlags().Lookup("time-label"))
	viper.BindPFlag("ltsv.apptime_label", cmd.PersistentFlags().Lookup("apptime-label"))
	viper.BindPFlag("ltsv.reqtime_label", cmd.PersistentFlags().Lookup("reqtime-label"))
	viper.BindPFlag("ltsv.size_label", cmd.PersistentFlags().Lookup("size-label"))
	viper.BindPFlag("ltsv.status_label", cmd.PersistentFlags().Lookup("status-label"))
}
