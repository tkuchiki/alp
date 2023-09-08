package testutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
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

func CreateTempDirAndFile(dir, content string) (string, error) {
	fpath := filepath.Join(dir, fmt.Sprint(time.Now().UnixNano()))
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

	return r.Replace(`{"__time__":"2015-09-06T05:58:05+09:00","__method__":"POST","__uri__":"/foo/bar/123?token=xxx&uuid=1234","__status__":200,"__body_bytes__":12,"__response_time__":0.057,"__request_time__":0.057}
{"__time__":"2015-09-06T05:58:05+09:00","__method__":"POST","__uri__":"/foo/bar/456?token=yyy","__status__":200,"__body_bytes__":34,"__response_time__":0.100,"__request_time__":0.100}`)
}

func LTSVLog(keys LogKeys) string {
	r := strings.NewReplacer("__time__", keys.Time,
		"__method__", keys.Method,
		"__uri__", keys.Uri,
		"__status__", keys.Status,
		"__size__", keys.BodyBytes,
		"__apptime__", keys.ResponseTime,
		"__reqtime__", keys.RequestTime)

	return r.Replace(`__time__:2015-09-06T05:58:05+09:00	__method__:POST	__uri__:/foo/bar/123?token=xxx&uuid=1234	__status__:200	__size__:12	__apptime__:0.057	__reqtime__:0.0057
__time__:2015-09-06T05:58:41+09:00	__method__:POST	__uri__:/foo/bar/456?token=yyy	__status__:200	__size__:34	__apptime__:0.100	__reqtime__:0.100
`)
}

func RegexpLog() string {
	return `127.0.0.1 - - [06/Sep/2015:05:58:05 +0900] "POST /foo/bar/123?token=xxx&uuid=1234 HTTP/1.1" 200 12 "-" "curl/7.54.0" "-" 0.057
127.0.0.1 - - [06/Sep/2015:05:58:41 +0900] "POST /foo/bar/456?token=yyy HTTP/1.1" 200 34 "-" "curl/7.54.0" "-" 0.100`
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
		`"(?:.+)"` + // user agent
		`\s(?P<__response_time__>\S+)(?:\s(?P<__request_time__>\S+))?$`)
}

func ConfigFile() string {
	return `sort: max
reverse: true
query_string: true`
}
