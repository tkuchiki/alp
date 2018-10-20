package flag

import (
	"fmt"

	"gopkg.in/alecthomas/kingpin.v2"
)

type Flags struct {
	Config            *string
	File              *string
	Dump              *string
	Load              *string
	Max               *bool
	Min               *bool
	Avg               *bool
	Sum               *bool
	Cnt               *bool
	SortUri           *bool
	Method            *bool
	MaxBody           *bool
	MinBody           *bool
	AvgBody           *bool
	SumBody           *bool
	P1                *bool
	P50               *bool
	P99               *bool
	Stddev            *bool
	Reverse           *bool
	QueryString       *bool
	Tsv               *bool
	NoHeaders         *bool
	ApptimeLabel      *string
	ReqtimeLabel      *string
	StatusLabel       *string
	SizeLabel         *string
	MethodLabel       *string
	UriLabel          *string
	TimeLabel         *string
	Limit             *int
	Location          *string
	Includes          *string
	Excludes          *string
	IncludeStatuses   *string
	ExcludeStatuses   *string
	Aggregates        *string
	StartTime         *string
	EndTime           *string
	StartTimeDuration *string
	EndTimeDuration   *string
}

func NewFlags(app *kingpin.Application) *Flags {
	flags := &Flags{}
	flags.Config = app.Flag("config", "config file").Short('c').String()
	flags.File = app.Flag("file", "access log file").Short('f').String()
	flags.Dump = app.Flag("dump", "dump profile data").Short('d').String()
	flags.Load = app.Flag("load", "load profile data").Short('l').String()
	flags.Max = app.Flag("max", "sort by max response time").Bool()
	flags.Min = app.Flag("min", "sort by min response time").Bool()
	flags.Avg = app.Flag("avg", "sort by avg response time").Bool()
	flags.Sum = app.Flag("sum", "sort by sum response time").Bool()
	flags.Cnt = app.Flag("cnt", "sort by count").Bool()
	flags.SortUri = app.Flag("uri", "sort by uri").Bool()
	flags.Method = app.Flag("method", "sort by method").Bool()
	flags.MaxBody = app.Flag("max-body", "sort by max body size").Bool()
	flags.MinBody = app.Flag("min-body", "sort by min body size").Bool()
	flags.AvgBody = app.Flag("avg-body", "sort by avg body size").Bool()
	flags.SumBody = app.Flag("sum-body", "sort by sum body size").Bool()
	flags.P1 = app.Flag("p1", "sort by 1 percentail response time").Bool()
	flags.P50 = app.Flag("p50", "sort by 50 percentail response time").Bool()
	flags.P99 = app.Flag("p99", "sort by 99 percentail response time").Bool()
	flags.Stddev = app.Flag("stddev", "sort by standard deviation response time").Bool()
	flags.Reverse = app.Flag("reverse", "reverse the result of comparisons").Short('r').Bool()
	flags.QueryString = app.Flag("query-string", "include query string").Short('q').Bool()
	flags.Tsv = app.Flag("tsv", "tsv format (default: table)").Bool()
	flags.ApptimeLabel = app.Flag("apptime-label", "apptime label").Default("apptime").String()
	flags.ReqtimeLabel = app.Flag("reqtime-label", "reqtime label").Default("reqtime").String()
	flags.StatusLabel = app.Flag("status-label", "status label").Default("status").String()
	flags.SizeLabel = app.Flag("size-label", "size label").Default("size").String()
	flags.MethodLabel = app.Flag("method-label", "method label").Default("method").String()
	flags.UriLabel = app.Flag("uri-label", "uri label").Default("uri").String()
	flags.TimeLabel = app.Flag("time-label", "time label").Default("time").String()
	flags.Limit = app.Flag("limit", "set an upper limit of the target uri").Default(fmt.Sprint(5000)).Int()
	flags.Location = app.Flag("location", "location name").String()
	flags.Includes = app.Flag("includes", "don't exclude uri matching PATTERN (comma separated)").PlaceHolder("PATTERN,...").String()
	flags.Excludes = app.Flag("excludes", "exclude uri matching PATTERN (comma separated)").PlaceHolder("PATTERN,...").String()
	flags.IncludeStatuses = app.Flag("include-statuses", "don't exclude status code matching PATTERN (comma separated)").PlaceHolder("PATTERN,...").String()
	flags.ExcludeStatuses = app.Flag("exclude-statuses", "exclude uri status code PATTERN (comma separated)").PlaceHolder("PATTERN,...").String()
	flags.NoHeaders = app.Flag("noheaders", "print no header line at all (only --tsv)").Bool()
	flags.Aggregates = app.Flag("aggregates", "aggregate uri matching PATTERN (comma separated)").PlaceHolder("PATTERN,...").String()
	flags.StartTime = app.Flag("start-time", "since the start time").PlaceHolder("TIME").String()
	flags.EndTime = app.Flag("end-time", "end time earlier").PlaceHolder("TIME").String()
	flags.StartTimeDuration = app.Flag("start-time-duration", "since the start time (now - time.Duration)").PlaceHolder("TIME_DURATION").String()
	flags.EndTimeDuration = app.Flag("end-time-duration", "end time earlier (now - time.Duration)").PlaceHolder("TIME_DURATION").String()

	return flags
}

func NewEmptyFlags() *Flags {
	var emptyStr string
	var emptyBool bool
	var emptyInt int

	return &Flags{
		Config:            &emptyStr,
		File:              &emptyStr,
		Dump:              &emptyStr,
		Load:              &emptyStr,
		Max:               &emptyBool,
		Min:               &emptyBool,
		Avg:               &emptyBool,
		Sum:               &emptyBool,
		Cnt:               &emptyBool,
		SortUri:           &emptyBool,
		Method:            &emptyBool,
		MaxBody:           &emptyBool,
		MinBody:           &emptyBool,
		AvgBody:           &emptyBool,
		SumBody:           &emptyBool,
		P1:                &emptyBool,
		P50:               &emptyBool,
		P99:               &emptyBool,
		Stddev:            &emptyBool,
		Reverse:           &emptyBool,
		QueryString:       &emptyBool,
		Tsv:               &emptyBool,
		NoHeaders:         &emptyBool,
		ApptimeLabel:      &emptyStr,
		ReqtimeLabel:      &emptyStr,
		StatusLabel:       &emptyStr,
		SizeLabel:         &emptyStr,
		MethodLabel:       &emptyStr,
		UriLabel:          &emptyStr,
		TimeLabel:         &emptyStr,
		Limit:             &emptyInt,
		Location:          &emptyStr,
		Includes:          &emptyStr,
		Excludes:          &emptyStr,
		IncludeStatuses:   &emptyStr,
		ExcludeStatuses:   &emptyStr,
		Aggregates:        &emptyStr,
		StartTime:         &emptyStr,
		EndTime:           &emptyStr,
		StartTimeDuration: &emptyStr,
		EndTimeDuration:   &emptyStr,
	}
}
