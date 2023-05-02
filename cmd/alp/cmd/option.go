package cmd

import (
	"os"

	"github.com/tkuchiki/alp/helpers"

	"github.com/spf13/cobra"
	"github.com/tkuchiki/alp/options"
	"github.com/tkuchiki/alp/stats"
)

func defineOptions(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP("config", "", "", "The configuration file")
	cmd.PersistentFlags().StringP("file", "", "", "The slowlog file")
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

func createOptions(cmd *cobra.Command, sortOptions *stats.SortOptions) (*options.Options, error) {
	config, err := cmd.PersistentFlags().GetString("config")
	if err != nil {
		return nil, err
	}

	file, err := cmd.PersistentFlags().GetString("file")
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

	var opts *options.Options
	if config != "" {
		cf, err := os.Open(config)
		if err != nil {
			return nil, err
		}
		defer cf.Close()

		opts, err = options.LoadOptionsFromReader(cf)
		if err != nil {
			return nil, err
		}

		err = sortOptions.SetAndValidate(opts.Sort)
		if err != nil {
			return nil, err
		}

		percentiles = opts.Percentiles
	} else {
		opts = options.NewOptions()
	}

	return options.SetOptions(opts,
		options.File(file),
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
