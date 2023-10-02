package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tkuchiki/alp/helpers"
	"github.com/tkuchiki/alp/options"
	"github.com/tkuchiki/alp/stats"
)

const (
	flagConfig                  = "config"
	flagFile                    = "file"
	flagDump                    = "dump"
	flagLoad                    = "load"
	flagFormat                  = "format"
	flagSort                    = "sort"
	flagReverse                 = "reverse"
	flagNoHeaders               = "noheaders"
	flagShowFooters             = "show-footers"
	flagLimit                   = "limit"
	flagOutput                  = "output"
	flagQueryString             = "query-string"
	flagQueryStringIgnoreValues = "qs-ignore-values"
	flagLocation                = "location"
	flagDecodeUri               = "decode-uri"
	flagMatchingGroups          = "matching-groups"
	flagFilters                 = "filters"
	flagPositionFile            = "pos"
	flagNoSavePositionFile      = "nosave-pos"
	flagPercentiles             = "percentiles"
	flagPage                    = "page"

	// json
	flagJSONUriKey       = "uri-key"
	flagJSONMethodKey    = "method-key"
	flagJSONTimeKey      = "time-key"
	flagJSONRestimeKey   = "restime-key"
	flagJSONReqtimeKey   = "reqtime-key"
	flagJSONBodyBytesKey = "body-bytes-key"
	flagJSONStatusKey    = "status-key"

	// ltsv
	flagLTSVUriLabel     = "uri-label"
	flagLTSVMethodLabel  = "method-label"
	flagLTSVTimeLabel    = "time-label"
	flagLTSVApptimeLabel = "apptime-label"
	flagLTSVReqtimeLabel = "reqtime-label"
	flagLTSVSizeLabel    = "size-label"
	flagLTSVStatusLabel  = "status-label"

	// regexp
	flagRegexpPattern         = "pattern"
	flagRegexpUriSubexp       = "uri-subexp"
	flagRegexpMethodSubexp    = "method-subexp"
	flagRegexpTimeSubexp      = "time-subexp"
	flagRegexpRestimeSubexp   = "restime-subexp"
	flagRegexpReqtimeSubexp   = "reqtime-subexp"
	flagRegexpBodyBytesSubexp = "body-bytes-subexp"
	flagRegexpStatusSubexp    = "status-subexp"

	// pcap
	flagPcapPcapServerIP   = "pcap-server-ip"
	flagPcapPcapServerPort = "pcap-server-port"

	// count
	flagCountKeys = "keys"
)

type flags struct {
	config      string
	sortOptions *stats.SortOptions
}

func newFlags() *flags {
	return &flags{
		sortOptions: stats.NewSortOptions(),
	}
}

func (f *flags) defineConfig(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&f.config, flagConfig, "", "The configuration file")
}

func (f *flags) defineFile(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagFile, "", "", "The access log file")
}

func (f *flags) defineDump(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagDump, "", "", "Dump profiled data as YAML")
}

func (f *flags) defineLoad(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagLoad, "", "", "Load the profiled YAML data")
}

func (f *flags) defineFormat(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagFormat, "", options.DefaultFormatOption, "The output format (table, markdown, tsv, csv, html, and json)")
}

func (f *flags) defineSort(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagSort, "", options.DefaultSortOption, "Output the results in sorted order")
}

func (f *flags) defineTopNSort(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagSort, "", options.DefaultTopNSortOption, "Output the results in sorted order (restime or bytes)")
}

func (f *flags) defineReverse(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolP(flagReverse, "r", false, "Sort results in reverse order")
}

func (f *flags) defineNoHeaders(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolP(flagNoHeaders, "", false, "Output no header line at all (only --format=tsv, csv)")
}

func (f *flags) defineShowFooters(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolP(flagShowFooters, "", false, "Output footer line at all (only --format=table, markdown)")
}

func (f *flags) defineLimit(cmd *cobra.Command) {
	cmd.PersistentFlags().IntP(flagLimit, "", options.DefaultLimitOption, "The maximum number of results to display")
}

func (f *flags) defineOutput(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagOutput, "o", options.DefaultOutputOption, "Specifies the results to display, separated by commas")
}

func (f *flags) defineQueryString(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolP(flagQueryString, "q", false, "Include the URI query string")
}

func (f *flags) defineQueryStringIgnoreValues(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolP(flagQueryStringIgnoreValues, "", false, "Ignore the value of the query string. Replace all values with xxx (only use with -q)")
}

func (f *flags) defineLocation(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagLocation, "", options.DefaultLocationOption, "Location name for the timezone")
}

func (f *flags) defineDecodeUri(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolP(flagDecodeUri, "", false, "Decode the URI")
}

func (f *flags) defineMatchingGroups(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagMatchingGroups, "m", "", "Specifies Query matching groups separated by commas")
}

func (f *flags) defineFilters(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagFilters, "f", "", "Only the log_reader are profiled that match the conditions")
}

func (f *flags) definePositionFile(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagPositionFile, "", "", "The position file")
}

func (f *flags) defineNoSavePositionFile(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolP(flagNoSavePositionFile, "", false, "Do not save position file")
}

func (f *flags) definePercentiles(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagPercentiles, "", "", "Specifies the percentiles separated by commas")
}

func (f *flags) definePage(cmd *cobra.Command) {
	cmd.PersistentFlags().IntP(flagPage, "", options.DefaultPaginationLimit, "Number of pages of pagination")
}

func (f *flags) defineJSONUriKey(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagJSONUriKey, "", options.DefaultUriKeyOption, "Change the uri key")
}

func (f *flags) defineJSONMethodKey(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagJSONMethodKey, "", options.DefaultMethodKeyOption, "Change the method key")
}

func (f *flags) defineJSONTimeKey(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagJSONTimeKey, "", options.DefaultTimeKeyOption, "Change the time key")
}

func (f *flags) defineJSONRestimeKey(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagJSONRestimeKey, "", options.DefaultResponseTimeKeyOption, "Change the response_time key")
}

func (f *flags) defineJSONReqtimeKey(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagJSONReqtimeKey, "", options.DefaultRequestTimeKeyOption, "Change the request_time key")
}

func (f *flags) defineJSONBodyBytesKey(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagJSONBodyBytesKey, "", options.DefaultBodyBytesKeyOption, "Change the body_bytes key")
}

func (f *flags) defineJSONStatusKey(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagJSONStatusKey, "", options.DefaultStatusKeyOption, "Change the status key")
}

func (f *flags) defineLTSVUriLabel(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagLTSVUriLabel, "", options.DefaultUriLabelOption, "Change the uri label")
}

func (f *flags) defineLTSVMethodLabel(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagLTSVMethodLabel, "", options.DefaultMethodLabelOption, "Change the method label")
}

func (f *flags) defineLTSVTimeLabel(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagLTSVTimeLabel, "", options.DefaultTimeLabelOption, "Change the time label")
}

func (f *flags) defineLTSVApptimeLabel(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagLTSVApptimeLabel, "", options.DefaultApptimeLabelOption, "Change the apptime label")
}

func (f *flags) defineLTSVReqtimeLabel(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagLTSVReqtimeLabel, "", options.DefaultReqtimeLabelOption, "Change the reqtime label")
}

func (f *flags) defineLTSVSizeLabel(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagLTSVSizeLabel, "", options.DefaultSizeLabelOption, "Change the size label")
}

func (f *flags) defineLTSVStatusLabel(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagLTSVStatusLabel, "", options.DefaultStatusLabelOption, "Change the status label")
}

func (f *flags) defineRegexpPattern(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagRegexpPattern, "", options.DefaultPatternOption, "Regular expressions pattern matching the log")
}

func (f *flags) defineRegexpUriSubexp(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagRegexpUriSubexp, "", options.DefaultUriSubexpOption, "Change the uri sub expression")
}

func (f *flags) defineRegexpMethodSubexp(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagRegexpMethodSubexp, "", options.DefaultMethodSubexpOption, "Change the method sub expression")
}

func (f *flags) defineRegexpTimeSubexp(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagRegexpTimeSubexp, "", options.DefaultTimeSubexpOption, "Change the time sub expression")
}

func (f *flags) defineRegexpRestimeSubexp(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagRegexpRestimeSubexp, "", options.DefaultResponseTimeSubexpOption, "Change the response_time sub expression")
}

func (f *flags) defineRegexpReqtimeSubexp(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagRegexpReqtimeSubexp, "", options.DefaultRequestTimeSubexpOption, "Change the request_time sub expression")
}

func (f *flags) defineRegexpBodyBytesSubexp(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagRegexpBodyBytesSubexp, "", options.DefaultBodyBytesSubexpOption, "Change the body_bytes sub expression")
}

func (f *flags) defineRegexpStatusSubexp(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagRegexpStatusSubexp, "", options.DefaultStatusSubexpOption, "Change the status sub expression")
}

func (f *flags) definePcapPcapServerIP(cmd *cobra.Command) {
	cmd.PersistentFlags().StringSliceP(flagPcapPcapServerIP, "", []string{options.DefaultPcapServerIPsOption[0]}, "HTTP server IP address of the captured packets")
}

func (f *flags) definePcapPcapServerPort(cmd *cobra.Command) {
	cmd.PersistentFlags().Uint16P(flagPcapPcapServerPort, "", options.DefaultPcapServerPortOption, "HTTP server TCP port of the captured packets")
}

func (f *flags) defineCountKeys(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagCountKeys, "", "", "Log key names (comma separated)")
	cmd.MarkPersistentFlagRequired(flagCountKeys)
}

func (f *flags) defineGlobalOptions(cmd *cobra.Command) {
	f.defineConfig(cmd)
}

func (f *flags) defineProfileOptions(cmd *cobra.Command) {
	f.defineFile(cmd)
	f.defineDump(cmd)
	f.defineLoad(cmd)
	f.defineFormat(cmd)
	f.defineSort(cmd)
	f.defineReverse(cmd)
	f.defineNoHeaders(cmd)
	f.defineShowFooters(cmd)
	f.defineLimit(cmd)
	f.defineOutput(cmd)
	f.defineQueryString(cmd)
	f.defineQueryStringIgnoreValues(cmd)
	f.defineLocation(cmd)
	f.defineDecodeUri(cmd)
	f.defineMatchingGroups(cmd)
	f.defineFilters(cmd)
	f.definePositionFile(cmd)
	f.defineNoSavePositionFile(cmd)
	f.definePercentiles(cmd)
	f.definePage(cmd)
}

func (f *flags) defineJSONOptions(cmd *cobra.Command) {
	f.defineJSONUriKey(cmd)
	f.defineJSONMethodKey(cmd)
	f.defineJSONTimeKey(cmd)
	f.defineJSONRestimeKey(cmd)
	f.defineJSONReqtimeKey(cmd)
	f.defineJSONBodyBytesKey(cmd)
	f.defineJSONStatusKey(cmd)
}

func (f *flags) defineLTSVOptions(cmd *cobra.Command) {
	f.defineLTSVUriLabel(cmd)
	f.defineLTSVMethodLabel(cmd)
	f.defineLTSVTimeLabel(cmd)
	f.defineLTSVApptimeLabel(cmd)
	f.defineLTSVReqtimeLabel(cmd)
	f.defineLTSVSizeLabel(cmd)
	f.defineLTSVStatusLabel(cmd)
}

func (f *flags) defineRegexpOptions(cmd *cobra.Command) {
	f.defineRegexpPattern(cmd)
	f.defineRegexpUriSubexp(cmd)
	f.defineRegexpMethodSubexp(cmd)
	f.defineRegexpTimeSubexp(cmd)
	f.defineRegexpRestimeSubexp(cmd)
	f.defineRegexpReqtimeSubexp(cmd)
	f.defineRegexpBodyBytesSubexp(cmd)
	f.defineRegexpStatusSubexp(cmd)
}

func (f *flags) definePcapOptions(cmd *cobra.Command) {
	f.definePcapPcapServerIP(cmd)
	f.definePcapPcapServerPort(cmd)
}

func (f *flags) defineDiffOptions(cmd *cobra.Command) {
	f.defineFormat(cmd)
	f.defineSort(cmd)
	f.defineReverse(cmd)
	f.defineNoHeaders(cmd)
	f.defineShowFooters(cmd)
	f.defineLimit(cmd)
	f.defineOutput(cmd)
	f.defineQueryString(cmd)
	f.defineQueryStringIgnoreValues(cmd)
	f.defineLocation(cmd)
	f.defineDecodeUri(cmd)
	f.defineMatchingGroups(cmd)
	f.defineFilters(cmd)
	f.definePercentiles(cmd)
	f.definePage(cmd)
}

func (f *flags) defineDiffSubCommandOptions(cmd *cobra.Command) {
	// overwrite and hidden => remove flag
	cmd.LocalFlags().String(flagFile, "", "")
	cmd.LocalFlags().MarkHidden(flagFile)

	f.defineDump(cmd)
	f.defineLoad(cmd)
	f.defineFormat(cmd)
	f.defineSort(cmd)
	f.defineReverse(cmd)
	f.defineNoHeaders(cmd)
	f.defineShowFooters(cmd)
	f.defineLimit(cmd)
	f.defineOutput(cmd)
	f.defineQueryString(cmd)
	f.defineQueryStringIgnoreValues(cmd)
	f.defineLocation(cmd)
	f.defineDecodeUri(cmd)
	f.defineMatchingGroups(cmd)
	f.defineFilters(cmd)
	f.definePositionFile(cmd)
	f.defineNoSavePositionFile(cmd)
	f.definePercentiles(cmd)
	f.definePage(cmd)
}

func (f *flags) defineTopNSubCommandOptions(cmd *cobra.Command) {
	// overwrite and hidden => remove flag
	cmd.LocalFlags().String(flagDump, "", "")
	cmd.LocalFlags().MarkHidden(flagDump)
	cmd.LocalFlags().String(flagLoad, "", "")
	cmd.LocalFlags().MarkHidden(flagLoad)
	cmd.LocalFlags().String(flagShowFooters, "", "")
	cmd.LocalFlags().MarkHidden(flagShowFooters)
	cmd.LocalFlags().String(flagLimit, "", "")
	cmd.LocalFlags().MarkHidden(flagLimit)
	cmd.LocalFlags().String(flagOutput, "", "")
	cmd.LocalFlags().MarkHidden(flagOutput)
	cmd.LocalFlags().String(flagQueryString, "", "")
	cmd.LocalFlags().MarkHidden(flagQueryString)
	cmd.LocalFlags().String(flagQueryStringIgnoreValues, "", "")
	cmd.LocalFlags().MarkHidden(flagQueryStringIgnoreValues)
	cmd.LocalFlags().String(flagMatchingGroups, "", "")
	cmd.LocalFlags().MarkHidden(flagMatchingGroups)
	cmd.LocalFlags().String(flagPercentiles, "", "")
	cmd.LocalFlags().MarkHidden(flagPercentiles)

	f.defineFile(cmd)
	f.defineFormat(cmd)
	f.defineTopNSort(cmd)
	f.defineReverse(cmd)
	f.defineNoHeaders(cmd)
	f.defineLocation(cmd)
	f.defineDecodeUri(cmd)
	f.defineFilters(cmd)
	f.definePositionFile(cmd)
	f.defineNoSavePositionFile(cmd)
	f.definePage(cmd)
}

func (f *flags) defineCountSubCommandOptions(cmd *cobra.Command) {
	// overwrite and hidden => remove flag
	cmd.LocalFlags().String(flagDump, "", "")
	cmd.LocalFlags().MarkHidden(flagDump)
	cmd.LocalFlags().String(flagLoad, "", "")
	cmd.LocalFlags().MarkHidden(flagLoad)
	cmd.LocalFlags().String(flagSort, "", "")
	cmd.LocalFlags().MarkHidden(flagSort)
	cmd.LocalFlags().String(flagShowFooters, "", "")
	cmd.LocalFlags().MarkHidden(flagShowFooters)
	cmd.LocalFlags().String(flagLimit, "", "")
	cmd.LocalFlags().MarkHidden(flagLimit)
	cmd.LocalFlags().String(flagOutput, "", "")
	cmd.LocalFlags().MarkHidden(flagOutput)
	cmd.LocalFlags().String(flagQueryString, "", "")
	cmd.LocalFlags().MarkHidden(flagQueryString)
	cmd.LocalFlags().String(flagQueryStringIgnoreValues, "", "")
	cmd.LocalFlags().MarkHidden(flagQueryStringIgnoreValues)
	cmd.LocalFlags().String(flagLocation, "", "")
	cmd.LocalFlags().MarkHidden(flagLocation)
	cmd.LocalFlags().String(flagDecodeUri, "", "")
	cmd.LocalFlags().MarkHidden(flagDecodeUri)
	cmd.LocalFlags().String(flagMatchingGroups, "", "")
	cmd.LocalFlags().MarkHidden(flagMatchingGroups)
	cmd.LocalFlags().String(flagFilters, "", "")
	cmd.LocalFlags().MarkHidden(flagFilters)
	cmd.LocalFlags().String(flagPositionFile, "", "")
	cmd.LocalFlags().MarkHidden(flagPositionFile)
	cmd.LocalFlags().String(flagNoSavePositionFile, "", "")
	cmd.LocalFlags().MarkHidden(flagNoSavePositionFile)
	cmd.LocalFlags().String(flagPercentiles, "", "")
	cmd.LocalFlags().MarkHidden(flagPercentiles)

	f.defineFile(cmd)
	f.defineReverse(cmd)
	f.defineFormat(cmd)
	f.defineNoHeaders(cmd)
	f.definePage(cmd)
	f.defineCountKeys(cmd)
}

func (f *flags) bindFlags(cmd *cobra.Command) {
	viper.BindPFlag("file", cmd.PersistentFlags().Lookup(flagFile))
	viper.BindPFlag("dump", cmd.PersistentFlags().Lookup(flagDump))
	viper.BindPFlag("load", cmd.PersistentFlags().Lookup(flagLoad))

	if !strings.Contains(cmd.Name(), "topN") {
		viper.BindPFlag("sort", cmd.PersistentFlags().Lookup(flagSort))
		viper.BindPFlag("reverse", cmd.PersistentFlags().Lookup(flagReverse))
	}

	viper.BindPFlag("query_string", cmd.PersistentFlags().Lookup(flagQueryString))
	viper.BindPFlag("query_string_ignore_values", cmd.PersistentFlags().Lookup(flagQueryStringIgnoreValues))
	viper.BindPFlag("decode_uri", cmd.PersistentFlags().Lookup(flagDecodeUri))
	viper.BindPFlag("format", cmd.PersistentFlags().Lookup(flagFormat))
	viper.BindPFlag("noheaders", cmd.PersistentFlags().Lookup(flagNoHeaders))
	viper.BindPFlag("show_footers", cmd.PersistentFlags().Lookup(flagShowFooters))
	viper.BindPFlag("limit", cmd.PersistentFlags().Lookup(flagLimit))
	viper.BindPFlag("matching_groups", cmd.PersistentFlags().Lookup(flagMatchingGroups))
	viper.BindPFlag("filters", cmd.PersistentFlags().Lookup(flagFilters))
	viper.BindPFlag("pos_file", cmd.PersistentFlags().Lookup(flagPositionFile))
	viper.BindPFlag("nosave_pos", cmd.PersistentFlags().Lookup(flagNoSavePositionFile))
	viper.BindPFlag("location", cmd.PersistentFlags().Lookup(flagLocation))
	viper.BindPFlag("output", cmd.PersistentFlags().Lookup(flagOutput))
	viper.BindPFlag("pagenation_limit", cmd.PersistentFlags().Lookup(flagPage))

	// json
	viper.BindPFlag("json.uri_key", cmd.PersistentFlags().Lookup(flagJSONUriKey))
	viper.BindPFlag("json.method_key", cmd.PersistentFlags().Lookup(flagJSONMethodKey))
	viper.BindPFlag("json.time_key", cmd.PersistentFlags().Lookup(flagJSONTimeKey))
	viper.BindPFlag("json.restime_key", cmd.PersistentFlags().Lookup(flagJSONRestimeKey))
	viper.BindPFlag("json.reqtime_key", cmd.PersistentFlags().Lookup(flagJSONReqtimeKey))
	viper.BindPFlag("json.body_bytes_key", cmd.PersistentFlags().Lookup(flagJSONBodyBytesKey))
	viper.BindPFlag("json.status_key", cmd.PersistentFlags().Lookup(flagJSONStatusKey))

	// ltsv
	viper.BindPFlag("ltsv.uri_label", cmd.PersistentFlags().Lookup(flagLTSVUriLabel))
	viper.BindPFlag("ltsv.method_label", cmd.PersistentFlags().Lookup(flagLTSVMethodLabel))
	viper.BindPFlag("ltsv.time_label", cmd.PersistentFlags().Lookup(flagLTSVTimeLabel))
	viper.BindPFlag("ltsv.apptime_label", cmd.PersistentFlags().Lookup(flagLTSVApptimeLabel))
	viper.BindPFlag("ltsv.reqtime_label", cmd.PersistentFlags().Lookup(flagLTSVReqtimeLabel))
	viper.BindPFlag("ltsv.size_label", cmd.PersistentFlags().Lookup(flagLTSVSizeLabel))
	viper.BindPFlag("ltsv.status_label", cmd.PersistentFlags().Lookup(flagLTSVStatusLabel))

	// regexp
	viper.BindPFlag("regexp.pattern", cmd.PersistentFlags().Lookup(flagRegexpPattern))
	viper.BindPFlag("regexp.uri_subexp", cmd.PersistentFlags().Lookup(flagRegexpUriSubexp))
	viper.BindPFlag("regexp.method_subexp", cmd.PersistentFlags().Lookup(flagRegexpMethodSubexp))
	viper.BindPFlag("regexp.time_subexp", cmd.PersistentFlags().Lookup(flagRegexpTimeSubexp))
	viper.BindPFlag("regexp.restime_subexp", cmd.PersistentFlags().Lookup(flagRegexpRestimeSubexp))
	viper.BindPFlag("regexp.reqtime_subexp", cmd.PersistentFlags().Lookup(flagRegexpReqtimeSubexp))
	viper.BindPFlag("regexp.body_bytes_subexp", cmd.PersistentFlags().Lookup(flagRegexpBodyBytesSubexp))
	viper.BindPFlag("regexp.status_subexp", cmd.PersistentFlags().Lookup(flagRegexpStatusSubexp))

	// pcap
	viper.BindPFlag("pcap.server_port", cmd.PersistentFlags().Lookup(flagPcapPcapServerPort))

	// count
	viper.BindPFlag("count.keys", cmd.PersistentFlags().Lookup(flagCountKeys))

	// topN
	if strings.Contains(cmd.Name(), "topN") {
		viper.BindPFlag("topN.sort", cmd.PersistentFlags().Lookup(flagSort))
		viper.BindPFlag("topN.reverse", cmd.PersistentFlags().Lookup(flagReverse))
	}
}

func (f *flags) createOptionsFromConfig(cmd *cobra.Command) (*options.Options, error) {
	opts := options.NewOptions()
	viper.SetConfigFile(f.config)
	viper.SetConfigType("yaml")

	// Start workaround
	// viper seems to merge slices, so we'll set empty slice and overwrite it manually.
	opts.Percentiles = []int{}
	opts.Pcap.ServerIPs = []string{}

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(opts); err != nil {
		return nil, err
	}

	if len(opts.Percentiles) == 0 {
		opts.Percentiles = options.DefaultPercentilesOption
	}

	if len(opts.Pcap.ServerIPs) == 0 {
		opts.Pcap.ServerIPs = []string{options.DefaultPcapServerIPsOption[0]}
	}
	// End workaround

	percentilesFlag := cmd.PersistentFlags().Lookup(flagPercentiles)
	if percentilesFlag != nil && percentilesFlag.Changed {
		ps := cmd.PersistentFlags().Lookup(flagPercentiles).Value.String()
		var percentiles []int
		var err error
		if ps != "" {
			percentiles, err = helpers.SplitCSVIntoInts(ps)
			if err != nil {
				return nil, err
			}

			if err = helpers.ValidatePercentiles(percentiles); err != nil {
				return nil, err
			}
		}
		opts.Percentiles = percentiles
	}

	srvIPFlag := cmd.PersistentFlags().Lookup(flagPcapPcapServerIP)
	if srvIPFlag != nil && srvIPFlag.Changed {
		ips := cmd.PersistentFlags().Lookup(flagPcapPcapServerIP).Value.String()
		opts.Pcap.ServerIPs = helpers.SplitCSV(ips)
	}

	if err := f.sortOptions.SetAndValidate(opts.Sort); err != nil {
		return nil, err
	}
	opts.Sort = f.sortOptions.SortType()

	return opts, nil
}

func (f *flags) setOptions(cmd *cobra.Command, opts *options.Options, flags []string) (*options.Options, error) {
	for _, flag := range flags {
		switch flag {
		case flagFile:
			file, err := cmd.PersistentFlags().GetString(flagFile)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.File(file))
		case flagDump:
			dump, err := cmd.PersistentFlags().GetString(flagDump)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.Dump(dump))
		case flagLoad:
			load, err := cmd.PersistentFlags().GetString(flagLoad)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.Load(load))
		case flagFormat:
			format, err := cmd.PersistentFlags().GetString(flagFormat)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.Format(format))
		case flagSort:
			sort, err := cmd.PersistentFlags().GetString(flagSort)
			if err != nil {
				return nil, err
			}

			err = f.sortOptions.SetAndValidate(sort)
			if err != nil {
				return nil, err
			}

			opts = options.SetOptions(opts, options.Sort(sort))
		case flagReverse:
			reverse, err := cmd.PersistentFlags().GetBool(flagReverse)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.Reverse(reverse))
		case flagNoHeaders:
			noHeaders, err := cmd.PersistentFlags().GetBool(flagNoHeaders)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.NoHeaders(noHeaders))
		case flagShowFooters:
			showFooters, err := cmd.PersistentFlags().GetBool(flagShowFooters)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.ShowFooters(showFooters))
		case flagLimit:
			limit, err := cmd.PersistentFlags().GetInt(flagLimit)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.Limit(limit))
		case flagOutput:
			output, err := cmd.PersistentFlags().GetString(flagOutput)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.Output(output))
		case flagQueryString:
			queryString, err := cmd.PersistentFlags().GetBool(flagQueryString)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.QueryString(queryString))
		case flagQueryStringIgnoreValues:
			queryStringIgnoreValues, err := cmd.PersistentFlags().GetBool(flagQueryStringIgnoreValues)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.QueryStringIgnoreValues(queryStringIgnoreValues))
		case flagDecodeUri:
			decodeUri, err := cmd.PersistentFlags().GetBool(flagDecodeUri)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.DecodeUri(decodeUri))
		case flagLocation:
			location, err := cmd.PersistentFlags().GetString(flagLocation)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.Location(location))
		case flagMatchingGroups:
			matchingGroups, err := cmd.PersistentFlags().GetString(flagMatchingGroups)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.CSVGroups(matchingGroups))
		case flagFilters:
			filters, err := cmd.PersistentFlags().GetString(flagFilters)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.Filters(filters))
		case flagPositionFile:
			pos, err := cmd.PersistentFlags().GetString(flagPositionFile)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.PosFile(pos))
		case flagNoSavePositionFile:
			noSavePos, err := cmd.PersistentFlags().GetBool(flagNoSavePositionFile)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.NoSavePos(noSavePos))
		case flagPercentiles:
			ps, err := cmd.PersistentFlags().GetString(flagPercentiles)
			if err != nil {
				return nil, err
			}

			var percentiles []int
			if ps != "" {
				percentiles, err = helpers.SplitCSVIntoInts(ps)
				if err != nil {
					return nil, err
				}

				if err = helpers.ValidatePercentiles(percentiles); err != nil {
					return nil, err
				}
			}
			opts = options.SetOptions(opts, options.Percentiles(percentiles))
		case flagPage:
			paginationLimit, err := cmd.PersistentFlags().GetInt(flagPage)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.PaginationLimit(paginationLimit))
		}
	}

	return opts, nil
}

func (f *flags) setProfileOptions(cmd *cobra.Command, opts *options.Options) (*options.Options, error) {
	_flags := []string{
		flagFile,
		flagDump,
		flagLoad,
		flagFormat,
		flagSort,
		flagReverse,
		flagNoHeaders,
		flagShowFooters,
		flagLimit,
		flagOutput,
		flagQueryString,
		flagQueryStringIgnoreValues,
		flagLocation,
		flagDecodeUri,
		flagMatchingGroups,
		flagFilters,
		flagPositionFile,
		flagNoSavePositionFile,
		flagPercentiles,
		flagPage,
	}

	return f.setOptions(cmd, opts, _flags)
}

func (f *flags) setJSONOptions(cmd *cobra.Command, opts *options.Options) (*options.Options, error) {
	uriKey, err := cmd.PersistentFlags().GetString(flagJSONUriKey)
	if err != nil {
		return nil, err
	}

	methodKey, err := cmd.PersistentFlags().GetString(flagJSONMethodKey)
	if err != nil {
		return nil, err
	}

	timeKey, err := cmd.PersistentFlags().GetString(flagJSONTimeKey)
	if err != nil {
		return nil, err
	}

	responseTimeKey, err := cmd.PersistentFlags().GetString(flagJSONRestimeKey)
	if err != nil {
		return nil, err
	}

	requestTimeKey, err := cmd.PersistentFlags().GetString(flagJSONReqtimeKey)
	if err != nil {
		return nil, err
	}

	bodyBytesKey, err := cmd.PersistentFlags().GetString(flagJSONBodyBytesKey)
	if err != nil {
		return nil, err
	}

	statusKey, err := cmd.PersistentFlags().GetString(flagJSONStatusKey)
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

func (f *flags) setLTSVOptions(cmd *cobra.Command, opts *options.Options) (*options.Options, error) {
	uriLabel, err := cmd.PersistentFlags().GetString(flagLTSVUriLabel)
	if err != nil {
		return nil, err
	}

	methodLabel, err := cmd.PersistentFlags().GetString(flagLTSVMethodLabel)
	if err != nil {
		return nil, err
	}

	timeLabel, err := cmd.PersistentFlags().GetString(flagLTSVTimeLabel)
	if err != nil {
		return nil, err
	}

	appTimeLabel, err := cmd.PersistentFlags().GetString(flagLTSVApptimeLabel)
	if err != nil {
		return nil, err
	}

	reqTimeLabel, err := cmd.PersistentFlags().GetString(flagLTSVReqtimeLabel)
	if err != nil {
		return nil, err
	}

	sizeLabel, err := cmd.PersistentFlags().GetString(flagLTSVSizeLabel)
	if err != nil {
		return nil, err
	}

	statusLabel, err := cmd.PersistentFlags().GetString(flagLTSVStatusLabel)
	if err != nil {
		return nil, err
	}

	return options.SetOptions(opts,
		options.UriLabel(uriLabel),
		options.MethodLabel(methodLabel),
		options.TimeLabel(timeLabel),
		options.ApptimeLabel(appTimeLabel),
		options.ReqtimeLabel(reqTimeLabel),
		options.SizeLabel(sizeLabel),
		options.StatusLabel(statusLabel),
	), nil
}

func (f *flags) setRegexpOptions(cmd *cobra.Command, opts *options.Options) (*options.Options, error) {
	pattern, err := cmd.PersistentFlags().GetString(flagRegexpPattern)
	if err != nil {
		return nil, err
	}

	uriSubexp, err := cmd.PersistentFlags().GetString(flagRegexpUriSubexp)
	if err != nil {
		return nil, err
	}

	methodSubexp, err := cmd.PersistentFlags().GetString(flagRegexpMethodSubexp)
	if err != nil {
		return nil, err
	}

	timeSubexp, err := cmd.PersistentFlags().GetString(flagRegexpTimeSubexp)
	if err != nil {
		return nil, err
	}
	restimeSubexp, err := cmd.PersistentFlags().GetString(flagRegexpRestimeSubexp)
	if err != nil {
		return nil, err
	}

	reqtimeSubexp, err := cmd.PersistentFlags().GetString(flagRegexpReqtimeSubexp)
	if err != nil {
		return nil, err
	}

	bodyBytesSubexp, err := cmd.PersistentFlags().GetString(flagRegexpBodyBytesSubexp)
	if err != nil {
		return nil, err
	}

	statusSubexp, err := cmd.PersistentFlags().GetString(flagRegexpStatusSubexp)
	if err != nil {
		return nil, err
	}

	return options.SetOptions(opts,
		options.Pattern(pattern),
		options.UriSubexp(uriSubexp),
		options.MethodSubexp(methodSubexp),
		options.TimeSubexp(timeSubexp),
		options.ResponseTimeSubexp(restimeSubexp),
		options.RequestTimeSubexp(reqtimeSubexp),
		options.BodyBytesSubexp(bodyBytesSubexp),
		options.StatusSubexp(statusSubexp),
	), nil
}

func (f *flags) setPcapOptions(cmd *cobra.Command, opts *options.Options) (*options.Options, error) {
	serverIPs, err := cmd.PersistentFlags().GetStringSlice(flagPcapPcapServerIP)
	if err != nil {
		return nil, err
	}

	serverPort, err := cmd.PersistentFlags().GetUint16(flagPcapPcapServerPort)
	if err != nil {
		return nil, err
	}

	return options.SetOptions(opts,
		options.PcapServerIPs(serverIPs),
		options.PcapServerPort(serverPort),
	), nil
}

func (f *flags) setDiffOptions(cmd *cobra.Command, opts *options.Options) (*options.Options, error) {
	_flags := []string{
		flagFormat,
		flagSort,
		flagReverse,
		flagNoHeaders,
		flagShowFooters,
		flagLimit,
		flagOutput,
		flagQueryString,
		flagQueryStringIgnoreValues,
		flagLocation,
		flagDecodeUri,
		flagMatchingGroups,
		flagFilters,
		flagPercentiles,
		flagPage,
	}

	return f.setOptions(cmd, opts, _flags)
}

func (f *flags) setDiffSubCommandOptions(cmd *cobra.Command, opts *options.Options) (*options.Options, error) {
	_flags := []string{
		flagDump,
		flagLoad,
		flagFormat,
		flagSort,
		flagReverse,
		flagNoHeaders,
		flagShowFooters,
		flagLimit,
		flagOutput,
		flagQueryString,
		flagQueryStringIgnoreValues,
		flagLocation,
		flagDecodeUri,
		flagMatchingGroups,
		flagFilters,
		flagPositionFile,
		flagNoSavePositionFile,
		flagPercentiles,
		flagPage,
	}

	return f.setOptions(cmd, opts, _flags)
}

func (f *flags) setTopNSubCommandOptions(cmd *cobra.Command, opts *options.Options) (*options.Options, error) {
	_flags := []string{
		flagFile,
		flagFormat,
		flagReverse,
		flagNoHeaders,
		flagLocation,
		flagDecodeUri,
		flagFilters,
		flagPositionFile,
		flagNoSavePositionFile,
		flagPage,
	}

	sort, err := cmd.PersistentFlags().GetString(flagSort)
	if err != nil {
		return nil, err
	}
	opts = options.SetOptions(opts, options.TopNSort(sort))

	reverse, err := cmd.PersistentFlags().GetBool(flagReverse)
	if err != nil {
		return nil, err
	}
	opts = options.SetOptions(opts, options.TopNReverse(reverse))

	return f.setOptions(cmd, opts, _flags)
}

func (f *flags) setCountSubCommandOptions(cmd *cobra.Command, opts *options.Options) (*options.Options, error) {
	keys, err := cmd.PersistentFlags().GetString(flagCountKeys)
	if err != nil {
		return nil, err
	}

	opts = options.SetOptions(opts,
		options.CountKeys(helpers.SplitCSV(keys)),
	)

	_flags := []string{
		flagFile,
		flagFormat,
		flagReverse,
		flagNoHeaders,
		flagPage,
	}

	return f.setOptions(cmd, opts, _flags)
}

// alp json
func (f *flags) createJSONOptions(cmd *cobra.Command) (*options.Options, error) {
	if f.config != "" {
		f.bindFlags(cmd)
		return f.createOptionsFromConfig(cmd)
	}

	opts, err := f.setProfileOptions(cmd, options.NewOptions())
	if err != nil {
		return nil, err
	}

	return f.setJSONOptions(cmd, opts)
}

// alp json diff
func (f *flags) createJSONDiffOptions(cmd *cobra.Command) (*options.Options, error) {
	if f.config != "" {
		f.bindFlags(cmd)
		return f.createOptionsFromConfig(cmd)
	}

	opts, err := f.setDiffSubCommandOptions(cmd, options.NewOptions())
	if err != nil {
		return nil, err
	}

	return f.setJSONOptions(cmd, opts)
}

// alp json topN
func (f *flags) createJSONTopNOptions(cmd *cobra.Command) (*options.Options, error) {
	if f.config != "" {
		f.bindFlags(cmd)
		return f.createOptionsFromConfig(cmd)
	}

	opts, err := f.setTopNSubCommandOptions(cmd, options.NewOptions())
	if err != nil {
		return nil, err
	}

	return f.setJSONOptions(cmd, opts)
}

// alp json count
func (f *flags) createJSONCountOptions(cmd *cobra.Command) (*options.Options, error) {
	if f.config != "" {
		f.bindFlags(cmd)
		return f.createOptionsFromConfig(cmd)
	}

	opts, err := f.setCountSubCommandOptions(cmd, options.NewOptions())
	if err != nil {
		return nil, err
	}

	return f.setJSONOptions(cmd, opts)
}

// alp ltsv
func (f *flags) createLTSVOptions(cmd *cobra.Command) (*options.Options, error) {
	if f.config != "" {
		f.bindFlags(cmd)
		return f.createOptionsFromConfig(cmd)
	}

	opts, err := f.setProfileOptions(cmd, options.NewOptions())
	if err != nil {
		return nil, err
	}

	return f.setLTSVOptions(cmd, opts)
}

// alp ltsv diff
func (f *flags) createLTSVDiffOptions(cmd *cobra.Command) (*options.Options, error) {
	if f.config != "" {
		f.bindFlags(cmd)
		return f.createOptionsFromConfig(cmd)
	}

	opts, err := f.setDiffSubCommandOptions(cmd, options.NewOptions())
	if err != nil {
		return nil, err
	}

	return f.setLTSVOptions(cmd, opts)
}

// alp ltsv topN
func (f *flags) createLTSVTopNOptions(cmd *cobra.Command) (*options.Options, error) {
	if f.config != "" {
		f.bindFlags(cmd)
		return f.createOptionsFromConfig(cmd)
	}

	opts, err := f.setTopNSubCommandOptions(cmd, options.NewOptions())
	if err != nil {
		return nil, err
	}

	return f.setLTSVOptions(cmd, opts)
}

// alp ltsv count
func (f *flags) createLTSVCountOptions(cmd *cobra.Command) (*options.Options, error) {
	if f.config != "" {
		f.bindFlags(cmd)
		return f.createOptionsFromConfig(cmd)
	}

	opts, err := f.setCountSubCommandOptions(cmd, options.NewOptions())
	if err != nil {
		return nil, err
	}

	return f.setLTSVOptions(cmd, opts)
}

// alp regexp
func (f *flags) createRegexpOptions(cmd *cobra.Command) (*options.Options, error) {
	if f.config != "" {
		f.bindFlags(cmd)
		return f.createOptionsFromConfig(cmd)
	}

	opts, err := f.setProfileOptions(cmd, options.NewOptions())
	if err != nil {
		return nil, err
	}

	return f.setRegexpOptions(cmd, opts)
}

// alp regexp diff
func (f *flags) createRegexpDiffOptions(cmd *cobra.Command) (*options.Options, error) {
	if f.config != "" {
		f.bindFlags(cmd)
		return f.createOptionsFromConfig(cmd)
	}

	opts, err := f.setDiffSubCommandOptions(cmd, options.NewOptions())
	if err != nil {
		return nil, err
	}

	return f.setRegexpOptions(cmd, opts)
}

// alp regexp topN
func (f *flags) createRegexpTopNOptions(cmd *cobra.Command) (*options.Options, error) {
	if f.config != "" {
		f.bindFlags(cmd)
		return f.createOptionsFromConfig(cmd)
	}

	opts, err := f.setTopNSubCommandOptions(cmd, options.NewOptions())
	if err != nil {
		return nil, err
	}

	return f.setRegexpOptions(cmd, opts)
}

// alp regexp count
func (f *flags) createRegexpCountOptions(cmd *cobra.Command) (*options.Options, error) {
	if f.config != "" {
		f.bindFlags(cmd)
		return f.createOptionsFromConfig(cmd)
	}

	opts, err := f.setCountSubCommandOptions(cmd, options.NewOptions())
	if err != nil {
		return nil, err
	}

	return f.setRegexpOptions(cmd, opts)
}

// alp pcap
func (f *flags) createPcapOptions(cmd *cobra.Command) (*options.Options, error) {
	if f.config != "" {
		f.bindFlags(cmd)
		return f.createOptionsFromConfig(cmd)
	}

	opts, err := f.setProfileOptions(cmd, options.NewOptions())
	if err != nil {
		return nil, err
	}

	return f.setPcapOptions(cmd, opts)
}

// alp pcap diff
func (f *flags) createPcapDiffOptions(cmd *cobra.Command) (*options.Options, error) {
	if f.config != "" {
		f.bindFlags(cmd)
		return f.createOptionsFromConfig(cmd)
	}

	opts, err := f.setDiffSubCommandOptions(cmd, options.NewOptions())
	if err != nil {
		return nil, err
	}

	return f.setPcapOptions(cmd, opts)
}

// alp pcap topN
func (f *flags) createPcapTopNOptions(cmd *cobra.Command) (*options.Options, error) {
	if f.config != "" {
		f.bindFlags(cmd)
		return f.createOptionsFromConfig(cmd)
	}

	opts, err := f.setTopNSubCommandOptions(cmd, options.NewOptions())
	if err != nil {
		return nil, err
	}

	return f.setPcapOptions(cmd, opts)
}

// alp diff
func (f *flags) createDiffOptions(cmd *cobra.Command) (*options.Options, error) {
	if f.config != "" {
		f.bindFlags(cmd)
		return f.createOptionsFromConfig(cmd)
	}

	return f.setDiffOptions(cmd, options.NewOptions())
}
