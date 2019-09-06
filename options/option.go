package options

import (
	"io"
	"io/ioutil"

	"github.com/tkuchiki/alp/helpers"

	"gopkg.in/yaml.v2"
)

const (
	DefaultSortOption     = "count"
	DefaultFormatOption   = "table"
	DefaultLimitOption    = 5000
	DefaultLocationOption = "Local"
	DefaultOutputOption   = "all"
	// ltsv
	DefaultApptimeLabelOption = "apptime"
	DefaultReqtimeLabelOption = "reqtime"
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
	DefaultRequestTimeKeyOption  = "request_time"
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
		`"(?:.+)"` + // user agent
		`\s(?P<response_time>\S+)(?:\s(?P<request_time>\S+))?$`
	DefaultUriSubexpOption          = "uri"
	DefaultMethodSubexpOption       = "method"
	DefaultTimeSubexpOption         = "time"
	DefaultResponseTimeSubexpOption = "response_time"
	DefaultRequestTimeSubexpOption  = "request_time"
	DefaultBodyBytesSubexpOption    = "body_bytes"
	DefaultStatusSubexpOption       = "status"
)

type Options struct {
	File           string         `yaml:"file"`
	Sort           string         `yaml:"sort"`
	Reverse        bool           `yaml:"reverse"`
	QueryString    bool           `yaml:"query_string"`
	Format         string         `yaml:"format"`
	NoHeaders      bool           `yaml:"noheaders"`
	ShowFooters    bool           `yaml:"show_footers"`
	Limit          int            `yaml:"limit"`
	MatchingGroups []string       `yaml:"matching_groups"`
	Filters        string         `yaml:"filters"`
	PosFile        string         `yaml:"pos_file"`
	NoSavePos      bool           `yaml:"nosave_pos"`
	Location       string         `yaml:"location"`
	Output         string         `yaml:"output"`
	LTSV           *LTSVOptions   `yaml:"ltsv"`
	Regexp         *RegexpOptions `yaml:"regexp"`
	JSON           *JSONOptions   `yaml:"json"`
}

type LTSVOptions struct {
	ApptimeLabel string `yaml:"apptime_label"`
	ReqtimeLabel string `yaml:"reqtime_label"`
	StatusLabel  string `yaml:"status_label"`
	SizeLabel    string `yaml:"size_label"`
	MethodLabel  string `yaml:"method_label"`
	UriLabel     string `yaml:"uri_label"`
	TimeLabel    string `yaml:"time_label"`
}

type RegexpOptions struct {
	Pattern            string `yaml:"pattern"`
	UriSubexp          string `yaml:"uri_subexp"`
	MethodSubexp       string `yaml:"method_subexp"`
	TimeSubexp         string `yaml:"time_subexp"`
	ResponseTimeSubexp string `yaml:"response_time_subexp"`
	RequestTimeSubexp  string `yaml:"request_time_subexp"`
	BodyBytesSubexp    string `yaml:"body_bytes_subexp"`
	StatusSubexp       string `yaml:"status_subexp"`
}

type JSONOptions struct {
	UriKey          string `yaml:"uri_key"`
	MethodKey       string `yaml:"method_key"`
	TimeKey         string `yaml:"time_key"`
	ResponseTimeKey string `yaml:"response_time_key"`
	RequestTimeKey  string `yaml:"request_time_key"`
	BodyBytesKey    string `yaml:"body_bytes_key"`
	StatusKey       string `yaml:"status_key"`
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

func Format(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.Format = s
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

func ShowFooters(b bool) Option {
	return func(opts *Options) {
		if b {
			opts.ShowFooters = b
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

func PosFile(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.PosFile = s
		}
	}
}

func NoSavePos(b bool) Option {
	return func(opts *Options) {
		if b {
			opts.NoSavePos = b
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

func ReqtimeLabel(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.LTSV.ReqtimeLabel = s
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

func RequestTimeSubexp(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.Regexp.RequestTimeSubexp = s
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

func RequestTimeKey(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.JSON.RequestTimeKey = s
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
		ReqtimeLabel: DefaultReqtimeLabelOption,
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
		RequestTimeSubexp:  DefaultRequestTimeSubexpOption,
		BodyBytesSubexp:    DefaultBodyBytesSubexpOption,
		StatusSubexp:       DefaultStatusSubexpOption,
	}

	json := &JSONOptions{
		UriKey:          DefaultUriKeyOption,
		MethodKey:       DefaultMethodKeyOption,
		TimeKey:         DefaultTimeKeyOption,
		ResponseTimeKey: DefaultResponseTimeKeyOption,
		RequestTimeKey:  DefaultRequestTimeKeyOption,
		BodyBytesKey:    DefaultBodyBytesKeyOption,
		StatusKey:       DefaultStatusKeyOption,
	}

	options := &Options{
		Sort:     DefaultSortOption,
		Format:   DefaultFormatOption,
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

	configs := &Options{}
	err = yaml.Unmarshal(buf, configs)

	opts = SetOptions(opts,
		Sort(configs.Sort),
		Limit(configs.Limit),
		Location(configs.Location),
		Output(configs.Output),
		Reverse(configs.Reverse),
		File(configs.File),
		QueryString(configs.QueryString),
		Format(configs.Format),
		NoHeaders(configs.NoHeaders),
		ShowFooters(configs.ShowFooters),
		PosFile(configs.PosFile),
		NoSavePos(configs.NoSavePos),
		MatchingGroups(configs.MatchingGroups),
		Filters(configs.Filters),
		// ltsv
		ApptimeLabel(configs.LTSV.ApptimeLabel),
		ReqtimeLabel(configs.LTSV.ReqtimeLabel),
		StatusLabel(configs.LTSV.StatusLabel),
		SizeLabel(configs.LTSV.SizeLabel),
		MethodLabel(configs.LTSV.MethodLabel),
		UriLabel(configs.LTSV.UriLabel),
		TimeLabel(configs.LTSV.TimeLabel),
		// json
		ResponseTimeKey(configs.JSON.ResponseTimeKey),
		RequestTimeKey(configs.JSON.RequestTimeKey),
		StatusKey(configs.JSON.StatusKey),
		BodyBytesKey(configs.JSON.BodyBytesKey),
		MethodKey(configs.JSON.MethodKey),
		UriKey(configs.JSON.UriKey),
		TimeKey(configs.JSON.TimeKey),
		// regexp
		Pattern(configs.Regexp.Pattern),
		ResponseTimeSubexp(configs.Regexp.ResponseTimeSubexp),
		RequestTimeSubexp(configs.Regexp.RequestTimeSubexp),
		StatusSubexp(configs.Regexp.StatusSubexp),
		BodyBytesSubexp(configs.Regexp.BodyBytesSubexp),
		MethodSubexp(configs.Regexp.MethodSubexp),
		UriSubexp(configs.Regexp.UriSubexp),
		TimeSubexp(configs.Regexp.TimeSubexp),
	)

	return opts, err
}
