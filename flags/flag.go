package flags

import (
	"fmt"

	"github.com/tkuchiki/alp/options"

	"gopkg.in/alecthomas/kingpin.v2"
)

type GlobalFlags struct {
	Config                  string
	File                    string
	Dump                    string
	Load                    string
	Sort                    string
	Reverse                 bool
	QueryString             bool
	QueryStringIgnoreValues bool
	DecodeUri               bool
	Format                  string
	NoHeaders               bool
	ShowFooters             bool
	Limit                   int
	Location                string
	Output                  string
	MatchingGroups          string
	Filters                 string
	PosFile                 string
	NoSavePos               bool
	Percentiles             string
	PaginationLimit         int
}

var Formats = []string{
	"table",
	"md",
	"markdown",
	"tsv",
	"csv",
	"html",
}

func NewGlobalFlags() *GlobalFlags {
	return &GlobalFlags{}
}

func (f *GlobalFlags) InitGlobalFlags(app *kingpin.Application) {
	app.Flag("config", "The configuration file").
		Short('c').StringVar(&f.Config)
	app.Flag("file", "The access log file").
		StringVar(&f.File)
	app.Flag("dump", "Dump profiled data as YAML").
		Short('d').StringVar(&f.Dump)
	app.Flag("load", "Load the profiled YAML data").
		Short('l').StringVar(&f.Load)
	app.Flag("sort", "Output the results in sorted order").
		PlaceHolder(options.DefaultSortOption).Default(options.DefaultSortOption).StringVar(&f.Sort)
	app.Flag("reverse", "Sort results in reverse order").
		Short('r').BoolVar(&f.Reverse)
	app.Flag("query-string", "Include the URI query string").
		Short('q').BoolVar(&f.QueryString)
	app.Flag("qs-ignore-values", "Ignore the value of the query string. Replace all values with xxx (only use with -q)").
		BoolVar(&f.QueryStringIgnoreValues)
	app.Flag("decode-uri", "Decode the URI").
		BoolVar(&f.DecodeUri)
	app.Flag("format", "The output format (table, markdown, tsv, csv and html)").
		PlaceHolder(options.DefaultFormatOption).EnumVar(&f.Format, Formats...)
	app.Flag("noheaders", "Output no header line at all (only --format=tsv, csv)").
		BoolVar(&f.NoHeaders)
	app.Flag("show-footers", "Output footer line at all (only --format=table, markdown)").
		BoolVar(&f.ShowFooters)
	app.Flag("limit", "The maximum number of results to display").
		PlaceHolder(fmt.Sprint(options.DefaultLimitOption)).IntVar(&f.Limit)
	app.Flag("location", "Location name for the timezone").
		PlaceHolder(options.DefaultLocationOption).StringVar(&f.Location)
	app.Flag("output", "Specifies the results to display, separated by commas").Short('o').
		PlaceHolder(options.DefaultOutputOption).StringVar(&f.Output)
	app.Flag("matching-groups", "Specifies URI matching groups separated by commas").
		Short('m').PlaceHolder("PATTERN,...").StringVar(&f.MatchingGroups)
	app.Flag("filters", "Only the logs are profiled that match the conditions").
		Short('f').StringVar(&f.Filters)
	app.Flag("pos", "The position file").
		PlaceHolder("POSITION_FILE").StringVar(&f.PosFile)
	app.Flag("nosave-pos", "Do not save position file").
		BoolVar(&f.NoSavePos)
	app.Flag("percentiles", "Specifies the percentiles separated by commas").Default("90,95,99").
		StringVar(&f.Percentiles)
	app.Flag("page", "Number of pages of pagination").
		PlaceHolder(fmt.Sprint(options.DefaultLimitOption)).IntVar(&f.PaginationLimit)
}
