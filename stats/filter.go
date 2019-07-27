package stats

import (
	"regexp"
	"time"

	"github.com/tkuchiki/alp/errors"

	"github.com/tkuchiki/alp/parsers"

	"github.com/tkuchiki/alp/options"
	"github.com/tkuchiki/parsetime"
)

type Filter struct {
	options       *options.Options
	includeGroups []*regexp.Regexp
	excludeGroups []*regexp.Regexp
	expeval       *ExpEval
	parseTime     parsetime.ParseTime
}

func NewFilter(options *options.Options) *Filter {
	return &Filter{
		options: options,
	}
}

func (f *Filter) Init() error {
	var err error

	if len(f.options.Includes) > 0 {
		f.includeGroups, err = compileIncludeGroups(f.options.Includes)
		if err != nil {
			return err
		}
	}

	if len(f.options.Excludes) > 0 {
		f.excludeGroups, err = compileExcludeGroups(f.options.Excludes)
		if err != nil {
			return err
		}
	}

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
	if len(f.includeGroups) > 0 || len(f.excludeGroups) > 0 || f.expeval != nil {
		return true
	}

	return false
}

func (f *Filter) Do(stat *parsers.ParsedHTTPStat) error {
	if !f.isEnable() {
		return nil
	}

	if len(f.includeGroups) > 0 {
		isnotMatched := true
		for _, re := range f.includeGroups {
			if ok := re.Match([]byte(stat.Uri)); ok {
				isnotMatched = false
			}
		}

		if isnotMatched {
			return errors.SkipReadLineErr
		}
	}

	if len(f.excludeGroups) > 0 {
		for _, re := range f.excludeGroups {
			if ok := re.Match([]byte(stat.Uri)); ok {
				return errors.SkipReadLineErr
			}
		}
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

func compileIncludeGroups(includes []string) ([]*regexp.Regexp, error) {
	includeGroups := make([]*regexp.Regexp, 0, len(includes))
	for _, pattern := range includes {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return []*regexp.Regexp{}, err
		}
		includeGroups = append(includeGroups, re)
	}

	return includeGroups, nil
}

func compileExcludeGroups(excludes []string) ([]*regexp.Regexp, error) {
	excludeGroups := make([]*regexp.Regexp, 0, len(excludes))
	for _, pattern := range excludes {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return []*regexp.Regexp{}, err
		}
		excludeGroups = append(excludeGroups, re)
	}

	return excludeGroups, nil
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

func subTimeDuration(duration string) (int64, error) {
	var d time.Duration
	var err error
	var t time.Time
	d, err = time.ParseDuration(duration)
	if err != nil {
		return 0, err
	}

	t = time.Now().Add(-1 * d)

	return t.UnixNano(), err
}
