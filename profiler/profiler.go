package profiler

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/tkuchiki/alp/helpers"

	"github.com/tkuchiki/alp/errors"
	"github.com/tkuchiki/alp/options"
	"github.com/tkuchiki/alp/parsers"
	"github.com/tkuchiki/alp/stats"
)

type Profiler struct {
	options   *options.Options
	outWriter io.Writer
	errWriter io.Writer
	inReader  *os.File
}

func NewProfiler(outw, errw io.Writer, opts *options.Options) *Profiler {
	return &Profiler{
		options:   opts,
		outWriter: outw,
		errWriter: errw,
		inReader:  os.Stdin,
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

func (p *Profiler) Run(sortOptions *stats.SortOptions, parser parsers.Parser) error {
	sts := stats.NewHTTPStats(true, false, false)

	err := sts.InitFilter(p.options)
	if err != nil {
		return err
	}

	sts.SetOptions(p.options)
	sts.SetSortOptions(sortOptions)

	printOptions := stats.NewPrintOptions(p.options.NoHeaders, p.options.ShowFooters, p.options.DecodeUri, p.options.PaginationLimit)
	printer := stats.NewPrinter(p.outWriter, p.options.Output, p.options.Format, p.options.Percentiles, printOptions)
	if err = printer.Validate(); err != nil {
		return err
	}

	if p.options.Load != "" {
		lf, err := os.Open(p.options.Load)
		if err != nil {
			return err
		}
		err = sts.LoadStats(lf)
		if err != nil {
			return err
		}
		defer lf.Close()

		sts.SortWithOptions()
		printer.Print(sts, nil)
		return nil
	}

	if len(p.options.MatchingGroups) > 0 {
		err = sts.SetURIMatchingGroups(p.options.MatchingGroups)
		if err != nil {
			return err
		}
	}

	var posfile *os.File
	if p.options.PosFile != "" {
		posfile, err = p.OpenPosFile(p.options.PosFile)
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

		if sts.CountUris() > p.options.Limit {
			return fmt.Errorf("Too many URI's (%d or less)", p.options.Limit)
		}
	}

	if p.options.Dump != "" {
		df, err := os.OpenFile(p.options.Dump, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		err = sts.DumpStats(df)
		if err != nil {
			return err
		}
		defer df.Close()
	}

	if !p.options.NoSavePos && p.options.PosFile != "" {
		posfile.Seek(0, 0)
		_, err = posfile.Write([]byte(fmt.Sprint(parser.ReadBytes())))
		if err != nil {
			return err
		}
	}

	sts.SortWithOptions()
	printer.Print(sts, nil)

	return nil
}
