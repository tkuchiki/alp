package flag

import (
	"fmt"

	"github.com/tkuchiki/alp/stats"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Flags struct {
	Config       string
	File         string
	Dump         string
	Load         string
	Sort         string
	Reverse      bool
	QueryString  bool
	Tsv          bool
	NoHeaders    bool
	ApptimeLabel string
	ReqtimeLabel string
	StatusLabel  string
	SizeLabel    string
	MethodLabel  string
	UriLabel     string
	TimeLabel    string
	Limit        int
	Location     string
	Includes     string
	Excludes     string
	Groups       string
	Filters      string
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
	"max-body": stats.SortMaxResponseBodySize,
	"min-body": stats.SortMinResponseBodySize,
	"avg-body": stats.SortAvgResponseBodySize,
	"sum-body": stats.SortSumResponseBodySize,
	"p1":       stats.SortP1ResponseTime,
	"p50":      stats.SortP50ResponseTime,
	"p99":      stats.SortP99ResponseTime,
	"stddev":   stats.SortStddevResponseTime,
}

func NewFlags() *Flags {
	return &Flags{}
}

func (f *Flags) InitFlags(app *kingpin.Application) {
	app.Flag("config", "config file").Short('c').StringVar(&f.Config)
	app.Flag("file", "access log file").Short('f').StringVar(&f.File)
	app.Flag("dump", "dump profile data").Short('d').StringVar(&f.Dump)
	app.Flag("load", "load profile data").Short('l').StringVar(&f.Load)
	app.Flag("sort", "sort").Default("max").EnumVar(&f.Sort, SortKeys...)
	app.Flag("reverse", "reverse the result of comparisons").Short('r').BoolVar(&f.Reverse)
	app.Flag("query-string", "include query string").Short('q').BoolVar(&f.QueryString)
	app.Flag("tsv", "tsv format (default: table)").BoolVar(&f.Tsv)
	app.Flag("apptime-label", "apptime label").Default("apptime").StringVar(&f.ApptimeLabel)
	app.Flag("reqtime-label", "reqtime label").Default("reqtime").StringVar(&f.ReqtimeLabel)
	app.Flag("status-label", "status label").Default("status").StringVar(&f.StatusLabel)
	app.Flag("size-label", "size label").Default("size").StringVar(&f.SizeLabel)
	app.Flag("method-label", "method label").Default("method").StringVar(&f.MethodLabel)
	app.Flag("uri-label", "uri label").Default("uri").StringVar(&f.UriLabel)
	app.Flag("time-label", "time label").Default("time").StringVar(&f.TimeLabel)
	app.Flag("limit", "set an upper limit of the target uri").Default(fmt.Sprint(5000)).IntVar(&f.Limit)
	app.Flag("location", "location name").StringVar(&f.Location)
	app.Flag("includes", "don't exclude uri matching PATTERN (comma separated)").PlaceHolder("PATTERN,...").StringVar(&f.Includes)
	app.Flag("excludes", "exclude uri matching PATTERN (comma separated)").PlaceHolder("PATTERN,...").StringVar(&f.Excludes)
	app.Flag("noheaders", "print no header line at all (only --tsv)").BoolVar(&f.NoHeaders)
	app.Flag("groups", "uri matching groups (comma separated)").PlaceHolder("PATTERN,...").StringVar(&f.Groups)
	app.Flag("filters", "the filters").StringVar(&f.Filters)
}