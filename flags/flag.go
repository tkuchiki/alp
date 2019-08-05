package flags

import (
	"fmt"

	"github.com/tkuchiki/alp/options"

	"github.com/tkuchiki/alp/stats"
	"gopkg.in/alecthomas/kingpin.v2"
)

type GlobalFlags struct {
	Config         string
	File           string
	Dump           string
	Load           string
	Sort           string
	Reverse        bool
	QueryString    bool
	Tsv            bool
	NoHeaders      bool
	Limit          int
	Location       string
	Output         string
	MatchingGroups string
	Filters        string
}

var SortKeys = []string{
	"max",
	"min",
	"avg",
	"sum",
	"count",
	"uri",
	"method",
	"max-body",
	"min-body",
	"avg-body",
	"sum-body",
	"p1",
	"p50",
	"p99",
	"stddev",
}

var SortOptions = map[string]string{
	"max":      stats.SortMaxResponseTime,
	"min":      stats.SortMinResponseTime,
	"avg":      stats.SortAvgResponseTime,
	"sum":      stats.SortSumResponseTime,
	"count":    stats.SortCount,
	"uri":      stats.SortUri,
	"method":   stats.SortMethod,
	"max-body": stats.SortMaxResponseBodyBytes,
	"min-body": stats.SortMinResponseBodyBytes,
	"avg-body": stats.SortAvgResponseBodyBytes,
	"sum-body": stats.SortSumResponseBodyBytes,
	"p1":       stats.SortP1ResponseTime,
	"p50":      stats.SortP50ResponseTime,
	"p99":      stats.SortP99ResponseTime,
	"stddev":   stats.SortStddevResponseTime,
}

func NewGlobalFlags() *GlobalFlags {
	return &GlobalFlags{}
}

func (f *GlobalFlags) InitGlobalFlags(app *kingpin.Application) {
	app.Flag("config", "config file").
		Short('c').StringVar(&f.Config)
	app.Flag("file", "access log file").
		Short('f').StringVar(&f.File)
	app.Flag("dump", "dump profile data").
		Short('d').StringVar(&f.Dump)
	app.Flag("load", "load profile data").
		Short('l').StringVar(&f.Load)
	app.Flag("sort", "sort").
		Default(options.DefaultSortOption).EnumVar(&f.Sort, SortKeys...)
	app.Flag("reverse", "reverse the result of comparisons").
		Short('r').BoolVar(&f.Reverse)
	app.Flag("query-string", "include query string").
		Short('q').BoolVar(&f.QueryString)
	app.Flag("tsv", "tsv format (default: table)").
		BoolVar(&f.Tsv)
	app.Flag("limit", "set an upper limit of the target uri").
		Default(fmt.Sprint(options.DefaultLimitOption)).IntVar(&f.Limit)
	app.Flag("location", "location name").
		Default(options.DefaultLocationOption).StringVar(&f.Location)
	app.Flag("output", "output").Short('o').
		Default(options.DefaultOutputOption).StringVar(&f.Output)
	app.Flag("noheaders", "print no header line at all (only --tsv)").
		BoolVar(&f.NoHeaders)
	app.Flag("matching-groups", "uri matching groups (comma separated)").
		Short('m').PlaceHolder("PATTERN,...").StringVar(&f.MatchingGroups)
	app.Flag("filters", "the filters").
		StringVar(&f.Filters)
}
