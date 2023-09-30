package options

import (
	"net"

	"github.com/tkuchiki/alp/helpers"
)

const (
	DefaultSortOption      = "count"
	DefaultFormatOption    = "table"
	DefaultLimitOption     = 5000
	DefaultLocationOption  = "Local"
	DefaultOutputOption    = "all"
	DefaultPaginationLimit = 100
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
	// pcap
	DefaultPcapServerPortOption = 80
	// topN
	DefaultTopNSortOption = "restime"
)

var DefaultPercentilesOption = []int{90, 95, 99}

var DefaultPcapServerIPsOption = getDefaultPcapServerIPsOption()

type Options struct {
	File                    string         `mapstructure:"file"`
	Dump                    string         `mapstructure:"dump"`
	Load                    string         `mapstructure:"load"`
	Sort                    string         `mapstructure:"sort"`
	Reverse                 bool           `mapstructure:"reverse"`
	QueryString             bool           `mapstructure:"query_string"`
	QueryStringIgnoreValues bool           `mapstructure:"query_string_ignore_values"`
	DecodeUri               bool           `mapstructure:"decode_uri"`
	Format                  string         `mapstructure:"format"`
	NoHeaders               bool           `mapstructure:"noheaders"`
	ShowFooters             bool           `mapstructure:"show_footers"`
	Limit                   int            `mapstructure:"limit"`
	MatchingGroups          []string       `mapstructure:"matching_groups"`
	Filters                 string         `mapstructure:"filters"`
	PosFile                 string         `mapstructure:"pos_file"`
	NoSavePos               bool           `mapstructure:"nosave_pos"`
	Location                string         `mapstructure:"location"`
	Output                  string         `mapstructure:"output"`
	Percentiles             []int          `mapstructure:"percentiles"`
	PaginationLimit         int            `mapstructure:"pagination_limit"`
	LTSV                    *LTSVOptions   `mapstructure:"ltsv"`
	Regexp                  *RegexpOptions `mapstructure:"regexp"`
	JSON                    *JSONOptions   `mapstructure:"json"`
	Pcap                    *PcapOptions   `mapstructure:"pcap"`
	Count                   *CountOptions  `mapstructure:"count"`
	TopN                    *TopNOptions   `mapstructure:"topN"`
}

type LTSVOptions struct {
	ApptimeLabel string `mapstructure:"apptime_label"`
	ReqtimeLabel string `mapstructure:"reqtime_label"`
	StatusLabel  string `mapstructure:"status_label"`
	SizeLabel    string `mapstructure:"size_label"`
	MethodLabel  string `mapstructure:"method_label"`
	UriLabel     string `mapstructure:"uri_label"`
	TimeLabel    string `mapstructure:"time_label"`
}

type RegexpOptions struct {
	Pattern            string `mapstructure:"pattern"`
	UriSubexp          string `mapstructure:"uri_subexp"`
	MethodSubexp       string `mapstructure:"method_subexp"`
	TimeSubexp         string `mapstructure:"time_subexp"`
	ResponseTimeSubexp string `mapstructure:"response_time_subexp"`
	RequestTimeSubexp  string `mapstructure:"request_time_subexp"`
	BodyBytesSubexp    string `mapstructure:"body_bytes_subexp"`
	StatusSubexp       string `mapstructure:"status_subexp"`
}

type JSONOptions struct {
	UriKey          string `mapstructure:"uri_key"`
	MethodKey       string `mapstructure:"method_key"`
	TimeKey         string `mapstructure:"time_key"`
	ResponseTimeKey string `mapstructure:"response_time_key"`
	RequestTimeKey  string `mapstructure:"request_time_key"`
	BodyBytesKey    string `mapstructure:"body_bytes_key"`
	StatusKey       string `mapstructure:"status_key"`
}

type PcapOptions struct {
	ServerIPs  []string `mapstructure:"server_ips"`
	ServerPort uint16   `mapstructure:"server_port"`
}

type CountOptions struct {
	Keys []string `mapstructure:"keys"`
}

type TopNOptions struct {
	Sort    string `mapstructure:"sort"`
	Reverse bool   `mapstructure:"reverse"`
}

type Option func(*Options)

func File(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.File = s
		}
	}
}

func Dump(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.Dump = s
		}
	}
}

func Load(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.Load = s
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

func QueryStringIgnoreValues(b bool) Option {
	return func(opts *Options) {
		if b {
			opts.QueryStringIgnoreValues = b
		}
	}
}

func DecodeUri(b bool) Option {
	return func(opts *Options) {
		if b {
			opts.DecodeUri = b
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

func Percentiles(i []int) Option {
	return func(opts *Options) {
		if len(i) > 0 {
			opts.Percentiles = i
		}
	}
}

func PaginationLimit(i int) Option {
	return func(opts *Options) {
		if i > 0 {
			opts.PaginationLimit = i
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

// pcap
func PcapServerIPs(ss []string) Option {
	return func(opts *Options) {
		if len(ss) > 0 {
			opts.Pcap.ServerIPs = ss
		}
	}
}

func PcapServerPort(n uint16) Option {
	return func(opts *Options) {
		if n != 0 {
			opts.Pcap.ServerPort = n
		}
	}
}

// count
func CountKeys(ss []string) Option {
	return func(opts *Options) {
		if len(ss) > 0 {
			opts.Count.Keys = ss
		}
	}
}

// topN
func TopNSort(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.TopN.Sort = s
		}
	}
}

func TopNReverse(b bool) Option {
	return func(opts *Options) {
		if b {
			opts.TopN.Reverse = b
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

	pcap := &PcapOptions{
		ServerIPs:  DefaultPcapServerIPsOption,
		ServerPort: DefaultPcapServerPortOption,
	}

	count := &CountOptions{}

	topN := &TopNOptions{
		Sort: DefaultTopNSortOption,
	}

	options := &Options{
		Sort:            DefaultSortOption,
		Format:          DefaultFormatOption,
		Limit:           DefaultLimitOption,
		Location:        DefaultLocationOption,
		Output:          DefaultOutputOption,
		Percentiles:     DefaultPercentilesOption,
		PaginationLimit: DefaultPaginationLimit,
		LTSV:            ltsv,
		Regexp:          regexp,
		JSON:            json,
		Pcap:            pcap,
		Count:           count,
		TopN:            topN,
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

func getDefaultPcapServerIPsOption() (ips []string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		// fallback
		return []string{"127.0.0.1", "::1"}
	}

	for _, addr := range addrs {
		ipNet, isIP := addr.(*net.IPNet)
		if !isIP {
			continue
		}

		ips = append(ips, ipNet.IP.String())
	}
	if len(ips) == 0 {
		// fallback
		return []string{"127.0.0.1", "::1"}
	}

	return
}
