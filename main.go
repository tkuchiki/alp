package main

import (
	"fmt"
	"github.com/najeira/ltsv"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"log"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

type Profile struct {
	Uri     string
	Cnt     int
	Max     float64
	Min     float64
	Sum     float64
	Avg     float64
	Method  string
	MaxBody float64
	MinBody float64
	SumBody float64
	AvgBody float64
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

func AbsPath(fname string) (f string, err error) {
	var fpath string
	matched, _ := regexp.Match("^~/", []byte(fname))
	if matched {
		usr, _ := user.Current()
		fpath = strings.Replace(fname, "~", usr.HomeDir, 1)
	} else {
		fpath, err = filepath.Abs(fname)
	}

	return fpath, err
}

func LoadFile(filename string) (f *os.File, err error) {
	fpath, err := AbsPath(filename)
	if err != nil {
		return f, err
	}
	f, err = os.Open(fpath)

	return f, err
}

func Round(f float64) string {
	return fmt.Sprintf("%.3f", f)
}

func Output(ps Profiles) {
	if *tsv {
		if !*noHeaders {
			fmt.Printf("Count\tMin\tMax\tSum\tAvg\tMax(Body)\tMin(Body)\tSum(Body)\tAvg(Body)\tMethod\tUri%v", eol)
		}

		for _, p := range ps {
			fmt.Printf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v%v",
				p.Cnt, Round(p.Min), Round(p.Max), Round(p.Sum), Round(p.Avg),
				Round(p.MinBody), Round(p.MaxBody), Round(p.SumBody), Round(p.AvgBody),
				p.Method, p.Uri, eol)
		}
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Count", "Min", "Max", "Sum", "Avg",
			"Max(Body)", "Min(Body)", "Sum(Body)", "Avg(Body)",
			"Method", "Uri"})
		for _, p := range ps {
			data := []string{
				fmt.Sprint(p.Cnt), Round(p.Min), Round(p.Max), Round(p.Sum), Round(p.Avg),
				Round(p.MinBody), Round(p.MaxBody), Round(p.SumBody), Round(p.AvgBody),
				p.Method, p.Uri}
			table.Append(data)
		}
		table.Render()
	}
}

func SortByMax(ps Profiles, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(ByMax{ps}))
	} else {
		sort.Sort(ByMax{ps})
	}
	Output(ps)
}

func SortByMin(ps Profiles, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(ByMin{ps}))
	} else {
		sort.Sort(ByMin{ps})
	}
	Output(ps)
}

func SortByAvg(ps Profiles, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(ByAvg{ps}))
	} else {
		sort.Sort(ByAvg{ps})
	}
	Output(ps)
}

func SortBySum(ps Profiles, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(BySum{ps}))
	} else {
		sort.Sort(BySum{ps})
	}
	Output(ps)
}

func SortByCnt(ps Profiles, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(ByCnt{ps}))
	} else {
		sort.Sort(ByCnt{ps})
	}
	Output(ps)
}

func SortByUri(ps Profiles, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(ByUri{ps}))
	} else {
		sort.Sort(ByUri{ps})
	}
	Output(ps)
}

func SortByMethod(ps Profiles, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(ByMethod{ps}))
	} else {
		sort.Sort(ByMethod{ps})
	}
	Output(ps)
}

func SortByMaxBody(ps Profiles, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(ByMaxBody{ps}))
	} else {
		sort.Sort(ByMaxBody{ps})
	}
	Output(ps)
}

func SortByMinBody(ps Profiles, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(ByMinBody{ps}))
	} else {
		sort.Sort(ByMinBody{ps})
	}
	Output(ps)
}

func SortByAvgBody(ps Profiles, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(ByAvgBody{ps}))
	} else {
		sort.Sort(ByAvgBody{ps})
	}
	Output(ps)
}

func SortBySumBody(ps Profiles, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(BySumBody{ps}))
	} else {
		sort.Sort(BySumBody{ps})
	}
	Output(ps)
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

var (
	file         = kingpin.Flag("file", "access log file").Short('f').String()
	max          = kingpin.Flag("max", "sort by max response time").Bool()
	min          = kingpin.Flag("min", "sort by min response time").Bool()
	avg          = kingpin.Flag("avg", "sort by avg response time").Bool()
	sum          = kingpin.Flag("sum", "sort by sum response time").Bool()
	cnt          = kingpin.Flag("cnt", "sort by count").Bool()
	sortUri      = kingpin.Flag("uri", "sort by uri").Bool()
	method       = kingpin.Flag("method", "sort by method").Bool()
	maxBody      = kingpin.Flag("max-body", "sort by max body size").Bool()
	minBody      = kingpin.Flag("min-body", "sort by min body size").Bool()
	avgBody      = kingpin.Flag("avg-body", "sort by avg body size").Bool()
	sumBody      = kingpin.Flag("sum-body", "sort by sum body size").Bool()
	reverse      = kingpin.Flag("reverse", "reverse the result of comparisons").Short('r').Bool()
	queryString  = kingpin.Flag("query-string", "include query string").Short('q').Bool()
	tsv          = kingpin.Flag("tsv", "tsv format (default: table)").Bool()
	apptimeLabel = kingpin.Flag("apptime-label", "apptime label").Default("apptime").String()
	sizeLabel    = kingpin.Flag("size-label", "size label").Default("size").String()
	methodLabel  = kingpin.Flag("method-label", "method label").Default("method").String()
	uriLabel     = kingpin.Flag("uri-label", "uri label").Default("uri").String()
	limit        = kingpin.Flag("limit", "set an upper limit of the target uri").Default("5000").Int()
	includes     = kingpin.Flag("includes", "don't exclude uri matching PATTERN (comma separated)").PlaceHolder("PATTERN,...").String()
	excludes     = kingpin.Flag("excludes", "exclude uri matching PATTERN (comma separated)").PlaceHolder("PATTERN,...").String()
	noHeaders    = kingpin.Flag("noheaders", "print no header line at all (only --tsv)").Bool()
	aggregates   = kingpin.Flag("aggregates", "aggregate uri matching PATTERN (comma separated)").PlaceHolder("PATTERN,...").String()

	eol = "\n"

	uri       string
	index     string
	accessLog Profiles
	uriHints      = make(map[string]int)
	length    int = 0
	cursor    int = 0
)

func main() {
	kingpin.Version("0.0.5")
	kingpin.Parse()

	var f *os.File
	var err error

	if *file != "" {
		f, err = LoadFile(*file)
		defer f.Close()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		f = os.Stdin
	}

	if runtime.GOOS == "windows" {
		eol = "\r\n"
	}

	var sortKey string

	if *max {
		sortKey = "max"
	} else if *min {
		sortKey = "min"
	} else if *avg {
		sortKey = "avg"
	} else if *sum {
		sortKey = "sum"
	} else if *cnt {
		sortKey = "cnt"
	} else if *sortUri {
		sortKey = "uri"
	} else if *method {
		sortKey = "method"
	} else if *maxBody {
		sortKey = "maxBody"
	} else if *minBody {
		sortKey = "minBody"
	} else if *avgBody {
		sortKey = "avgBody"
	} else if *sumBody {
		sortKey = "sumBody"
	} else {
		sortKey = "max"
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

		resTime, err := strconv.ParseFloat(line[*apptimeLabel], 64)
		if err != nil {
			continue
		}

		bodySize, err := strconv.ParseFloat(line[*sizeLabel], 64)
		if err != nil {
			continue
		}

		u, err := url.Parse(line[*uriLabel])
		if err != nil {
			log.Fatal(err)
		}

		if *queryString {
			v := url.Values{}
			values := u.Query()
			for q, _ := range values {
				v.Set(q, "xxx")
			}
			uri = fmt.Sprintf("%s?%s", u.Path, v.Encode())
			index = fmt.Sprintf("%s_%s?%s", line[*methodLabel], u.Path, v.Encode())
		} else {
			uri = u.Path
			index = fmt.Sprintf("%s_%s", line[*methodLabel], u.Path)
		}

		if *includes != "" {
			isnotMatched := true
			includePatterns := strings.Split(*includes, ",")
			for _, pattern := range includePatterns {
				if ok, err := regexp.Match(pattern, []byte(uri)); ok && err == nil {
					isnotMatched = false
				} else if err != nil {
					log.Fatal(err)
				}
			}

			if isnotMatched {
				continue Loop
			}
		}

		if *excludes != "" {
			excludePatterns := strings.Split(*excludes, ",")
			for _, pattern := range excludePatterns {
				if ok, err := regexp.Match(pattern, []byte(uri)); ok && err == nil {
					continue Loop
				} else if err != nil {
					log.Fatal(err)
				}
			}
		}

		isMatched := false
		if *aggregates != "" {
			aggregatePatterns := strings.Split(*aggregates, ",")
			for _, pattern := range aggregatePatterns {
				if ok, err := regexp.Match(pattern, []byte(uri)); ok && err == nil {
					isMatched = true
					index = fmt.Sprintf("%s_%s", line[*methodLabel], pattern)
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

		if len(uriHints) > *limit {
			log.Fatal(fmt.Sprintf("Too many uri (%d or less)", *limit))
		}

		if accessLog[cursor].Max < resTime {
			accessLog[cursor].Max = resTime
		}

		if accessLog[cursor].Min >= resTime || accessLog[cursor].Min == 0 {
			accessLog[cursor].Min = resTime
		}

		accessLog[cursor].Cnt++
		accessLog[cursor].Sum += resTime
		accessLog[cursor].Method = line[*methodLabel]

		if accessLog[cursor].MaxBody < bodySize {
			accessLog[cursor].MaxBody = bodySize
		}

		if accessLog[cursor].MinBody >= bodySize || accessLog[cursor].MinBody == 0 {
			accessLog[cursor].MinBody = bodySize
		}

		accessLog[cursor].SumBody += bodySize
	}

	for i, _ := range accessLog {
		accessLog[i].Avg = accessLog[i].Sum / float64(accessLog[i].Cnt)
		accessLog[i].AvgBody = accessLog[i].SumBody / float64(accessLog[i].Cnt)
	}

	switch sortKey {
	case "max":
		SortByMax(accessLog, *reverse)
	case "min":
		SortByMin(accessLog, *reverse)
	case "avg":
		SortByAvg(accessLog, *reverse)
	case "sum":
		SortBySum(accessLog, *reverse)
	case "cnt":
		SortByCnt(accessLog, *reverse)
	case "uri":
		SortByUri(accessLog, *reverse)
	case "method":
		SortByMethod(accessLog, *reverse)
	case "maxBody":
		SortByMaxBody(accessLog, *reverse)
	case "minBody":
		SortByMinBody(accessLog, *reverse)
	case "avgBody":
		SortByAvgBody(accessLog, *reverse)
	case "sumBody":
		SortBySumBody(accessLog, *reverse)
	}
}
