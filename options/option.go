package options

import (
	"io"
	"io/ioutil"

	"github.com/tkuchiki/alp/helpers"

	"gopkg.in/yaml.v2"
)

const (
	DefaultSortOption     = "max"
	DefaultLimitOption    = 5000
	DefaultLocationOption = "Local"
	DefaultOutputOption   = "all"
	// ltsv
	DefaultApptimeLabelOption = "apptime"
	DefaultStatusLabelOption  = "status"
	DefaultSizeLabelOption    = "size"
	DefaultMethodLabelOption  = "method"
	DefaultUriLabelOption     = "uri"
	DefaultTimeLabelOption    = "time"
	// json
	DefaultUriKeyOption          = "uri"
	DefaultMethodKeyOption       = "method"
	DefaultTimeKeyOption         = "time"
	DefaultResponseTimeKeyOption = "response_time"
	DefaultBodyBytesKeyOption    = "body_bytes"
	DefaultStatusKeyOption       = "status"
	// regexp
	DefaultPatternOption = `^(\S+)\s` + // remote host
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
	DefaultUriSubexpOption          = "uri"
	DefaultMethodSubexpOption       = "method"
	DefaultTimeSubexpOption         = "time"
	DefaultResponseTimeSubexpOption = "response_time"
	DefaultBodyBytesSubexpOption    = "body_bytes"
	DefaultStatusSubexpOption       = "status"
)

type Options struct {
	File              string         `yaml:"file"`
	Sort              string         `yaml:"sort"`
	Reverse           bool           `yaml:"reverse"`
	QueryString       bool           `yaml:"query_string"`
	Tsv               bool           `yaml:"tsv"`
	NoHeaders         bool           `yaml:no_headers`
	Limit             int            `yaml:"limit"`
	MatchingGroups    []string       `yaml:"matching_groups"`
	StartTime         string         `yaml:"start_time"`
	EndTime           string         `yaml:"end_time"`
	StartTimeDuration string         `yaml:"start_time_duration"`
	EndTimeDuration   string         `yaml:"end_time_duration"`
	Filters           string         `yaml:"filters"`
	Location          string         `yaml:location`
	Output            string         `yaml:output`
	LTSV              *LTSVOptions   `yaml:ltsv`
	Regexp            *RegexpOptions `yaml:regexp`
	JSON              *JSONOptions   `yaml:json`
}

type LTSVOptions struct {
	ApptimeLabel string `yaml:"apptime_label"`
	StatusLabel  string `yaml:"status_label"`
	SizeLabel    string `yaml:"size_label"`
	MethodLabel  string `yaml:"method_label"`
	UriLabel     string `yaml:"uri_label"`
	TimeLabel    string `yaml:"time_label"`
}

type RegexpOptions struct {
	Pattern            string `yaml:pattern`
	UriSubexp          string `yaml:uri_subexp`
	MethodSubexp       string `yaml:method_subexp`
	TimeSubexp         string `yaml:time_subexp`
	ResponseTimeSubexp string `yaml:response_time_subexp`
	BodyBytesSubexp    string `yaml:body_bytes_subexp`
	StatusSubexp       string `yaml:status_subexp`
}

type JSONOptions struct {
	UriKey          string `yaml:uri_key`
	MethodKey       string `yaml:method_key`
	TimeKey         string `yaml:time_key`
	ResponseTimeKey string `yaml:response_time_key`
	BodyBytesKey    string `yaml:body_bytes_key`
	StatusKey       string `yaml:status_key`
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

func Limit(i int) Option {
	return func(opts *Options) {
		if i > 0 {
			opts.Limit = i
		}
	}
}

func MatchingGroups(values []string) Option {
	return func(opts *Options) {
		if len(values) > 0 {
			opts.MatchingGroups = values
		}
	}
}

func CSVGroups(csv string) Option {
	return func(opts *Options) {
		a := helpers.SplitCSV(csv)
		if len(a) > 0 {
			opts.MatchingGroups = a
		}
	}
}

func Filters(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.Filters = s
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

func Output(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.Output = s
		}
	}
}

// ltsv
func ApptimeLabel(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.LTSV.ApptimeLabel = s
		}
	}
}

func StatusLabel(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.LTSV.StatusLabel = s
		}
	}
}

func SizeLabel(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.LTSV.SizeLabel = s
		}
	}
}

func MethodLabel(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.LTSV.MethodLabel = s
		}
	}
}

func UriLabel(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.LTSV.UriLabel = s
		}
	}
}

func TimeLabel(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.LTSV.TimeLabel = s
		}
	}
}

// regexp
func Pattern(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.Regexp.Pattern = s
		}
	}
}

func UriSubexp(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.Regexp.UriSubexp = s
		}
	}
}

func MethodSubexp(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.Regexp.MethodSubexp = s
		}
	}
}

func TimeSubexp(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.Regexp.TimeSubexp = s
		}
	}
}

func ResponseTimeSubexp(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.Regexp.ResponseTimeSubexp = s
		}
	}
}

func BodyBytesSubexp(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.Regexp.BodyBytesSubexp = s
		}
	}
}

func StatusSubexp(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.Regexp.StatusSubexp = s
		}
	}
}

// json
func UriKey(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.JSON.UriKey = s
		}
	}
}

func MethodKey(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.JSON.MethodKey = s
		}
	}
}

func TimeKey(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.JSON.TimeKey = s
		}
	}
}

func ResponseTimeKey(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.JSON.ResponseTimeKey = s
		}
	}
}

func BodyBytesKey(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.JSON.BodyBytesKey = s
		}
	}
}

func StatusKey(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.JSON.StatusKey = s
		}
	}
}

func NewOptions(opt ...Option) *Options {
	ltsv := &LTSVOptions{
		ApptimeLabel: DefaultApptimeLabelOption,
		StatusLabel:  DefaultStatusLabelOption,
		SizeLabel:    DefaultSizeLabelOption,
		MethodLabel:  DefaultMethodLabelOption,
		UriLabel:     DefaultUriLabelOption,
		TimeLabel:    DefaultTimeLabelOption,
	}

	regexp := &RegexpOptions{
		Pattern:            DefaultPatternOption,
		UriSubexp:          DefaultUriSubexpOption,
		MethodSubexp:       DefaultMethodSubexpOption,
		TimeSubexp:         DefaultTimeSubexpOption,
		ResponseTimeSubexp: DefaultResponseTimeSubexpOption,
		BodyBytesSubexp:    DefaultBodyBytesSubexpOption,
		StatusSubexp:       DefaultStatusSubexpOption,
	}

	json := &JSONOptions{
		UriKey:          DefaultUriKeyOption,
		MethodKey:       DefaultMethodKeyOption,
		TimeKey:         DefaultTimeKeyOption,
		ResponseTimeKey: DefaultResponseTimeKeyOption,
		BodyBytesKey:    DefaultBodyBytesKeyOption,
		StatusKey:       DefaultStatusKeyOption,
	}

	options := &Options{
		Sort:     DefaultSortOption,
		Limit:    DefaultLimitOption,
		Location: DefaultLocationOption,
		Output:   DefaultOutputOption,
		LTSV:     ltsv,
		Regexp:   regexp,
		JSON:     json,
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
