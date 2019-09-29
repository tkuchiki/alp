package alp

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/tkuchiki/alp/helpers"

	"github.com/tkuchiki/alp/errors"

	"github.com/tkuchiki/alp/stats"

	"github.com/tkuchiki/alp/flags"
	"github.com/tkuchiki/alp/options"
	"github.com/tkuchiki/alp/parsers"
	"gopkg.in/alecthomas/kingpin.v2"
)

const version = "1.0.2"

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
	app := kingpin.New("alp", "alp is the access log profiler for LTSV, JSON, and others.")
	p := &Profiler{
		outWriter:    outw,
		errWriter:    errw,
		inReader:     os.Stdin,
		optionParser: app,
	}

	p.globalFlags = flags.NewGlobalFlags()
	p.globalFlags.InitGlobalFlags(p.optionParser)

	p.subcmdLTSV = app.Command("ltsv", "Profile the logs for LTSV")
	p.ltsvFlags = flags.NewLTSVFlags()
	p.ltsvFlags.InitFlags(p.subcmdLTSV)

	p.subcmdJSON = app.Command("json", "Profile the logs for JSON")
	p.jsonFlags = flags.NewJSONFlags()
	p.jsonFlags.InitFlags(p.subcmdJSON)

	p.subcmdRegexp = app.Command("regexp", "Profile the logs that match a regular expression")
	p.regexpFlags = flags.NewRegexpFlags()
	p.regexpFlags.InitFlags(p.subcmdRegexp)

	return p
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

func (p *Profiler) OpenPosFile(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
}

func (p *Profiler) ReadPosFile(f *os.File) (int, error) {
	reader := bufio.NewReader(f)
	pos, _, err := reader.ReadLine()
	if err != nil {
		return 0, err
	}

	return helpers.StringToInt(string(pos))
}

func (p *Profiler) Run(args []string) error {
	var command string
	p.optionParser.Version(version)
	switch kingpin.MustParse(p.optionParser.Parse(args)) {
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
		sort = flags.SortOptions[opts.Sort]

	} else {
		opts = options.NewOptions()
	}

	opts = options.SetOptions(opts,
		options.File(p.globalFlags.File),
		options.Sort(sort),
		options.Reverse(p.globalFlags.Reverse),
		options.QueryString(p.globalFlags.QueryString),
		options.Format(p.globalFlags.Format),
		options.Limit(p.globalFlags.Limit),
		options.Location(p.globalFlags.Location),
		options.Output(p.globalFlags.Output),
		options.NoHeaders(p.globalFlags.NoHeaders),
		options.ShowFooters(p.globalFlags.ShowFooters),
		options.CSVGroups(p.globalFlags.MatchingGroups),
		options.Filters(p.globalFlags.Filters),
		options.PosFile(p.globalFlags.PosFile),
		options.NoSavePos(p.globalFlags.NoSavePos),
		// ltsv
		options.UriLabel(p.ltsvFlags.UriLabel),
		options.MethodLabel(p.ltsvFlags.MethodLabel),
		options.TimeLabel(p.ltsvFlags.TimeLabel),
		options.ApptimeLabel(p.ltsvFlags.ApptimeLabel),
		options.ReqtimeLabel(p.ltsvFlags.ReqtimeLabel),
		options.SizeLabel(p.ltsvFlags.SizeLabel),
		options.StatusLabel(p.ltsvFlags.StatusLabel),
		// json
		options.UriKey(p.jsonFlags.UriKey),
		options.MethodKey(p.jsonFlags.MethodKey),
		options.TimeKey(p.jsonFlags.TimeKey),
		options.ResponseTimeKey(p.jsonFlags.ResponseTimeKey),
		options.RequestTimeKey(p.jsonFlags.RequestTimeKey),
		options.BodyBytesKey(p.jsonFlags.BodyBytesKey),
		options.StatusKey(p.jsonFlags.StatusKey),
		// regexp
		options.UriSubexp(p.regexpFlags.UriSubexp),
		options.MethodSubexp(p.regexpFlags.MethodSubexp),
		options.TimeSubexp(p.regexpFlags.TimeSubexp),
		options.ResponseTimeSubexp(p.regexpFlags.ResponseTimeSubexp),
		options.RequestTimeSubexp(p.regexpFlags.RequestTimeSubexp),
		options.BodyBytesSubexp(p.regexpFlags.BodyBytesSubexp),
		options.StatusSubexp(p.regexpFlags.StatusSubexp),
	)

	sts := stats.NewHTTPStats(true, false, false)

	err = sts.InitFilter(opts)
	if err != nil {
		return err
	}

	sts.SetOptions(opts)

	printer := stats.NewPrinter(p.outWriter, opts.Output, opts.Format, opts.NoHeaders, opts.ShowFooters)
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
			opts.LTSV.ApptimeLabel, opts.LTSV.ReqtimeLabel, opts.LTSV.SizeLabel, opts.LTSV.StatusLabel,
		)
		parser = parsers.NewLTSVParser(f, label, opts.QueryString)
	case "json":
		keys := parsers.NewJSONKeys(opts.JSON.UriKey, opts.JSON.MethodKey, opts.JSON.TimeKey,
			opts.JSON.ResponseTimeKey, opts.JSON.RequestTimeKey, opts.JSON.BodyBytesKey, opts.JSON.StatusKey)
		parser = parsers.NewJSONParser(f, keys, opts.QueryString)
	case "regexp":
		names := parsers.NewSubexpNames(opts.Regexp.UriSubexp, opts.Regexp.MethodSubexp, opts.Regexp.TimeSubexp,
			opts.Regexp.ResponseTimeSubexp, opts.Regexp.RequestTimeSubexp, opts.Regexp.BodyBytesSubexp, opts.Regexp.StatusSubexp)
		parser, err = parsers.NewRegexpParser(f, options.DefaultPatternOption, names, opts.QueryString)

		if err != nil {
			return err
		}
	}

	var posfile *os.File
	if opts.PosFile != "" {
		posfile, err = p.OpenPosFile(opts.PosFile)
		if err != nil {
			return err
		}
		defer posfile.Close()

		pos, err := p.ReadPosFile(posfile)
		if err != nil && err != io.EOF {
			return err
		}

		err = parser.Seek(pos)
		if err != nil {
			return err
		}

		parser.SetReadBytes(pos)
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

	if !opts.NoSavePos && opts.PosFile != "" {
		posfile.Seek(0, 0)
		_, err = posfile.Write([]byte(fmt.Sprint(parser.ReadBytes())))
		if err != nil {
			return err
		}
	}

	sts.SortWithOptions()
	printer.Print(sts)

	return nil
}
