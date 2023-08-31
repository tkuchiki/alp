package cmd

import (
	"github.com/tkuchiki/alp/helpers"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tkuchiki/alp/options"
	"github.com/tkuchiki/alp/stats"
)

func defineOptions(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP("config", "", "", "The configuration file")
	cmd.PersistentFlags().StringP("file", "", "", "The access log file")
	cmd.PersistentFlags().StringP("dump", "", "", "Dump profiled data as YAML")
	cmd.PersistentFlags().StringP("load", "", "", "Load the profiled YAML data")
	cmd.PersistentFlags().StringP("format", "", options.DefaultFormatOption, "The output format (table, markdown, tsv, csv and html)")
	cmd.PersistentFlags().StringP("sort", "", options.DefaultSortOption, "Output the results in sorted order")
	cmd.PersistentFlags().BoolP("reverse", "r", false, "Sort results in reverse order")
	cmd.PersistentFlags().BoolP("noheaders", "", false, "Output no header line at all (only --format=tsv, csv)")
	cmd.PersistentFlags().BoolP("show-footers", "", false, "Output footer line at all (only --format=table, markdown)")
	cmd.PersistentFlags().IntP("limit", "", options.DefaultLimitOption, "The maximum number of results to display")
	cmd.PersistentFlags().StringP("output", "o", options.DefaultOutputOption, "Specifies the results to display, separated by commas")
	cmd.PersistentFlags().BoolP("query-string", "q", false, "Include the URI query string")
	cmd.PersistentFlags().BoolP("qs-ignore-values", "", false, "Ignore the value of the query string. Replace all values with xxx (only use with -q)")
	cmd.PersistentFlags().StringP("location", "", options.DefaultLocationOption, "Location name for the timezone")
	cmd.PersistentFlags().BoolP("decode-uri", "", false, "Decode the URI")
	cmd.PersistentFlags().StringP("matching-groups", "m", "", "Specifies Query matching groups separated by commas")
	cmd.PersistentFlags().StringP("filters", "f", "", "Only the logs are profiled that match the conditions")
	cmd.PersistentFlags().StringP("pos", "", "", "The position file")
	cmd.PersistentFlags().BoolP("nosave-pos", "", false, "Do not save position file")
	cmd.PersistentFlags().StringP("percentiles", "", "", "Specifies the percentiles separated by commas")
	cmd.PersistentFlags().IntP("page", "", options.DefaultPaginationLimit, "Number of pages of pagination")

}

func createCommonOptionsFromFlags(cmd *cobra.Command, sortOptions *stats.SortOptions) (*options.Options, error) {
	file, err := cmd.PersistentFlags().GetString("file")
	if err != nil {
		return nil, err
	}

	dump, err := cmd.PersistentFlags().GetString("dump")
	if err != nil {
		return nil, err
	}

	load, err := cmd.PersistentFlags().GetString("load")
	if err != nil {
		return nil, err
	}

	format, err := cmd.PersistentFlags().GetString("format")
	if err != nil {
		return nil, err
	}

	sort, err := cmd.PersistentFlags().GetString("sort")
	if err != nil {
		return nil, err
	}

	err = sortOptions.SetAndValidate(sort)
	if err != nil {
		return nil, err
	}

	reverse, err := cmd.PersistentFlags().GetBool("reverse")
	if err != nil {
		return nil, err
	}

	noHeaders, err := cmd.PersistentFlags().GetBool("noheaders")
	if err != nil {
		return nil, err
	}

	showFooters, err := cmd.PersistentFlags().GetBool("show-footers")
	if err != nil {
		return nil, err
	}

	limit, err := cmd.PersistentFlags().GetInt("limit")
	if err != nil {
		return nil, err
	}

	output, err := cmd.PersistentFlags().GetString("output")
	if err != nil {
		return nil, err
	}

	queryString, err := cmd.PersistentFlags().GetBool("query-string")
	if err != nil {
		return nil, err
	}

	queryStringIgnoreValues, err := cmd.PersistentFlags().GetBool("qs-ignore-values")
	if err != nil {
		return nil, err
	}

	decodeUri, err := cmd.PersistentFlags().GetBool("decode-uri")
	if err != nil {
		return nil, err
	}

	location, err := cmd.PersistentFlags().GetString("location")
	if err != nil {
		return nil, err
	}

	matchingGroups, err := cmd.PersistentFlags().GetString("matching-groups")
	if err != nil {
		return nil, err
	}

	filters, err := cmd.PersistentFlags().GetString("filters")
	if err != nil {
		return nil, err
	}

	pos, err := cmd.PersistentFlags().GetString("pos")
	if err != nil {
		return nil, err
	}

	noSavePos, err := cmd.PersistentFlags().GetBool("nosave-pos")
	if err != nil {
		return nil, err
	}

	ps, err := cmd.PersistentFlags().GetString("percentiles")
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

	paginationLimit, err := cmd.PersistentFlags().GetInt("page")
	if err != nil {
		return nil, err
	}

	opts := options.NewOptions()

	return options.SetOptions(opts,
		options.File(file),
		options.Dump(dump),
		options.Load(load),
		options.Sort(sortOptions.SortType()),
		options.Reverse(reverse),
		options.Format(format),
		options.Limit(limit),
		options.Output(output),
		options.QueryString(queryString),
		options.QueryStringIgnoreValues(queryStringIgnoreValues),
		options.DecodeUri(decodeUri),
		options.Location(location),
		options.NoHeaders(noHeaders),
		options.ShowFooters(showFooters),
		options.CSVGroups(matchingGroups),
		options.Filters(filters),
		options.PosFile(pos),
		options.NoSavePos(noSavePos),
		options.Percentiles(percentiles),
		options.PaginationLimit(paginationLimit),
	), nil
}

func createOptionsFromConfig(cmd *cobra.Command, sortOptions *stats.SortOptions, config string) (*options.Options, error) {
	opts := options.NewOptions()
	viper.SetConfigFile(config)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(opts); err != nil {
		return nil, err
	}

	if err := sortOptions.SetAndValidate(opts.Sort); err != nil {
		return nil, err
	}
	opts.Sort = sortOptions.SortType()

	return opts, nil
}

func bindCommonFlags(cmd *cobra.Command) {
	viper.BindPFlag("file", cmd.PersistentFlags().Lookup("file"))
	viper.BindPFlag("dump", cmd.PersistentFlags().Lookup("dump"))
	viper.BindPFlag("load", cmd.PersistentFlags().Lookup("load"))
	viper.BindPFlag("sort", cmd.PersistentFlags().Lookup("sort"))
	viper.BindPFlag("reverse", cmd.PersistentFlags().Lookup("reverse"))
	viper.BindPFlag("query_string", cmd.PersistentFlags().Lookup("query-string"))
	viper.BindPFlag("query_string_ignore_values", cmd.PersistentFlags().Lookup("qs-ignore-values"))
	viper.BindPFlag("decode_uri", cmd.PersistentFlags().Lookup("decode-uri"))
	viper.BindPFlag("format", cmd.PersistentFlags().Lookup("format"))
	viper.BindPFlag("noheaders", cmd.PersistentFlags().Lookup("noheaders"))
	viper.BindPFlag("show_footers", cmd.PersistentFlags().Lookup("show-footers"))
	viper.BindPFlag("limit", cmd.PersistentFlags().Lookup("limit"))
	viper.BindPFlag("matching_groups", cmd.PersistentFlags().Lookup("matching-groups"))
	viper.BindPFlag("filters", cmd.PersistentFlags().Lookup("filters"))
	viper.BindPFlag("pos_file", cmd.PersistentFlags().Lookup("pos"))
	viper.BindPFlag("nosave_pos", cmd.PersistentFlags().Lookup("nosave-pos"))
	viper.BindPFlag("location", cmd.PersistentFlags().Lookup("location"))
	viper.BindPFlag("output", cmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("percentiles", cmd.PersistentFlags().Lookup("percentiles"))
	viper.BindPFlag("pagenation_limit", cmd.PersistentFlags().Lookup("page"))
}
