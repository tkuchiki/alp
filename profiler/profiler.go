package profiler

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/tkuchiki/alp/errors"
	"github.com/tkuchiki/alp/helpers"
	"github.com/tkuchiki/alp/options"
	"github.com/tkuchiki/alp/parsers"
	"github.com/tkuchiki/alp/stats"
)

type Profiler struct {
	options     *options.Options
	outWriter   io.Writer
	errWriter   io.Writer
	inReader    *os.File
	printer     *stats.Printer
	loadEnabled bool
}

func NewProfiler(outw, errw io.Writer, opts *options.Options) *Profiler {
	printOptions := stats.NewPrintOptions(opts.NoHeaders, opts.ShowFooters, opts.DecodeUri, opts.PaginationLimit)
	printer := stats.NewPrinter(outw, opts.Output, opts.Format, opts.Percentiles, printOptions)

	return &Profiler{
		options:     opts,
		outWriter:   outw,
		errWriter:   errw,
		inReader:    os.Stdin,
		printer:     printer,
		loadEnabled: true,
	}
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

func (p *Profiler) SetPrinter(printer *stats.Printer) {
	p.printer = printer
}

func (p *Profiler) ValidatePrinter() error {
	return p.printer.Validate()
}

func (p *Profiler) Options() *options.Options {
	return p.options
}

func (p *Profiler) DisableLoad() {
	p.loadEnabled = false
}

func (p *Profiler) Load(sortOptions *stats.SortOptions) (*stats.HTTPStats, error) {
	sts := stats.NewHTTPStats(true, false, false)
	err := sts.InitFilter(p.options)
	if err != nil {
		return nil, err
	}

	sts.SetOptions(p.options)
	sts.SetSortOptions(sortOptions)

	lf, err := os.Open(p.options.Load)
	if err != nil {
		return nil, err
	}
	err = sts.LoadStats(lf)
	if err != nil {
		return nil, err
	}
	defer lf.Close()

	sts.SortWithOptions()

	return sts, nil
}

func (p *Profiler) profile(sortOptions *stats.SortOptions, parser parsers.Parser) (*stats.HTTPStats, error) {
	sts := stats.NewHTTPStats(true, false, false)

	if p.options.Load != "" && p.loadEnabled {
		return p.Load(sortOptions)
	}

	err := sts.InitFilter(p.options)
	if err != nil {
		return nil, err
	}

	sts.SetOptions(p.options)
	sts.SetSortOptions(sortOptions)

	if len(p.options.MatchingGroups) > 0 {
		err = sts.SetURIMatchingGroups(p.options.MatchingGroups)
		if err != nil {
			return nil, err
		}
	}

	var posfile *os.File
	if p.options.PosFile != "" {
		posfile, err = p.OpenPosFile(p.options.PosFile)
		if err != nil {
			return nil, err
		}
		defer posfile.Close()

		pos, err := p.ReadPosFile(posfile)
		if err != nil && err != io.EOF {
			return nil, err
		}

		err = parser.Seek(pos)
		if err != nil {
			return nil, err
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

			return nil, err
		}

		var b bool
		b, err = sts.DoFilter(s)
		if err != nil {
			return nil, err
		}

		if !b {
			continue Loop
		}

		sts.Set(s.Uri, s.Method, s.Status, s.ResponseTime, s.BodyBytes, 0)

		if sts.CountUris() > p.options.Limit {
			return nil, fmt.Errorf("Too many URI's (%d or less)", p.options.Limit)
		}
	}

	if !p.options.NoSavePos && p.options.PosFile != "" {
		posfile.Seek(0, 0)
		_, err = posfile.Write([]byte(fmt.Sprint(parser.ReadBytes())))
		if err != nil {
			return nil, err
		}
	}

	return sts, nil
}

func (p *Profiler) Profile(sortOptions *stats.SortOptions, parser parsers.Parser) (*stats.HTTPStats, error) {
	sts, err := p.profile(sortOptions, parser)
	if err != nil {
		return nil, err
	}

	sts.SortWithOptions()

	return sts, nil
}

func (p *Profiler) Run(sortOptions *stats.SortOptions, parser parsers.Parser, from *stats.HTTPStats) error {
	sts, err := p.profile(sortOptions, parser)
	if err != nil {
		return err
	}

	if p.options.Load != "" {
		if from == nil {
			p.printer.Print(sts, nil)
		} else {
			// diff
			p.printer.Print(from, sts)
		}

		return nil
	}

	if p.options.Dump != "" {
		df, err := os.OpenFile(p.options.Dump, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		err = sts.DumpStats(df)
		if err != nil {
			return err
		}
		defer df.Close()
	}

	sts.SortWithOptions()

	if from == nil {
		p.printer.Print(sts, nil)
	} else {
		// diff
		p.printer.Print(from, sts)
	}

	return nil
}
