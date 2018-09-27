package alp

import (
	"fmt"
	"io"
	"net/url"
	"regexp"

	"github.com/tkuchiki/http-profile-helper"
	"gopkg.in/alecthomas/kingpin.v2"
)

func Run() error {
	var (
		config = kingpin.Flag("config", "config file").Short('c').String()
		file   = kingpin.Flag("file", "access log file").Short('f').String()
		dump              = kingpin.Flag("dump", "dump profile data").Short('d').String()
		load              = kingpin.Flag("load", "load profile data").Short('l').String()
		max               = kingpin.Flag("max", "sort by max response time").Bool()
		min               = kingpin.Flag("min", "sort by min response time").Bool()
		avg               = kingpin.Flag("avg", "sort by avg response time").Bool()
		sum               = kingpin.Flag("sum", "sort by sum response time").Bool()
		cnt               = kingpin.Flag("cnt", "sort by count").Bool()
		sortUri           = kingpin.Flag("uri", "sort by uri").Bool()
		method            = kingpin.Flag("method", "sort by method").Bool()
		maxBody           = kingpin.Flag("max-body", "sort by max body size").Bool()
		minBody           = kingpin.Flag("min-body", "sort by min body size").Bool()
		avgBody           = kingpin.Flag("avg-body", "sort by avg body size").Bool()
		sumBody           = kingpin.Flag("sum-body", "sort by sum body size").Bool()
		p1                = kingpin.Flag("p1", "sort by 1 percentail response time").Bool()
		p50               = kingpin.Flag("p50", "sort by 50 percentail response time").Bool()
		p99               = kingpin.Flag("p99", "sort by 99 percentail response time").Bool()
		stddev            = kingpin.Flag("stddev", "sort by standard deviation response time").Bool()
		reverse           = kingpin.Flag("reverse", "reverse the result of comparisons").Short('r').Bool()
		queryString       = kingpin.Flag("query-string", "include query string").Short('q').Bool()
		tsv               = kingpin.Flag("tsv", "tsv format (default: table)").Bool()
		apptimeLabel      = kingpin.Flag("apptime-label", "apptime label").Default("apptime").String()
		reqtimeLabel      = kingpin.Flag("reqtime-label", "reqtime label").Default("reqtime").String()
		statusLabel       = kingpin.Flag("status-label", "status label").Default("status").String()
		sizeLabel         = kingpin.Flag("size-label", "size label").Default("size").String()
		methodLabel       = kingpin.Flag("method-label", "method label").Default("method").String()
		uriLabel          = kingpin.Flag("uri-label", "uri label").Default("uri").String()
		timeLabel         = kingpin.Flag("time-label", "time label").Default("time").String()
		limit             = kingpin.Flag("limit", "set an upper limit of the target uri").Default(fmt.Sprint(5000)).Int()
		location          = kingpin.Flag("location", "location name").String()
		includes          = kingpin.Flag("includes", "don't exclude uri matching PATTERN (comma separated)").PlaceHolder("PATTERN,...").String()
		excludes          = kingpin.Flag("excludes", "exclude uri matching PATTERN (comma separated)").PlaceHolder("PATTERN,...").String()
		includeStatuses   = kingpin.Flag("include-statuses", "don't exclude status code matching PATTERN (comma separated)").PlaceHolder("PATTERN,...").String()
		excludeStatuses   = kingpin.Flag("exclude-statuses", "exclude uri status code PATTERN (comma separated)").PlaceHolder("PATTERN,...").String()
		noHeaders         = kingpin.Flag("noheaders", "print no header line at all (only --tsv)").Bool()
		aggregates        = kingpin.Flag("aggregates", "aggregate uri matching PATTERN (comma separated)").PlaceHolder("PATTERN,...").String()
		startTime         = kingpin.Flag("start-time", "since the start time").PlaceHolder("TIME").String()
		endTime           = kingpin.Flag("end-time", "end time earlier").PlaceHolder("TIME").String()
		startTimeDuration = kingpin.Flag("start-time-duration", "since the start time (now - time.Duration)").PlaceHolder("TIME_DURATION").String()
		endTimeDuration   = kingpin.Flag("end-time-duration", "end time earlier (now - time.Duration)").PlaceHolder("TIME_DURATION").String()

		err error
		sort string
		uri string
	)

	kingpin.CommandLine.Help = "Access Log Profiler for LTSV (read from file or stdin)."
	kingpin.Version("0.4.0")
	kingpin.Parse()

	if *max {
		sort = "max"
	} else if *min {
		sort = "min"
	} else if *avg {
		sort = "avg"
	} else if *sum {
		sort = "sum"
	} else if *cnt {
		sort = "cnt"
	} else if *sortUri {
		sort = "uri"
	} else if *method {
		sort = "method"
	} else if *maxBody {
		sort = "max-body"
	} else if *minBody {
		sort = "min-body"
	} else if *avgBody {
		sort = "avg-body"
	} else if *sumBody {
		sort = "sum-body"
	} else if *p1 {
		sort = "p1"
	} else if *p50 {
		sort = "p50"
	} else if *p99 {
		sort = "p99"
	} else if *stddev {
		sort = "stddev"
	} else {
		sort = "max"
	}

	cliOptions := helper.NewCliOptions(&helper.Options{
			File:               *file,
			Sort:               sort,
			Reverse:            *reverse,
			QueryString:        *queryString,
			Tsv:                *tsv,
			ApptimeLabel:       *apptimeLabel,
			ReqtimeLabel:       *reqtimeLabel,
			StatusLabel:        *statusLabel,
			SizeLabel:          *sizeLabel,
			MethodLabel:        *methodLabel,
			UriLabel:           *uriLabel,
			TimeLabel:          *timeLabel,
			Limit:              *limit,
			NoHeaders:          *noHeaders,
			StartTime:          *startTime,
			EndTime:            *endTime,
			StartTimeDuration:  *startTimeDuration,
			EndTimeDuration:    *endTimeDuration,
		},
		*includes, *excludes, *includeStatuses, *excludeStatuses, *aggregates)

	var fileOptions *helper.Options
	if *config != "" {
		fileOptions, err = helper.LoadOptionsFromYAML(*config)
		if err != nil {
			return err
		}
	} else {
		fileOptions = &helper.Options{}
	}

	defaultOptions := helper.NewDefaultOptions()
	options := helper.SetOptions(cliOptions, fileOptions, defaultOptions)

	prof := NewProfiler()

	if *load != "" {
		err = prof.stats.LoadStats(*load)
		if err != nil {
			return err
		}

		prof.stats.Print()
		return nil
	}

	f, err := prof.helper.OpenFD(options.File)
	if err != nil {
		return err
	}
	defer f.Close()

	var includeGroups []*regexp.Regexp

	if len(options.Includes) > 0 {
		includeGroups, err = prof.helper.CompileIncludeGroups(options.Includes)
		if err != nil {
			return err
		}
	}

	var excludeGroups []*regexp.Regexp
	if len(options.Excludes) > 0 {
		excludeGroups, err = prof.helper.CompileExcludeGroups(options.Excludes)
		if err != nil {
			return err
		}
	}

	var includeStatusGroups []*regexp.Regexp
	if len(options.IncludeStatuses) > 0 {
		includeGroups, err = prof.helper.CompileIncludeStatusGroups(options.IncludeStatuses)
		if err != nil {
			return err
		}
	}

	var excludeStatusGroups []*regexp.Regexp
	if len(options.ExcludeStatuses) > 0 {
		includeGroups, err = prof.helper.CompileExcludeStatusGroups(options.ExcludeStatuses)
		if err != nil {
			return err
		}
	}

	var uriGroups []*regexp.Regexp
	if len(options.Aggregates) > 0 {
		uriGroups, err = prof.helper.CompileUriGroups(options.Aggregates)
		if err != nil {
			return err
		}
	}

	err = prof.helper.InitParseTime(*location)
	if err != nil {
		return err
	}

	var sTimeNano int64
	if options.StartTime != "" {
		sTimeNano, err = prof.helper.TimeStrToUnixNano(options.StartTime)
		if err != nil {
			return err
		}
	}

	if options.StartTimeDuration != "" {
		sTimeNano, err = prof.helper.SubTimeDuration(options.StartTimeDuration)
		if err != nil {
			return err
		}
	}

	var eTimeNano int64
	if options.EndTime != "" {
		eTimeNano, err = prof.helper.TimeStrToUnixNano(options.EndTime)
		if err != nil {
			return err
		}
	}

	if options.EndTimeDuration != "" {
		eTimeNano, err = prof.helper.SubTimeDuration(options.EndTimeDuration)
		if err != nil {
			return err
		}
	}

	prof.helper.InitReader(f)
Loop:
	for {
		line, err := prof.helper.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		resTime, err := prof.helper.ParseResponseTime(line[options.ApptimeLabel])
		if err != nil {
			var reqTime float64
			reqTime, err = prof.helper.ParseRequestTime(line[options.ReqtimeLabel])
			if err != nil {
				continue
			}

			resTime = reqTime
		}

		bodySize, err := prof.helper.ParseBodySize(line[options.SizeLabel])
		if err != nil {
			continue
		}

		status, err := prof.helper.ParseStatusCode(line[options.StatusLabel])
		if err != nil {
			continue
		}

		if sTimeNano != 0 || eTimeNano != 0 {
			t, err := prof.helper.ParseTime(line[options.TimeLabel])
			if err != nil {
				continue
			}
			timeNano := t.UnixNano()
			if !prof.helper.IsIncludedInTime(sTimeNano, eTimeNano, timeNano) {
				continue
			}
		}

		u, err := url.Parse(line[options.UriLabel])
		if err != nil {
			continue
		}

		if options.QueryString {
			v := url.Values{}
			values := u.Query()
			for q := range values {
				v.Set(q, "xxx")
			}
			uri = fmt.Sprintf("%s?%s", u.Path, v.Encode())
		} else {
			uri = u.Path
		}

		if len(options.Includes) > 0 {
			isnotMatched := true
			for _, re := range includeGroups {
				if ok := re.Match([]byte(uri)); ok && err == nil {
					isnotMatched = false
				} else if err != nil {
					return err
				}
			}

			if isnotMatched {
				continue Loop
			}
		}

		if len(options.Excludes) > 0 {
			for _, re := range excludeGroups {
				if ok := re.Match([]byte(uri)); ok && err == nil {
					continue Loop
				} else if err != nil {
					return err
				}
			}
		}

		if len(options.IncludeStatuses) > 0 {
			isnotMatched := true
			for _, re := range includeStatusGroups {
				if ok := re.Match([]byte(line[options.StatusLabel])); ok && err == nil {
					isnotMatched = false
				} else if err != nil {
					return err
				}
			}

			if isnotMatched {
				continue Loop
			}
		}

		if len(options.ExcludeStatuses) > 0 {
			for _, re := range excludeStatusGroups {
				if ok := re.Match([]byte(line[options.StatusLabel])); ok && err == nil {
					continue Loop
				} else if err != nil {
					return err
				}
			}
		}

		isMatched := false
		if len(options.Aggregates) > 0 {
			for _, re := range uriGroups {
				if ok := re.Match([]byte(uri)); ok && err == nil {
					isMatched = true
					pattern := re.String()
					uri = pattern
					prof.stats.Set(uri, line[options.MethodLabel], status,
						resTime, bodySize, 0)
				} else if err != nil {
					return err
				}
			}
		}

		if !isMatched {
			prof.stats.Set(uri, line[options.MethodLabel], status,
				resTime, bodySize, 0)
		}

		if prof.stats.CountUris() > options.Limit {
			return fmt.Errorf("Too many URI's (%d or less)", options.Limit)
		}
	}

	if *dump != "" {
		err = prof.stats.DumpStats(*dump)
		if err != nil {
			return err
		}
	}

	prof.stats.Print()

	return nil
}
