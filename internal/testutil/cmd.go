package testutil

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/tkuchiki/alp/options"
	"github.com/tkuchiki/alp/stats"
)

type LogKeys struct {
	Uri          string
	Method       string
	Time         string
	ResponseTime string
	RequestTime  string
	BodyBytes    string
	Status       string
}

func NewJsonLogKeys() LogKeys {
	return LogKeys{
		Uri:          "uri",
		Method:       "method",
		Time:         "time",
		ResponseTime: "response_time",
		RequestTime:  "request_time",
		BodyBytes:    "body_bytes",
		Status:       "status",
	}
}

func NewLTSVLogKeys() LogKeys {
	return LogKeys{
		Uri:          "uri",
		Method:       "method",
		Time:         "time",
		ResponseTime: "apptime",
		RequestTime:  "reqtime",
		BodyBytes:    "size",
		Status:       "status",
	}
}

func NewRegexpLogKeys() LogKeys {
	return LogKeys{
		Uri:          "uri",
		Method:       "method",
		Time:         "time",
		ResponseTime: "response_time",
		RequestTime:  "request_time",
		BodyBytes:    "body_bytes",
		Status:       "status",
	}
}

func CreateTempDirAndFile(dir, filename, content string) (string, error) {
	fpath := filepath.Join(dir, filename)
	err := os.WriteFile(fpath, []byte(content), 0644)

	return fpath, err
}

func JsonLog(keys LogKeys) string {
	r := strings.NewReplacer("__time__", keys.Time,
		"__method__", keys.Method,
		"__uri__", keys.Uri,
		"__status__", keys.Status,
		"__body_bytes__", keys.BodyBytes,
		"__response_time__", keys.ResponseTime,
		"__request_time__", keys.RequestTime)

	return r.Replace(`{"__time__":"2015-09-06T05:58:05+09:00","__method__":"POST","__uri__":"/foo/bar/123?token=xxx&uuid=1234","__status__":200,"__body_bytes__":12,"__response_time__":0.057,"__request_time__":0.057,"ua":"UA1"}
{"__time__":"2015-09-06T05:58:05+09:00","__method__":"POST","__uri__":"/foo/bar/456?token=yyy","__status__":200,"__body_bytes__":34,"__response_time__":0.100,"__request_time__":0.100,"ua":"UA2"}`)
}

func LTSVLog(keys LogKeys) string {
	r := strings.NewReplacer("__time__", keys.Time,
		"__method__", keys.Method,
		"__uri__", keys.Uri,
		"__status__", keys.Status,
		"__size__", keys.BodyBytes,
		"__apptime__", keys.ResponseTime,
		"__reqtime__", keys.RequestTime)

	return r.Replace(`__time__:2015-09-06T05:58:05+09:00	__method__:POST	__uri__:/foo/bar/123?token=xxx&uuid=1234	__status__:200	__size__:12	__apptime__:0.057	__reqtime__:0.0057	ua:UA1
__time__:2015-09-06T05:58:41+09:00	__method__:POST	__uri__:/foo/bar/456?token=yyy	__status__:200	__size__:34	__apptime__:0.100	__reqtime__:0.100	ua:UA2
`)
}

func RegexpLog() string {
	return `127.0.0.1 - - [06/Sep/2015:05:58:05 +0900] "POST /foo/bar/123?token=xxx&uuid=1234 HTTP/1.1" 200 12 "-" "UA1" "-" 0.057
127.0.0.1 - - [06/Sep/2015:05:58:41 +0900] "POST /foo/bar/456?token=yyy HTTP/1.1" 200 34 "-" "UA2" "-" 0.100`
}

func RegexpPattern(keys LogKeys) string {
	r := strings.NewReplacer("__time__", keys.Time,
		"__method__", keys.Method,
		"__uri__", keys.Uri,
		"__status__", keys.Status,
		"__body_bytes__", keys.BodyBytes,
		"__response_time__", keys.ResponseTime,
		"__request_time__", keys.RequestTime)

	return r.Replace(`\S+\s+` +
		`(\S+\s+)+` + // user
		`\[(?P<__time__>[^]]+)\]\s` + // time
		`"(?P<__method__>\S*)\s?` + // method
		`(?P<__uri__>(?:[^"]*(?:\\")?)*)\s` + // URL
		`([^"]*)"\s` + // protocol
		`(?P<__status__>\S+)\s` + // status code
		`(?P<__body_bytes__>\S+)\s` + // bytes
		`"((?:[^"]*(?:\\")?)*)"\s` + // referer
		`"(?P<ua>[^"]*)"` + // user agent
		`\s(?P<__response_time__>\S+)(?:\s(?P<__request_time__>\S+))?$`)
}

func ConfigFile() string {
	return `sort: max
reverse: true
query_string: true`
}

func DummyOptions(sort string) *options.Options {
	sortOptions := stats.NewSortOptions()
	sortOptions.SetAndValidate(sort)

	return &options.Options{
		File:                    "/path/to/file",
		Sort:                    sortOptions.SortType(),
		Location:                "dummy",
		Reverse:                 false,
		QueryString:             false,
		QueryStringIgnoreValues: false,
		DecodeUri:               false,
		Format:                  "markdown",
		Limit:                   100,
		NoHeaders:               false,
		ShowFooters:             false,
		MatchingGroups: []string{
			"/foo/.+",
		},
		Filters:         ".Uri == '/foo/bar'",
		Output:          "count,uri,min,max",
		PosFile:         "/path/to/pos",
		NoSavePos:       false,
		Percentiles:     []int{1, 5},
		PaginationLimit: 10,
		LTSV: &options.LTSVOptions{
			UriLabel:     "u",
			MethodLabel:  "m",
			TimeLabel:    "t",
			ApptimeLabel: "a",
			ReqtimeLabel: "r",
			SizeLabel:    "sz",
			StatusLabel:  "st",
		},
		JSON: &options.JSONOptions{
			UriKey:          "u",
			MethodKey:       "m",
			TimeKey:         "t",
			ResponseTimeKey: "res",
			RequestTimeKey:  "req",
			BodyBytesKey:    "b",
			StatusKey:       "s",
		},
		Regexp: &options.RegexpOptions{
			Pattern:            "dummy pattern",
			UriSubexp:          "u",
			MethodSubexp:       "m",
			TimeSubexp:         "t",
			ResponseTimeSubexp: "res",
			RequestTimeSubexp:  "req",
			BodyBytesSubexp:    "b",
			StatusSubexp:       "s",
		},
		Pcap: &options.PcapOptions{
			ServerIPs: []string{
				"192.168.1.10",
			},
			ServerPort: 12345,
		},
		Count: &options.CountOptions{
			Keys: []string{
				"ua",
			},
		},
		TopN: &options.TopNOptions{
			Sort:    "restime",
			Reverse: false,
		},
	}
}

func DummyOverwrittenOptions(sort string) *options.Options {
	sortOptions := stats.NewSortOptions()
	sortOptions.SetAndValidate(sort)

	return &options.Options{
		File:                    "/path/to/overwritten/file",
		Sort:                    sortOptions.SortType(),
		Location:                "overwritten location",
		Reverse:                 true,
		QueryString:             true,
		QueryStringIgnoreValues: true,
		DecodeUri:               true,
		Format:                  "tsv",
		Limit:                   200,
		NoHeaders:               true,
		ShowFooters:             true,
		MatchingGroups: []string{
			"/foo/bar/.+",
			"/bar/.+",
		},
		Filters:         ".Status == 200",
		Output:          "uri,avg",
		PosFile:         "/path/to/overwritten/pos",
		NoSavePos:       true,
		Percentiles:     []int{5, 9},
		PaginationLimit: 20,
		LTSV: &options.LTSVOptions{
			UriLabel:     "u2",
			MethodLabel:  "m2",
			TimeLabel:    "t2",
			ApptimeLabel: "a2",
			ReqtimeLabel: "r2",
			SizeLabel:    "sz2",
			StatusLabel:  "st2",
		},
		JSON: &options.JSONOptions{
			UriKey:          "u2",
			MethodKey:       "m2",
			TimeKey:         "t2",
			ResponseTimeKey: "res2",
			RequestTimeKey:  "req2",
			BodyBytesKey:    "b2",
			StatusKey:       "s2",
		},
		Regexp: &options.RegexpOptions{
			UriSubexp:          "u2",
			MethodSubexp:       "m2",
			TimeSubexp:         "t2",
			ResponseTimeSubexp: "res2",
			RequestTimeSubexp:  "req2",
			BodyBytesSubexp:    "b2",
			StatusSubexp:       "s2",
		},
		Pcap: &options.PcapOptions{
			ServerIPs: []string{
				"192.168.1.20",
			},
			ServerPort: 54321,
		},
		Count: &options.CountOptions{
			Keys: []string{
				"host",
				"user_agent",
			},
		},
		TopN: &options.TopNOptions{
			Sort:    "bytes",
			Reverse: true,
		},
	}
}

func DummyConfigFile(sort string, dummyOpts *options.Options) string {
	configTmpl := `file: {{ .File }}
sort: ` + sort + `
reverse: {{ .Reverse }}
query_string: {{ .QueryString }}
query_string_ignore_values: {{ .QueryStringIgnoreValues }}
decode_uri: {{ .DecodeUri }}
format: {{ .Format }}
limit: {{ .Limit }}
noheaders: {{ .NoHeaders }}
show_footers: {{ .ShowFooters }}
matching_groups:
{{ range .MatchingGroups }}
  - {{ . }}
{{ end }}
filters: {{ .Filters }}
output: {{ .Output }}
pos_file: {{ .PosFile }}
nosave_pos: {{ .NoSavePos }}
location: {{ .Location }}
percentiles:
{{ range .Percentiles }}
  - {{ . }}
{{ end }}
pagination_limit: {{ .PaginationLimit }}
ltsv:
  uri_label: {{ .LTSV.UriLabel }}
  method_label: {{ .LTSV.MethodLabel }}
  time_label: {{ .LTSV.TimeLabel }}
  apptime_label: {{ .LTSV.ApptimeLabel }}
  reqtime_label: {{ .LTSV.ReqtimeLabel }}
  size_label: {{ .LTSV.SizeLabel }}
  status_label: {{ .LTSV.StatusLabel }}
json:
  uri_key: {{ .JSON.UriKey }}
  method_key: {{ .JSON.MethodKey }}
  time_key: {{ .JSON.TimeKey }}
  response_time_key: {{ .JSON.ResponseTimeKey }}
  request_time_key: {{ .JSON.RequestTimeKey }}
  body_bytes_key: {{ .JSON.BodyBytesKey }}
  status_key: {{ .JSON.StatusKey }}
regexp:
  pattern: {{ .Regexp.Pattern }}
  uri_subexp: {{ .Regexp.UriSubexp }}
  method_subexp: {{ .Regexp.MethodSubexp }}
  time_subexp: {{ .Regexp.TimeSubexp }}
  response_time_subexp: {{ .Regexp.ResponseTimeSubexp }}
  request_time_subexp: {{ .Regexp.RequestTimeSubexp }}
  body_bytes_subexp: {{ .Regexp.BodyBytesSubexp }}
  status_subexp: {{ .Regexp.StatusSubexp }}
pcap:
  server_ips:
{{ range .Pcap.ServerIPs }}
    - {{ . }}
{{ end }}
  server_port: {{ .Pcap.ServerPort }}
count:
  keys:
{{ range .Count.Keys }}
    - {{ . }}
{{ end }}
topN:
  sort: {{ .TopN.Sort }}
  reverse: {{ .TopN.Reverse }}
`
	t, err := template.New("dummy_config").Parse(configTmpl)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	if err = t.Execute(&buf, dummyOpts); err != nil {
		panic(err)
	}

	return buf.String()
}
