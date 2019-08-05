package stats

import (
	"time"

	"github.com/tkuchiki/alp/errors"
	"github.com/tkuchiki/alp/options"
	"github.com/tkuchiki/alp/parsers"
	"github.com/tkuchiki/parsetime"
)

type Filter struct {
	options   *options.Options
	expeval   *ExpEval
	parseTime parsetime.ParseTime
}

func NewFilter(options *options.Options) *Filter {
	return &Filter{
		options: options,
	}
}

func (f *Filter) Init() error {
	var err error

	err = f.InitParseTime(f.options.Location)
	if err != nil {
		return err
	}

	if f.options.Filters != "" {
		var ee *ExpEval
		ee, err = NewExpEval(f.options.Filters, f.parseTime)
		if err != nil {
			return err
		}

		f.expeval = ee
	}

	return nil
}

func (f *Filter) isEnable() bool {
	if f.expeval != nil {
		return true
	}

	return false
}

func (f *Filter) Do(stat *parsers.ParsedHTTPStat) error {
	if !f.isEnable() {
		return nil
	}

	if f.expeval != nil {
		matched, err := f.expeval.Run(stat)
		if err != nil {
			return err
		}

		if !matched {
			return errors.SkipReadLineErr
		}
	}

	return nil
}

func (f *Filter) InitParseTime(loc string) error {
	p, err := parsetime.NewParseTime(loc)
	f.parseTime = p
	return err
}

func (f *Filter) ParseTime(val string) (time.Time, error) {
	return f.parseTime.Parse(val)
}

func (f *Filter) TimeStrToUnixNano(val string) (int64, error) {
	t, err := f.parseTime.Parse(val)
	if err != nil {
		return 0, err
	}

	return t.UnixNano(), nil
}
