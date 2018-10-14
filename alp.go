package alp

import (
	"fmt"
	"io"
	"os"

	"github.com/tkuchiki/gohttpstats"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Profiler struct {
	outWriter    io.Writer
	errWriter    io.Writer
	inReader     *os.File
	optionParser *kingpin.Application
	flags        *httpstats.Flags
	version      string
}

func NewProfiler(outw, errw io.Writer) *Profiler {
	p := &Profiler{
		outWriter:    outw,
		errWriter:    errw,
		inReader:     os.Stdin,
		optionParser: kingpin.New("alp", "Access Log Profiler for LTSV (read from file or stdin)."),
		version:      "0.4.0",
	}
	p.flags = httpstats.NewFlags(p.optionParser)

	return p
}

func (p *Profiler) SetFlags(flags *httpstats.Flags) {
	p.flags = flags
}

func (p *Profiler) SetInReader(f *os.File) {
	p.inReader = f
}

func (p *Profiler) Open(filename string) (*os.File, error) {
	var f *os.File
	var err error

	if filename != "" {
		f, err = os.Open(filename)
	} else {
		f = p.inReader
	}

	return f, err
}

func (p *Profiler) Run() error {
	p.optionParser.Version(p.version)
	p.optionParser.Parse(os.Args[1:])

	var sort string
	if *p.flags.Max {
		sort = httpstats.SortMaxResponseTime
	} else if *p.flags.Min {
		sort = httpstats.SortMinResponseTime
	} else if *p.flags.Avg {
		sort = httpstats.SortAvgResponseTime
	} else if *p.flags.Sum {
		sort = httpstats.SortSumResponseTime
	} else if *p.flags.Cnt {
		sort = httpstats.SortCount
	} else if *p.flags.P1 {
		sort = httpstats.SortP1ResponseTime
	} else if *p.flags.P50 {
		sort = httpstats.SortP50ResponseTime
	} else if *p.flags.P99 {
		sort = httpstats.SortP99ResponseTime
	} else if *p.flags.Stddev {
		sort = httpstats.SortStddevResponseTime
	} else if *p.flags.SortUri {
		sort = httpstats.SortUri
	} else if *p.flags.Method {
		sort = httpstats.SortMethod
	} else if *p.flags.MaxBody {
		sort = httpstats.SortMaxResponseBodySize
	} else if *p.flags.MinBody {
		sort = httpstats.SortMinResponseBodySize
	} else if *p.flags.AvgBody {
		sort = httpstats.SortAvgResponseBodySize
	} else if *p.flags.SumBody {
		sort = httpstats.SortSumResponseBodySize
	} else {
		sort = httpstats.SortMaxResponseTime
	}

	cliOptions := httpstats.NewCliOptions(&httpstats.Options{
		File:              *p.flags.File,
		Sort:              sort,
		Reverse:           *p.flags.Reverse,
		QueryString:       *p.flags.QueryString,
		Tsv:               *p.flags.Tsv,
		ApptimeLabel:      *p.flags.ApptimeLabel,
		ReqtimeLabel:      *p.flags.ReqtimeLabel,
		StatusLabel:       *p.flags.StatusLabel,
		SizeLabel:         *p.flags.SizeLabel,
		MethodLabel:       *p.flags.MethodLabel,
		UriLabel:          *p.flags.UriLabel,
		TimeLabel:         *p.flags.TimeLabel,
		Limit:             *p.flags.Limit,
		NoHeaders:         *p.flags.NoHeaders,
		StartTime:         *p.flags.StartTime,
		EndTime:           *p.flags.EndTime,
		StartTimeDuration: *p.flags.StartTimeDuration,
		EndTimeDuration:   *p.flags.EndTimeDuration,
	},
		*p.flags.Includes, *p.flags.Excludes, *p.flags.IncludeStatuses, *p.flags.ExcludeStatuses, *p.flags.Aggregates)

	var err error
	var fileOptions *httpstats.Options
	if *p.flags.Config != "" {
		cf, err := os.Open(*p.flags.Config)
		if err != nil {
			return err
		}
		fileOptions, err = httpstats.LoadOptionsFromReader(cf)
		if err != nil {
			return err
		}
		defer cf.Close()
	} else {
		fileOptions = &httpstats.Options{}
	}

	defaultOptions := httpstats.NewDefaultOptions()
	options := httpstats.MergeOptions(cliOptions, fileOptions, defaultOptions)

	po := httpstats.NewPrintOptions()
	po.SetWriter(p.outWriter)
	stats := httpstats.NewHTTPStats(true, false, false, po)

	err = stats.InitFilter(p.flags, options)
	if err != nil {
		return err
	}

	stats.SetOptions(options)

	if *p.flags.Load != "" {
		lf, err := os.Open(*p.flags.Load)
		if err != nil {
			return err
		}
		err = stats.LoadStats(lf)
		if err != nil {
			return err
		}
		defer lf.Close()

		stats.SortWithOptions()
		stats.Print()
		return nil
	}

	f, err := p.Open(options.File)
	if err != nil {
		return err
	}
	defer f.Close()

	if len(options.Aggregates) > 0 {
		err = stats.SetURICapturingGroups(options.Aggregates)
		if err != nil {
			return err
		}
	}

	err = stats.InitParser("ltsv", f)
	if err != nil {
		return err
	}

Loop:
	for {
		uri, method, timestr, resTime, bodySize, status, err := stats.Parse()
		if err == io.EOF {
			break
		} else if err == httpstats.SkipReadLineErr {
			continue Loop
		} else if err != nil {
			return err
		}

		if !stats.DoFilter(uri, method, timestr) {
			continue Loop
		}

		stats.Set(uri, method, status, resTime, bodySize, 0)

		if stats.CountUris() > options.Limit {
			return fmt.Errorf("Too many URI's (%d or less)", options.Limit)
		}
	}

	if *p.flags.Dump != "" {
		df, err := os.OpenFile(*p.flags.Dump, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		err = stats.DumpStats(df)
		if err != nil {
			return err
		}
		defer df.Close()
	}

	stats.SortWithOptions()
	stats.Print()

	return nil
}
