package main

import (
	"fmt"
	"github.com/najeira/ltsv"
	"github.com/tkuchiki/parsetime"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"log"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
)

const (
	ApptimeLabel = "apptime"
	ReqtimeLabel = "reqtime"
	StatusLabel  = "status"
	SizeLabel    = "size"
	MethodLabel  = "method"
	UriLabel     = "uri"
	TimeLabel    = "time"
	Limit        = 5000
)

var (
	config            = kingpin.Flag("config", "config file").Short('c').String()
	file              = kingpin.Flag("file", "access log file").Short('f').String()
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
	apptimeLabel      = kingpin.Flag("apptime-label", "apptime label").Default(ApptimeLabel).String()
	reqtimeLabel      = kingpin.Flag("reqtime-label", "reqtime label").Default(ReqtimeLabel).String()
	statusLabel       = kingpin.Flag("status-label", "status label").Default(StatusLabel).String()
	sizeLabel         = kingpin.Flag("size-label", "size label").Default(SizeLabel).String()
	methodLabel       = kingpin.Flag("method-label", "method label").Default(MethodLabel).String()
	uriLabel          = kingpin.Flag("uri-label", "uri label").Default(UriLabel).String()
	timeLabel         = kingpin.Flag("time-label", "time label").Default(TimeLabel).String()
	limit             = kingpin.Flag("limit", "set an upper limit of the target uri").Default(fmt.Sprint(Limit)).Int()
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

	uri       string
	index     string
	accessLog Profiles
	uriHints      = make(map[string]int)
	length    int = 0
	cursor    int = 0
)

func main() {
	kingpin.CommandLine.Help = "Access Log Profiler for LTSV (read from file or stdin)."
	kingpin.Version("0.3.0")
	kingpin.Parse()

	var f *os.File
	var err error
	var c Config

	if *config != "" {
		c, err = LoadYAML(*config)
		if err != nil {
			log.Fatal(err)
		}
	}

	option := Config{
		File:               *file,
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
		IncludesStr:        *includes,
		ExcludesStr:        *excludes,
		IncludeStatusesStr: *includeStatuses,
		ExcludeStatusesStr: *excludeStatuses,
		NoHeaders:          *noHeaders,
		AggregatesStr:      *aggregates,
		StartTime:          *startTime,
		EndTime:            *endTime,
		StartTimeDuration:  *startTimeDuration,
		EndTimeDuration:    *endTimeDuration,
	}

	if *max {
		c.Sort = "max"
	} else if *min {
		c.Sort = "min"
	} else if *avg {
		c.Sort = "avg"
	} else if *sum {
		c.Sort = "sum"
	} else if *cnt {
		c.Sort = "cnt"
	} else if *sortUri {
		c.Sort = "uri"
	} else if *method {
		c.Sort = "method"
	} else if *maxBody {
		c.Sort = "max-body"
	} else if *minBody {
		c.Sort = "min-body"
	} else if *avgBody {
		c.Sort = "avg-body"
	} else if *sumBody {
		c.Sort = "sum-body"
	} else if *p1 {
		c.Sort = "p1"
	} else if *p50 {
		c.Sort = "p50"
	} else if *p99 {
		c.Sort = "p99"
	} else if *stddev {
		c.Sort = "stddev"
	} else {
		if c.Sort == "" {
			c.Sort = "max"
		}
	}

	c = SetConfig(c, option)

	if *load != "" {
		accessLog, err = LoadProfiles(*load)
		if err != nil {
			log.Fatal(err)
		}

		SortProfiles(accessLog, c)
		return
	}

	fileinfo, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}

	if fileinfo.Mode()&os.ModeNamedPipe == 0 {
		f, err = os.Open(c.File)
		defer f.Close()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		f = os.Stdin
	}

	accessLog = make(Profiles, 0, c.Limit)

	var includeRegexps []*regexp.Regexp
	if len(c.Includes) > 0 {
		includeRegexps = make([]*regexp.Regexp, 0, len(c.Includes))
		for _, pattern := range c.Includes {
			re, rerr := regexp.Compile(pattern)
			if rerr != nil {
				log.Fatal(err)
			}
			includeRegexps = append(includeRegexps, re)
		}
	}

	var excludeRegexps []*regexp.Regexp
	if len(c.Excludes) > 0 {
		excludeRegexps = make([]*regexp.Regexp, 0, len(c.Excludes))
		for _, pattern := range c.Excludes {
			re, rerr := regexp.Compile(pattern)
			if rerr != nil {
				log.Fatal(err)
			}
			excludeRegexps = append(excludeRegexps, re)
		}
	}

	var includeStatusRegexps []*regexp.Regexp
	if len(c.IncludeStatuses) > 0 {
		includeRegexps = make([]*regexp.Regexp, 0, len(c.Includes))
		for _, pattern := range c.IncludeStatuses {
			re, rerr := regexp.Compile(pattern)
			if rerr != nil {
				log.Fatal(err)
			}
			includeStatusRegexps = append(includeStatusRegexps, re)
		}
	}

	var excludeStatusRegexps []*regexp.Regexp
	if len(c.ExcludeStatuses) > 0 {
		excludeRegexps = make([]*regexp.Regexp, 0, len(c.Excludes))
		for _, pattern := range c.ExcludeStatuses {
			re, rerr := regexp.Compile(pattern)
			if rerr != nil {
				log.Fatal(err)
			}
			excludeStatusRegexps = append(excludeStatusRegexps, re)
		}
	}

	var aggregateRegexps []*regexp.Regexp
	if len(c.Aggregates) > 0 {
		aggregateRegexps = make([]*regexp.Regexp, 0, len(c.Aggregates))
		for _, pattern := range c.Aggregates {
			re, rerr := regexp.Compile(pattern)
			if rerr != nil {
				log.Fatal(err)
			}
			aggregateRegexps = append(aggregateRegexps, re)
		}
	}

	var p parsetime.ParseTime
	p, err = parsetime.NewParseTime(*location)

	var sTimeNano int64
	if c.StartTime != "" {
		sTime, err := p.Parse(c.StartTime)
		if err != nil {
			log.Fatal(err)
		}
		sTimeNano = sTime.UnixNano()
	}

	if c.StartTimeDuration != "" {
		sTime, err := TimeDurationSub(c.StartTimeDuration)
		if err != nil {
			log.Fatal(err)
		}
		sTimeNano = sTime.UnixNano()
	}

	var eTimeNano int64
	if c.EndTime != "" {
		eTime, err := p.Parse(c.EndTime)
		if err != nil {
			log.Fatal(err)
		}
		eTimeNano = eTime.UnixNano()
	}

	if c.EndTimeDuration != "" {
		eTime, err := TimeDurationSub(c.EndTimeDuration)
		if err != nil {
			log.Fatal(err)
		}
		eTimeNano = eTime.UnixNano()
	}

	r := ltsv.NewReader(f)
Loop:
	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		resTime, err := strconv.ParseFloat(line[c.ApptimeLabel], 64)
		if err != nil {
			var reqTime float64
			reqTime, err = strconv.ParseFloat(line[c.ReqtimeLabel], 64)
			if err != nil {
				continue
			}

			resTime = reqTime
		}

		bodySize, err := strconv.ParseFloat(line[c.SizeLabel], 64)
		if err != nil {
			continue
		}

		if sTimeNano != 0 || eTimeNano != 0 {
			t, err := p.Parse(line[c.TimeLabel])
			if err != nil {
				continue
			}
			timeNano := t.UnixNano()
			if !TimeCmp(sTimeNano, eTimeNano, timeNano) {
				continue
			}
		}

		u, err := url.Parse(line[c.UriLabel])
		if err != nil {
			continue
		}

		if c.QueryString {
			v := url.Values{}
			values := u.Query()
			for q, _ := range values {
				v.Set(q, "xxx")
			}
			uri = fmt.Sprintf("%s?%s", u.Path, v.Encode())
			index = fmt.Sprintf("%s_%s?%s", line[c.MethodLabel], u.Path, v.Encode())
		} else {
			uri = u.Path
			index = fmt.Sprintf("%s_%s", line[c.MethodLabel], u.Path)
		}

		if len(c.Includes) > 0 {
			isnotMatched := true
			for _, re := range includeRegexps {
				if ok := re.Match([]byte(uri)); ok && err == nil {
					isnotMatched = false
				} else if err != nil {
					log.Fatal(err)
				}
			}

			if isnotMatched {
				continue Loop
			}
		}

		if len(c.Excludes) > 0 {
			for _, re := range excludeRegexps {
				if ok := re.Match([]byte(uri)); ok && err == nil {
					continue Loop
				} else if err != nil {
					log.Fatal(err)
				}
			}
		}

		if len(c.IncludeStatuses) > 0 {
			isnotMatched := true
			for _, re := range includeStatusRegexps {
				if ok := re.Match([]byte(line[c.StatusLabel])); ok && err == nil {
					isnotMatched = false
				} else if err != nil {
					log.Fatal(err)
				}
			}

			if isnotMatched {
				continue Loop
			}
		}

		if len(c.ExcludeStatuses) > 0 {
			for _, re := range excludeStatusRegexps {
				if ok := re.Match([]byte(line[c.StatusLabel])); ok && err == nil {
					continue Loop
				} else if err != nil {
					log.Fatal(err)
				}
			}
		}

		isMatched := false
		if len(c.Aggregates) > 0 {
			for _, re := range aggregateRegexps {
				if ok := re.Match([]byte(uri)); ok && err == nil {
					isMatched = true
					pattern := re.String()
					index = fmt.Sprintf("%s_%s", line[c.MethodLabel], pattern)
					uri = pattern
					SetCursor(index, uri)
				} else if err != nil {
					log.Fatal(err)
				}
			}
		}

		if !isMatched {
			SetCursor(index, uri)
		}

		if len(uriHints) > c.Limit {
			log.Fatal(fmt.Sprintf("Too many uri (%d or less)", c.Limit))
		}

		if accessLog[cursor].Max < resTime {
			accessLog[cursor].Max = resTime
		}

		if accessLog[cursor].Min >= resTime || accessLog[cursor].Min == 0 {
			accessLog[cursor].Min = resTime
		}

		accessLog[cursor].Cnt++
		accessLog[cursor].Sum += resTime
		accessLog[cursor].Method = line[c.MethodLabel]
		accessLog[cursor].Percentails = append(accessLog[cursor].Percentails, Percentail{RequestTime: resTime})

		if accessLog[cursor].MaxBody < bodySize {
			accessLog[cursor].MaxBody = bodySize
		}

		if accessLog[cursor].MinBody >= bodySize || accessLog[cursor].MinBody == 0 {
			accessLog[cursor].MinBody = bodySize
		}

		accessLog[cursor].SumBody += bodySize
	}

	for i, _ := range accessLog {
		sort.Sort(ByRequestTime{accessLog[i].Percentails})
		accessLog[i].Avg = accessLog[i].Sum / float64(accessLog[i].Cnt)
		accessLog[i].AvgBody = accessLog[i].SumBody / float64(accessLog[i].Cnt)

		p1Len := LenPercentail(len(accessLog[i].Percentails), 1)
		accessLog[i].P1 = accessLog[i].Percentails[p1Len].RequestTime
		p50Len := LenPercentail(len(accessLog[i].Percentails), 50)
		accessLog[i].P50 = accessLog[i].Percentails[p50Len].RequestTime
		p99Len := LenPercentail(len(accessLog[i].Percentails), 99)
		accessLog[i].P99 = accessLog[i].Percentails[p99Len].RequestTime

		accessLog[i].Stddev = RequestTimeStddev(accessLog[i].Percentails, accessLog[i].Sum, accessLog[i].Avg)
	}

	SortProfiles(accessLog, c)

	if *dump != "" {
		err = DumpProfiles(*dump, accessLog)
		if err != nil {
			log.Fatal(err)
		}
	}
}
