package alp

import (
	"fmt"
	"io"
	"os"

	"github.com/tkuchiki/alp/errors"

	"github.com/tkuchiki/alp/stats"

	"github.com/tkuchiki/alp/flags"
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
	subcmdLTSV   *kingpin.CmdClause
	subcmdRegexp *kingpin.CmdClause
	subcmdJSON   *kingpin.CmdClause
	globalFlags  *flags.GlobalFlags
	ltsvFlags    *flags.LTSVFlags
	regexpFlags  *flags.RegexpFlags
	jsonFlags    *flags.JSONFlags
}

func NewProfiler(outw, errw io.Writer) *Profiler {
	app := kingpin.New("alp", "Access Log Profiler for LTSV (read from file or stdin).")
	p := &Profiler{
		outWriter:    outw,
		errWriter:    errw,
		inReader:     os.Stdin,
		optionParser: app,
	}

	p.subcmdLTSV = app.Command("ltsv", "ltsv")
	p.subcmdRegexp = app.Command("regexp", "regexp")
	p.subcmdJSON = app.Command("json", "json")

	p.globalFlags = flags.NewGlobalFlags()
	p.globalFlags.InitGlobalFlags(p.optionParser)

	p.ltsvFlags = flags.NewLTSVFlags()
	p.ltsvFlags.InitFlags(p.subcmdLTSV)

	p.regexpFlags = flags.NewRegexpFlags()
	p.regexpFlags.InitFlags(p.subcmdRegexp)

	p.jsonFlags = flags.NewJSONFlags()
	p.jsonFlags.InitFlags(p.subcmdJSON)

	return p
}

func (p *Profiler) SetFlags(flags *flags.GlobalFlags) {
	p.globalFlags = flags
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
	var command string
	p.optionParser.Version(version)
	switch kingpin.MustParse(p.optionParser.Parse(os.Args[1:])) {
	case p.subcmdLTSV.FullCommand():
		command = p.subcmdLTSV.FullCommand()
	case p.subcmdRegexp.FullCommand():
		command = p.subcmdRegexp.FullCommand()
	case p.subcmdJSON.FullCommand():
		command = p.subcmdJSON.FullCommand()
	}

	sort := flags.SortOptions[p.globalFlags.Sort]

	var err error
	var opts *options.Options
	if p.globalFlags.Config != "" {
		cf, err := os.Open(p.globalFlags.Config)
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
		options.File(p.globalFlags.File),
		options.Sort(sort),
		options.Reverse(p.globalFlags.Reverse),
		options.QueryString(p.globalFlags.QueryString),
		options.Tsv(p.globalFlags.Tsv),
		options.Limit(p.globalFlags.Limit),
		options.Location(p.globalFlags.Location),
		options.Output(p.globalFlags.Output),
		options.NoHeaders(p.globalFlags.NoHeaders),
		options.CSVGroups(p.globalFlags.MatchingGroups),
		options.Filters(p.globalFlags.Filters),
		//ltsv
		options.ApptimeLabel(p.ltsvFlags.ApptimeLabel),
		options.StatusLabel(p.ltsvFlags.StatusLabel),
		options.SizeLabel(p.ltsvFlags.SizeLabel),
		options.MethodLabel(p.ltsvFlags.MethodLabel),
		options.UriLabel(p.ltsvFlags.UriLabel),
		options.TimeLabel(p.ltsvFlags.TimeLabel),
	)

	sts := stats.NewHTTPStats(true, false, false)

	err = sts.InitFilter(opts)
	if err != nil {
		return err
	}

	sts.SetOptions(opts)

	printFormat := "table"
	if opts.Tsv {
		printFormat = "tsv"
	}
	printer := stats.NewPrinter(p.outWriter, opts.Output, printFormat)
	if err = printer.Validate(); err != nil {
		return err
	}

	if p.globalFlags.Load != "" {
		lf, err := os.Open(p.globalFlags.Load)
		if err != nil {
			return err
		}
		err = sts.LoadStats(lf)
		if err != nil {
			return err
		}
		defer lf.Close()

		sts.SortWithOptions()
		printer.Print(sts)
		return nil
	}

	f, err := p.Open(opts.File)
	if err != nil {
		return err
	}
	defer f.Close()

	if len(opts.MatchingGroups) > 0 {
		err = sts.SetURIMatchingGroups(opts.MatchingGroups)
		if err != nil {
			return err
		}
	}

	var parser parsers.Parser
	switch command {
	case "ltsv":
		label := parsers.NewLTSVLabel(opts.LTSV.UriLabel, opts.LTSV.MethodLabel, opts.LTSV.TimeLabel,
			opts.LTSV.ApptimeLabel, opts.LTSV.SizeLabel, opts.LTSV.StatusLabel,
		)
		parser = parsers.NewLTSVParser(f, label, opts.QueryString)
	case "regexp":
		expr := `^(\S+)\s` + // remote host
			`\S+\s+` +
			`(\S+\s+)+` + // user
			`\[(?P<time>[^]]+)\]\s` + // time
			`"(?P<method>\S*)\s?` + // method
			`(?P<uri>(?:[^"]*(?:\\")?)*)\s` + // URL
			`([^"]*)"\s` + // protocol
			`(?P<status>\S+)\s` + // status code
			`(?P<body_bytes>\S+)\s` + // bytes
			`"((?:[^"]*(?:\\")?)*)"\s` + // referer
			`"(.*)"` + // user agent
			`\s(?P<response_time>.*)$`
		names := parsers.NewSubexpNames("", "", "", "", "", "")
		parser, err = parsers.NewRegexpParser(f, expr, names, opts.QueryString)
	case "json":
		keys := parsers.NewJSONKeys("", "", "", "", "", "")
		parser = parsers.NewJSONParser(f, keys, opts.QueryString)
	}

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

		sts.Set(s.Uri, s.Method, s.Status, s.ResponseTime, s.BodyBytes, 0)

		if sts.CountUris() > opts.Limit {
			return fmt.Errorf("Too many URI's (%d or less)", opts.Limit)
		}
	}

	if p.globalFlags.Dump != "" {
		df, err := os.OpenFile(p.globalFlags.Dump, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		err = sts.DumpStats(df)
		if err != nil {
			return err
		}
		defer df.Close()
	}

	sts.SortWithOptions()
	printer.Print(sts)

	return nil
}
