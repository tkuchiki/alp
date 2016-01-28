package main

import (
	"fmt"
	"github.com/najeira/ltsv"
	"github.com/olekukonko/tablewriter"
	"github.com/tkuchiki/parsetime"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"time"
)

type Percentail struct {
	RequestTime float64
}

type Percentails []Percentail

type Profile struct {
	Uri         string       `yaml:"uri"`
	Cnt         int          `yaml:"cnt"`
	Max         float64      `yaml:"max"`
	Min         float64      `yaml:"min"`
	Sum         float64      `yaml:"sum"`
	Avg         float64      `yaml:"anv"`
	Method      string       `yaml:"method"`
	MaxBody     float64      `yaml:"max_body"`
	MinBody     float64      `yaml:"min_body"`
	SumBody     float64      `yaml:"sum_body"`
	AvgBody     float64      `yaml:"avg_body"`
	Percentails []Percentail `yaml:"percentails"`
	P1          float64      `yaml:"p1"`
	P50         float64      `yaml:"p50"`
	P99         float64      `yaml:"p99"`
	Stddev      float64      `yaml:"stddev"`
}

type Profiles []Profile
type ByMax struct{ Profiles }
type ByMin struct{ Profiles }
type BySum struct{ Profiles }
type ByAvg struct{ Profiles }
type ByUri struct{ Profiles }
type ByCnt struct{ Profiles }
type ByMethod struct{ Profiles }
type ByMaxBody struct{ Profiles }
type ByMinBody struct{ Profiles }
type BySumBody struct{ Profiles }
type ByAvgBody struct{ Profiles }
type ByP1 struct{ Profiles }
type ByP50 struct{ Profiles }
type ByP99 struct{ Profiles }
type ByStddev struct{ Profiles }

func (s Profiles) Len() int      { return len(s) }
func (s Profiles) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ByMax) Less(i, j int) bool     { return s.Profiles[i].Max < s.Profiles[j].Max }
func (s ByMin) Less(i, j int) bool     { return s.Profiles[i].Min < s.Profiles[j].Min }
func (s BySum) Less(i, j int) bool     { return s.Profiles[i].Sum < s.Profiles[j].Sum }
func (s ByAvg) Less(i, j int) bool     { return s.Profiles[i].Avg < s.Profiles[j].Avg }
func (s ByUri) Less(i, j int) bool     { return s.Profiles[i].Uri < s.Profiles[j].Uri }
func (s ByCnt) Less(i, j int) bool     { return s.Profiles[i].Cnt < s.Profiles[j].Cnt }
func (s ByMethod) Less(i, j int) bool  { return s.Profiles[i].Method < s.Profiles[j].Method }
func (s ByMaxBody) Less(i, j int) bool { return s.Profiles[i].MaxBody < s.Profiles[j].MaxBody }
func (s ByMinBody) Less(i, j int) bool { return s.Profiles[i].MinBody < s.Profiles[j].MinBody }
func (s BySumBody) Less(i, j int) bool { return s.Profiles[i].SumBody < s.Profiles[j].SumBody }
func (s ByAvgBody) Less(i, j int) bool { return s.Profiles[i].AvgBody < s.Profiles[j].AvgBody }
func (s ByP1) Less(i, j int) bool      { return s.Profiles[i].P1 < s.Profiles[j].P1 }
func (s ByP50) Less(i, j int) bool     { return s.Profiles[i].P50 < s.Profiles[j].P50 }
func (s ByP99) Less(i, j int) bool     { return s.Profiles[i].P99 < s.Profiles[j].P99 }
func (s ByStddev) Less(i, j int) bool  { return s.Profiles[i].Stddev < s.Profiles[j].Stddev }

type ByRequestTime struct{ Percentails }

func (s Percentails) Len() int      { return len(s) }
func (s Percentails) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ByRequestTime) Less(i, j int) bool {
	return s.Percentails[i].RequestTime < s.Percentails[j].RequestTime
}

const (
	ApptimeLabel = "apptime"
	SizeLabel    = "size"
	MethodLabel  = "method"
	UriLabel     = "uri"
	TimeLabel    = "time"
	Limit        = 5000
)

func LoadYAML(filename string) (config Config, err error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(buf, &config)

	return config, err
}

func DumpProfiles(filename string, ps Profiles) (err error) {
	buf, err := yaml.Marshal(&ps)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, buf, os.ModePerm)

	return err
}

func LoadProfiles(filename string) (ps Profiles, err error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return ps, err
	}

	err = yaml.Unmarshal(buf, &ps)

	return ps, err
}

func Round(f float64) string {
	return fmt.Sprintf("%.3f", f)
}

func Output(ps Profiles, c Config) {
	if c.Tsv {
		if !c.NoHeaders {
			fmt.Printf("Count\tMin\tMax\tSum\tAvg\tP1\tP50\tP99\tStddev\tMin(Body)\tMax(Body)\tSum(Body)\tAvg(Body)\tMethod\tUri%v", eol)
		}

		for _, p := range ps {
			fmt.Printf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v%v",
				p.Cnt, Round(p.Min), Round(p.Max), Round(p.Sum), Round(p.Avg),
				Round(p.P1), Round(p.P50), Round(p.P99), Round(p.Stddev),
				Round(p.MinBody), Round(p.MaxBody), Round(p.SumBody), Round(p.AvgBody),
				p.Method, p.Uri, eol)
		}
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Count", "Min", "Max", "Sum", "Avg",
			"P1", "P50", "P99", "Stddev",
			"Min(Body)", "Max(Body)", "Sum(Body)", "Avg(Body)",
			"Method", "Uri"})
		for _, p := range ps {
			data := []string{
				fmt.Sprint(p.Cnt), Round(p.Min), Round(p.Max), Round(p.Sum), Round(p.Avg),
				Round(p.P1), Round(p.P50), Round(p.P99), Round(p.Stddev),
				Round(p.MinBody), Round(p.MaxBody), Round(p.SumBody), Round(p.AvgBody),
				p.Method, p.Uri}
			table.Append(data)
		}
		table.Render()
	}
}

func SortByMax(ps Profiles, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByMax{ps}))
	} else {
		sort.Sort(ByMax{ps})
	}
	Output(ps, c)
}

func SortByMin(ps Profiles, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByMin{ps}))
	} else {
		sort.Sort(ByMin{ps})
	}
	Output(ps, c)
}

func SortByAvg(ps Profiles, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByAvg{ps}))
	} else {
		sort.Sort(ByAvg{ps})
	}
	Output(ps, c)
}

func SortBySum(ps Profiles, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(BySum{ps}))
	} else {
		sort.Sort(BySum{ps})
	}
	Output(ps, c)
}

func SortByCnt(ps Profiles, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByCnt{ps}))
	} else {
		sort.Sort(ByCnt{ps})
	}
	Output(ps, c)
}

func SortByUri(ps Profiles, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByUri{ps}))
	} else {
		sort.Sort(ByUri{ps})
	}
	Output(ps, c)
}

func SortByMethod(ps Profiles, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByMethod{ps}))
	} else {
		sort.Sort(ByMethod{ps})
	}
	Output(ps, c)
}

func SortByMaxBody(ps Profiles, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByMaxBody{ps}))
	} else {
		sort.Sort(ByMaxBody{ps})
	}
	Output(ps, c)
}

func SortByMinBody(ps Profiles, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByMinBody{ps}))
	} else {
		sort.Sort(ByMinBody{ps})
	}
	Output(ps, c)
}

func SortByAvgBody(ps Profiles, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByAvgBody{ps}))
	} else {
		sort.Sort(ByAvgBody{ps})
	}
	Output(ps, c)
}

func SortBySumBody(ps Profiles, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(BySumBody{ps}))
	} else {
		sort.Sort(BySumBody{ps})
	}
	Output(ps, c)
}

func SortByP1(ps Profiles, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByP1{ps}))
	} else {
		sort.Sort(ByP1{ps})
	}
	Output(ps, c)
}

func SortByP50(ps Profiles, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByP50{ps}))
	} else {
		sort.Sort(ByP50{ps})
	}
	Output(ps, c)
}

func SortByP99(ps Profiles, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByP99{ps}))
	} else {
		sort.Sort(ByP99{ps})
	}
	Output(ps, c)
}

func SortByStddev(ps Profiles, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByStddev{ps}))
	} else {
		sort.Sort(ByStddev{ps})
	}
	Output(ps, c)
}

func SortProfiles(accessLog Profiles, c Config) {
	switch c.Sort {
	case "max":
		SortByMax(accessLog, c)
	case "min":
		SortByMin(accessLog, c)
	case "avg":
		SortByAvg(accessLog, c)
	case "sum":
		SortBySum(accessLog, c)
	case "cnt":
		SortByCnt(accessLog, c)
	case "uri":
		SortByUri(accessLog, c)
	case "method":
		SortByMethod(accessLog, c)
	case "max-body":
		SortByMaxBody(accessLog, c)
	case "min-body":
		SortByMinBody(accessLog, c)
	case "avg-body":
		SortByAvgBody(accessLog, c)
	case "sum-body":
		SortBySumBody(accessLog, c)
	case "p1":
		SortByP1(accessLog, c)
	case "p50":
		SortByP50(accessLog, c)
	case "p99":
		SortByP99(accessLog, c)
	case "stddev":
		SortByStddev(accessLog, c)
	}
}

func SetCursor(index string, uri string) {
	if _, ok := uriHints[index]; ok {
		cursor = uriHints[index]
	} else {
		uriHints[index] = length
		cursor = length
		length++
		accessLog = append(accessLog, Profile{Uri: uri})
	}
}

func TimeCmp(startTimeNano int64, endTimeNano int64, timeNano int64) bool {
	if startTimeNano > 0 && endTimeNano == 0 {
		return startTimeNano <= timeNano
	} else if endTimeNano > 0 && startTimeNano == 0 {
		return endTimeNano >= timeNano
	} else if startTimeNano > 0 && endTimeNano > 0 {
		return startTimeNano <= timeNano && endTimeNano >= timeNano
	}

	return false
}

func TimeDurationSub(duration string) (t time.Time, err error) {
	var d time.Duration
	d, err = time.ParseDuration(duration)
	if err != nil {
		return t, err
	}

	t = time.Now().Add(-1 * d)

	return t, err
}

func LenPercentail(l int, n int) (pLen int) {
	pLen = (l * n / 100) - 1
	if pLen < 0 {
		pLen = 0
	}

	return pLen
}

func RequestTimeStddev(requests Percentails, sum, avg float64) (stddev float64) {
	n := float64(len(requests))
	for _, r := range requests {
		stddev += (r.RequestTime - avg) * (r.RequestTime - avg)
	}

	return math.Sqrt(stddev / n)
}

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
	sizeLabel         = kingpin.Flag("size-label", "size label").Default(SizeLabel).String()
	methodLabel       = kingpin.Flag("method-label", "method label").Default(MethodLabel).String()
	uriLabel          = kingpin.Flag("uri-label", "uri label").Default(UriLabel).String()
	timeLabel         = kingpin.Flag("time-label", "time label").Default(TimeLabel).String()
	limit             = kingpin.Flag("limit", "set an upper limit of the target uri").Default(fmt.Sprint(Limit)).Int()
	includes          = kingpin.Flag("includes", "don't exclude uri matching PATTERN (comma separated)").PlaceHolder("PATTERN,...").String()
	excludes          = kingpin.Flag("excludes", "exclude uri matching PATTERN (comma separated)").PlaceHolder("PATTERN,...").String()
	noHeaders         = kingpin.Flag("noheaders", "print no header line at all (only --tsv)").Bool()
	aggregates        = kingpin.Flag("aggregates", "aggregate uri matching PATTERN (comma separated)").PlaceHolder("PATTERN,...").String()
	startTime         = kingpin.Flag("start-time", "since the start time").PlaceHolder("TIME").String()
	endTime           = kingpin.Flag("end-time", "end time earlier").PlaceHolder("TIME").String()
	startTimeDuration = kingpin.Flag("start-time-duration", "since the start time (now - time.Duration)").PlaceHolder("TIME_DURATION").String()
	endTimeDuration   = kingpin.Flag("end-time-duration", "end time earlier (now - time.Duration)").PlaceHolder("TIME_DURATION").String()

	eol = "\n"

	uri       string
	index     string
	accessLog Profiles
	uriHints      = make(map[string]int)
	length    int = 0
	cursor    int = 0
)

func main() {
	kingpin.CommandLine.Help = "Access Log Profiler for LTSV (read from file or stdin)."
	kingpin.Version("0.2.1")
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
		File:              *file,
		Reverse:           *reverse,
		QueryString:       *queryString,
		Tsv:               *tsv,
		ApptimeLabel:      *apptimeLabel,
		SizeLabel:         *sizeLabel,
		MethodLabel:       *methodLabel,
		UriLabel:          *uriLabel,
		TimeLabel:         *timeLabel,
		Limit:             *limit,
		IncludesStr:       *includes,
		ExcludesStr:       *excludes,
		NoHeaders:         *noHeaders,
		AggregatesStr:     *aggregates,
		StartTime:         *startTime,
		EndTime:           *endTime,
		StartTimeDuration: *startTimeDuration,
		EndTimeDuration:   *endTimeDuration,
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

	if runtime.GOOS == "windows" {
		eol = "\r\n"
	}

	var includeRegexps []*regexp.Regexp
	if len(c.Includes) > 0 {
		includeRegexps = make([]*regexp.Regexp, 0, len(c.Includes))
		for _, pattern := range c.Includes {
			includeRegexps = append(includeRegexps, regexp.MustCompile(pattern))
		}
	}

	var excludeRegexps []*regexp.Regexp
	if len(c.Excludes) > 0 {
		excludeRegexps = make([]*regexp.Regexp, 0, len(c.Excludes))
		for _, pattern := range c.Excludes {
			excludeRegexps = append(excludeRegexps, regexp.MustCompile(pattern))
		}
	}

	var aggregateRegexps []*regexp.Regexp
	if len(c.Aggregates) > 0 {
		aggregateRegexps = make([]*regexp.Regexp, 0, len(c.Aggregates))
		for _, pattern := range c.Aggregates {
			aggregateRegexps = append(aggregateRegexps, regexp.MustCompile(pattern))
		}
	}

	var sTimeNano int64
	if c.StartTime != "" {
		sTime, err := parsetime.Parse(c.StartTime)
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
		eTime, err := parsetime.Parse(c.EndTime)
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
			continue
		}

		bodySize, err := strconv.ParseFloat(line[c.SizeLabel], 64)
		if err != nil {
			continue
		}

		if sTimeNano != 0 || eTimeNano != 0 {
			t, err := parsetime.Log(line[c.TimeLabel])
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
			log.Fatal(err)
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
