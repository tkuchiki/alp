package options

import (
	"io"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	DefaultSortOption         = "max"
	DefaultApptimeLabelOption = "apptime"
	DefaultReqtimeLabelOption = "reqtime"
	DefaultStatusLabelOption  = "status"
	DefaultSizeLabelOption    = "size"
	DefaultMethodLabelOption  = "method"
	DefaultUriLabelOption     = "uri"
	DefaultTimeLabelOption    = "time"
	DefaultLimitOption        = 5000
)

func splitCSV(val string) []string {
	strs := strings.Split(val, ",")
	if len(strs) == 1 && strs[0] == "" {
		return []string{}
	}

	trimedStrs := make([]string, 0, len(strs))

	for _, s := range strs {
		trimedStrs = append(trimedStrs, strings.Trim(s, " "))
	}

	return trimedStrs
}

type Options struct {
	File              string   `yaml:"file"`
	Sort              string   `yaml:"sort"`
	Reverse           bool     `yaml:"reverse"`
	QueryString       bool     `yaml:"query_string"`
	Tsv               bool     `yaml:"tsv"`
	NoHeaders         bool     `yaml:no_headers`
	ApptimeLabel      string   `yaml:"apptime_label"`
	ReqtimeLabel      string   `yaml:"reqtime_label"`
	StatusLabel       string   `yaml:"status_label"`
	SizeLabel         string   `yaml:"size_label"`
	MethodLabel       string   `yaml:"method_label"`
	UriLabel          string   `yaml:"uri_label"`
	TimeLabel         string   `yaml:"time_label"`
	Limit             int      `yaml:"limit"`
	Includes          []string `yaml:"includes"`
	Excludes          []string `yaml:"excludes"`
	Groups            []string `yaml:"groups"`
	StartTime         string   `yaml:"start_time"`
	EndTime           string   `yaml:"end_time"`
	StartTimeDuration string   `yaml:"start_time_duration"`
	EndTimeDuration   string   `yaml:"end_time_duration"`
	Location          string   `yaml:location`
}

type Option func(*Options)

func File(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.File = s
		}
	}
}

func Sort(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.Sort = s
		}
	}
}

func Reverse(b bool) Option {
	return func(opts *Options) {
		if b {
			opts.Reverse = b
		}
	}
}

func QueryString(b bool) Option {
	return func(opts *Options) {
		if b {
			opts.QueryString = b
		}
	}
}

func Tsv(b bool) Option {
	return func(opts *Options) {
		if b {
			opts.Tsv = b
		}
	}
}

func NoHeaders(b bool) Option {
	return func(opts *Options) {
		if b {
			opts.NoHeaders = b
		}
	}
}

func ApptimeLabel(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.ApptimeLabel = s
		}
	}
}

func ReqtimeLabel(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.ReqtimeLabel = s
		}
	}
}

func StatusLabel(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.StatusLabel = s
		}
	}
}

func SizeLabel(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.SizeLabel = s
		}
	}
}

func MethodLabel(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.MethodLabel = s
		}
	}
}

func UriLabel(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.UriLabel = s
		}
	}
}

func TimeLabel(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.TimeLabel = s
		}
	}
}

func Limit(i int) Option {
	return func(opts *Options) {
		if i > 0 {
			opts.Limit = i
		}
	}
}

func Includes(values []string) Option {
	return func(opts *Options) {
		if len(values) > 0 {
			opts.Includes = values
		}
	}
}

func CSVIncludes(csv string) Option {
	return func(opts *Options) {
		i := splitCSV(csv)
		if len(i) > 0 {
			opts.Includes = i
		}
	}

}

func Excludes(values []string) Option {
	return func(opts *Options) {
		if len(values) > 0 {
			opts.Excludes = values
		}
	}
}

func CSVExcludes(csv string) Option {
	return func(opts *Options) {
		e := splitCSV(csv)
		if len(e) > 0 {
			opts.Excludes = e
		}
	}

}

func Groups(values []string) Option {
	return func(opts *Options) {
		if len(values) > 0 {
			opts.Groups = values
		}
	}
}

func CSVGroups(csv string) Option {
	return func(opts *Options) {
		a := splitCSV(csv)
		if len(a) > 0 {
			opts.Groups = a
		}
	}
}

func StartTime(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.StartTime = s
		}
	}
}

func EndTime(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.EndTime = s
		}
	}
}

func StartTimeDuration(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.StartTimeDuration = s
		}
	}
}

func EndTimeDuration(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.EndTimeDuration = s
		}
	}
}

func Location(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.Location = s
		}
	}
}

func NewOptions(opt ...Option) *Options {
	options := &Options{
		Sort:         DefaultSortOption,
		ApptimeLabel: DefaultApptimeLabelOption,
		ReqtimeLabel: DefaultReqtimeLabelOption,
		StatusLabel:  DefaultStatusLabelOption,
		SizeLabel:    DefaultSizeLabelOption,
		MethodLabel:  DefaultMethodLabelOption,
		UriLabel:     DefaultUriLabelOption,
		TimeLabel:    DefaultTimeLabelOption,
		Limit:        DefaultLimitOption,
	}

	for _, o := range opt {
		o(options)
	}

	return options
}

func SetOptions(options *Options, opt ...Option) *Options {
	for _, o := range opt {
		o(options)
	}

	return options
}

func LoadOptionsFromReader(r io.Reader) (*Options, error) {
	opts := NewOptions()
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return opts, err
	}

	err = yaml.Unmarshal(buf, opts)

	return opts, err
}
