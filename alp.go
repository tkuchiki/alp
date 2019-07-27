package alp

import (
	"fmt"
	"io"
	"os"

	"github.com/tkuchiki/alp/errors"

	"github.com/tkuchiki/alp/stats"

	"github.com/tkuchiki/alp/flag"
	"github.com/tkuchiki/alp/options"
	"github.com/tkuchiki/alp/parsers"
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
	var opts *options.Options
	if p.flags.Config != "" {
		cf, err := os.Open(p.flags.Config)
		if err != nil {
			return err
		}
		defer cf.Close()

		opts, err = options.LoadOptionsFromReader(cf)
		if err != nil {
			return err
		}
	} else {
		opts = options.NewOptions()
	}

	opts = options.SetOptions(opts,
		options.File(p.flags.File),
		options.Sort(sort),
		options.Reverse(p.flags.Reverse),
		options.QueryString(p.flags.QueryString),
		options.Tsv(p.flags.Tsv),
		options.ApptimeLabel(p.flags.ApptimeLabel),
		options.ReqtimeLabel(p.flags.ReqtimeLabel),
		options.StatusLabel(p.flags.StatusLabel),
		options.SizeLabel(p.flags.SizeLabel),
		options.MethodLabel(p.flags.MethodLabel),
		options.UriLabel(p.flags.UriLabel),
		options.TimeLabel(p.flags.TimeLabel),
		options.Limit(p.flags.Limit),
		options.Location(p.flags.Location),
		options.NoHeaders(p.flags.NoHeaders),
		options.CSVIncludes(p.flags.Includes),
		options.CSVExcludes(p.flags.Excludes),
		options.CSVGroups(p.flags.Groups),
		options.Filters(p.flags.Filters),
	)

	po := stats.NewPrintOptions()
	po.SetWriter(p.outWriter)
	if opts.Tsv {
		po.SetFormat("tsv")
	}
	sts := stats.NewHTTPStats(true, false, false, po)

	err = sts.InitFilter(opts)
	if err != nil {
		return err
	}

	sts.SetOptions(opts)

	if p.flags.Load != "" {
		lf, err := os.Open(p.flags.Load)
		if err != nil {
			return err
		}
		err = sts.LoadStats(lf)
		if err != nil {
			return err
		}
		defer lf.Close()

		sts.SortWithOptions()
		sts.Print()
		return nil
	}

	f, err := p.Open(opts.File)
	if err != nil {
		return err
	}
	defer f.Close()

	if len(opts.Groups) > 0 {
		err = sts.SetURICapturingGroups(opts.Groups)
		if err != nil {
			return err
		}
	}

	label := parsers.NewLTSVLabel(opts.UriLabel, opts.ApptimeLabel, opts.ReqtimeLabel,
		opts.SizeLabel, opts.StatusLabel, opts.MethodLabel, opts.TimeLabel,
	)

	parser := parsers.NewLTSVParser(f, label, opts.QueryString)
	if err != nil {
		return err
	}

Loop:
	for {
		s, err := parser.Parse()
		if err != nil {
			if err == io.EOF {
				break
			} else if err == errors.SkipReadLineErr {
				continue Loop
			}

			return err
		}

		var b bool
		b, err = sts.DoFilter(s)
		if err != nil {
			return err
		}

		if !b {
			continue Loop
		}

		sts.Set(s.Uri, s.Method, s.Status, s.ResponseTime, s.BodySize, 0)

		if sts.CountUris() > opts.Limit {
			return fmt.Errorf("Too many URI's (%d or less)", opts.Limit)
		}
	}

	if p.flags.Dump != "" {
		df, err := os.OpenFile(p.flags.Dump, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		err = sts.DumpStats(df)
		if err != nil {
			return err
		}
		defer df.Close()
	}

	sts.SortWithOptions()
	sts.Print()

	return nil
}
