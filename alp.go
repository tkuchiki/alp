package alp

import (
	"fmt"
	"io"
	"os"

	"github.com/tkuchiki/alp/flag"
	"github.com/tkuchiki/gohttpstats"
	"github.com/tkuchiki/gohttpstats/options"
	"github.com/tkuchiki/gohttpstats/parsers"
	"gopkg.in/alecthomas/kingpin.v2"
)

const version = "0.4.0"

type Profiler struct {
	outWriter    io.Writer
	errWriter    io.Writer
	inReader     *os.File
	optionParser *kingpin.Application
	flags        *flag.Flags
}

func NewProfiler(outw, errw io.Writer) *Profiler {
	p := &Profiler{
		outWriter:    outw,
		errWriter:    errw,
		inReader:     os.Stdin,
		optionParser: kingpin.New("alp", "Access Log Profiler for LTSV (read from file or stdin)."),
	}
	p.flags = flag.NewFlags()

	p.flags.InitFlags(p.optionParser)

	return p
}

func (p *Profiler) SetFlags(flags *flag.Flags) {
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
	p.optionParser.Version(version)
	kingpin.MustParse(p.optionParser.Parse(os.Args[1:]))

	sort := flag.SortOptions[p.flags.Sort]

	var err error
	var options *stats_options.Options
	if p.flags.Config != "" {
		cf, err := os.Open(p.flags.Config)
		if err != nil {
			return err
		}
		defer cf.Close()

		options, err = stats_options.LoadOptionsFromReader(cf)
		if err != nil {
			return err
		}
	} else {
		options = stats_options.NewOptions()
	}

	options = stats_options.SetOptions(options,
		stats_options.File(p.flags.File),
		stats_options.Sort(sort),
		stats_options.Reverse(p.flags.Reverse),
		stats_options.QueryString(p.flags.QueryString),
		stats_options.Tsv(p.flags.Tsv),
		stats_options.ApptimeLabel(p.flags.ApptimeLabel),
		stats_options.ReqtimeLabel(p.flags.ReqtimeLabel),
		stats_options.StatusLabel(p.flags.StatusLabel),
		stats_options.SizeLabel(p.flags.SizeLabel),
		stats_options.MethodLabel(p.flags.MethodLabel),
		stats_options.UriLabel(p.flags.UriLabel),
		stats_options.TimeLabel(p.flags.TimeLabel),
		stats_options.Limit(p.flags.Limit),
		stats_options.NoHeaders(p.flags.NoHeaders),
		stats_options.StartTime(p.flags.StartTime),
		stats_options.EndTime(p.flags.EndTime),
		stats_options.StartTimeDuration(p.flags.StartTimeDuration),
		stats_options.EndTimeDuration(p.flags.EndTimeDuration),
		stats_options.CSVIncludes(p.flags.Includes),
		stats_options.CSVExcludes(p.flags.Excludes),
		stats_options.CSVIncludeStatuses(p.flags.IncludeStatuses),
		stats_options.CSVExcludeStatuses(p.flags.ExcludeStatuses),
		stats_options.CSVAggregates(p.flags.Aggregates),
	)

	po := httpstats.NewPrintOptions()
	po.SetWriter(p.outWriter)
	if options.Tsv {
		po.SetFormat("tsv")
	}
	stats := httpstats.NewHTTPStats(true, false, false, po)

	err = stats.InitFilter(options)
	if err != nil {
		return err
	}

	stats.SetOptions(options)

	if p.flags.Load != "" {
		lf, err := os.Open(p.flags.Load)
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

	label := parsers.NewLTSVLabel(options.UriLabel, options.ApptimeLabel, options.ReqtimeLabel,
		options.SizeLabel, options.StatusLabel, options.MethodLabel, options.TimeLabel,
	)

	parser := parsers.NewLTSVParser(f, label, options.QueryString)
	if err != nil {
		return err
	}

Loop:
	for {
		s, err := parser.Parse()
		if err != nil {
			if err == io.EOF {
				break
			} else if err == httpstats.SkipReadLineErr {
				continue Loop
			}

			return err
		}

		if !stats.DoFilter(s.Uri, s.Method, s.Time) {
			continue Loop
		}

		stats.Set(s.Uri, s.Method, s.Status, s.ResponseTime, s.BodySize, 0)

		if stats.CountUris() > options.Limit {
			return fmt.Errorf("Too many URI's (%d or less)", options.Limit)
		}
	}

	if p.flags.Dump != "" {
		df, err := os.OpenFile(p.flags.Dump, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
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
