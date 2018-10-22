package flag

import (
	"fmt"

	"gopkg.in/alecthomas/kingpin.v2"
)

type Flags struct {
	Config            string
	File              string
	Dump              string
	Load              string
	Max               bool
	Min               bool
	Avg               bool
	Sum               bool
	Cnt               bool
	SortUri           bool
	Method            bool
	MaxBody           bool
	MinBody           bool
	AvgBody           bool
	SumBody           bool
	P1                bool
	P50               bool
	P99               bool
	Stddev            bool
	Reverse           bool
	QueryString       bool
	Tsv               bool
	NoHeaders         bool
	ApptimeLabel      string
	ReqtimeLabel      string
	StatusLabel       string
	SizeLabel         string
	MethodLabel       string
	UriLabel          string
	TimeLabel         string
	Limit             int
	Location          string
	Includes          string
	Excludes          string
	IncludeStatuses   string
	ExcludeStatuses   string
	Aggregates        string
	StartTime         string
	EndTime           string
	StartTimeDuration string
	EndTimeDuration   string
}

func NewFlags() *Flags {
	return &Flags{}
}

func (f *Flags) InitFlags(app *kingpin.Application) {
	app.Flag("config", "config file").Short('c').StringVar(&f.Config)
	app.Flag("file", "access log file").Short('f').StringVar(&f.File)
	app.Flag("dump", "dump profile data").Short('d').StringVar(&f.Dump)
	app.Flag("load", "load profile data").Short('l').StringVar(&f.Load)
	app.Flag("max", "sort by max response time").BoolVar(&f.Max)
	app.Flag("min", "sort by min response time").BoolVar(&f.Min)
	app.Flag("avg", "sort by avg response time").BoolVar(&f.Avg)
	app.Flag("sum", "sort by sum response time").BoolVar(&f.Sum)
	app.Flag("cnt", "sort by count").BoolVar(&f.Cnt)
	app.Flag("uri", "sort by uri").BoolVar(&f.SortUri)
	app.Flag("method", "sort by method").BoolVar(&f.Method)
	app.Flag("max-body", "sort by max body size").BoolVar(&f.MaxBody)
	app.Flag("min-body", "sort by min body size").BoolVar(&f.MinBody)
	app.Flag("avg-body", "sort by avg body size").BoolVar(&f.AvgBody)
	app.Flag("sum-body", "sort by sum body size").BoolVar(&f.SumBody)
	app.Flag("p1", "sort by 1 percentail response time").BoolVar(&f.P1)
	app.Flag("p50", "sort by 50 percentail response time").BoolVar(&f.P50)
	app.Flag("p99", "sort by 99 percentail response time").BoolVar(&f.P99)
	app.Flag("stddev", "sort by standard deviation response time").BoolVar(&f.Stddev)
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
	app.Flag("include-statuses", "don't exclude status code matching PATTERN (comma separated)").PlaceHolder("PATTERN,...").StringVar(&f.IncludeStatuses)
	app.Flag("exclude-statuses", "exclude uri status code PATTERN (comma separated)").PlaceHolder("PATTERN,...").StringVar(&f.ExcludeStatuses)
	app.Flag("noheaders", "print no header line at all (only --tsv)").BoolVar(&f.NoHeaders)
	app.Flag("aggregates", "aggregate uri matching PATTERN (comma separated)").PlaceHolder("PATTERN,...").StringVar(&f.Aggregates)
	app.Flag("start-time", "since the start time").PlaceHolder("TIME").StringVar(&f.StartTime)
	app.Flag("end-time", "end time earlier").PlaceHolder("TIME").StringVar(&f.EndTime)
	app.Flag("start-time-duration", "since the start time (now - time.Duration)").PlaceHolder("TIME_DURATION").StringVar(&f.StartTimeDuration)
	app.Flag("end-time-duration", "end time earlier (now - time.Duration)").PlaceHolder("TIME_DURATION").StringVar(&f.EndTimeDuration)
}
