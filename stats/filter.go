package stats

import (
	"regexp"
	"time"

	"github.com/tkuchiki/alp/options"
	"github.com/tkuchiki/parsetime"
)

type Filter struct {
	options       *options.Options
	includeGroups []*regexp.Regexp
	excludeGroups []*regexp.Regexp
	sTimeNano     int64
	eTimeNano     int64
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

	if f.options.StartTime != "" {
		f.sTimeNano, err = f.TimeStrToUnixNano(f.options.StartTime)
		if err != nil {
			return err
		}
	}

	if f.options.StartTimeDuration != "" {
		f.sTimeNano, err = subTimeDuration(f.options.StartTimeDuration)
		if err != nil {
			return err
		}
	}

	if f.options.EndTime != "" {
		f.eTimeNano, err = f.TimeStrToUnixNano(f.options.EndTime)
		if err != nil {
			return err
		}
	}

	if f.options.EndTimeDuration != "" {
		f.eTimeNano, err = subTimeDuration(f.options.EndTimeDuration)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *Filter) isEnable() bool {
	if len(f.includeGroups) > 0 || len(f.excludeGroups) > 0 ||
		f.options.StartTime != "" || f.options.StartTimeDuration != "" ||
		f.options.EndTime != "" || f.options.EndTimeDuration != "" {
		return true
	}

	return false
}

func (f *Filter) Do(uri, status, timestr string) error {
	if !f.isEnable() {
		return nil
	}

	if len(f.includeGroups) > 0 {
		isnotMatched := true
		for _, re := range f.includeGroups {
			if ok := re.Match([]byte(uri)); ok {
				isnotMatched = false
			}
		}

		if isnotMatched {
			return SkipReadLineErr
		}
	}

	if len(f.excludeGroups) > 0 {
		for _, re := range f.excludeGroups {
			if ok := re.Match([]byte(uri)); ok {
				return SkipReadLineErr
			}
		}
	}

	if f.sTimeNano != 0 || f.eTimeNano != 0 {
		t, err := f.ParseTime(timestr)
		if err != nil {
			return SkipReadLineErr
		}
		timeNano := t.UnixNano()
		if !IsIncludedInTime(f.sTimeNano, f.eTimeNano, timeNano) {
			return SkipReadLineErr
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
