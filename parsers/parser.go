package parsers

import (
	"fmt"
	"net/url"

	"github.com/tkuchiki/alp/errors"
	"github.com/tkuchiki/alp/helpers"
)

type Parser interface {
	Parse() (*ParsedHTTPStat, error)
}

type ParsedHTTPStat struct {
	Uri          string
	Method       string
	Time         string
	ResponseTime float64
	BodyBytes    float64
	Status       int
}

type statKeys struct {
	uri          string
	method       string
	time         string
	responseTime string
	bodyBytes    string
	status       string
}

type statKey func(*statKeys)

func uriKey(s string) statKey {
	return func(sk *statKeys) {
		if s != "" {
			sk.uri = s
		}
	}
}

func methodKey(s string) statKey {
	return func(sk *statKeys) {
		if s != "" {
			sk.method = s
		}
	}
}

func timeKey(s string) statKey {
	return func(sk *statKeys) {
		if s != "" {
			sk.time = s
		}
	}
}

func responseTimeKey(s string) statKey {
	return func(sk *statKeys) {
		if s != "" {
			sk.responseTime = s
		}
	}
}

func bodyBytesKey(s string) statKey {
	return func(sk *statKeys) {
		if s != "" {
			sk.bodyBytes = s
		}
	}
}

func statusKey(s string) statKey {
	return func(sk *statKeys) {
		if s != "" {
			sk.status = s
		}
	}
}

func newStatKeys(sk ...statKey) *statKeys {
	sks := &statKeys{
		uri:          "uri",
		method:       "method",
		time:         "time",
		responseTime: "response_time",
		bodyBytes:    "body_bytes",
		status:       "status",
	}

	for _, s := range sk {
		s(sks)
	}

	return sks
}

func NewParsedHTTPStat(uri, method, time string, resTime, bodyBytes float64, status int) *ParsedHTTPStat {
	return &ParsedHTTPStat{
		Uri:          uri,
		Method:       method,
		Time:         time,
		ResponseTime: resTime,
		BodyBytes:    bodyBytes,
		Status:       status,
	}
}

func toStats(parsedValue map[string]string, keys *statKeys, strictMode, queryString bool) (*ParsedHTTPStat, error) {
	u, err := url.Parse(parsedValue[keys.uri])
	if err != nil {
		return nil, errSkipReadLine(strictMode, err)
	}
	var uri string
	if queryString {
		v := url.Values{}
		values := u.Query()
		for q := range values {
			v.Set(q, "xxx")
		}

		qs := v.Encode()
		if qs != "" {
			uri = fmt.Sprintf("%s?%s", u.Path, qs)
		} else {
			uri = u.Path
		}
	} else {
		uri = u.Path
	}

	resTime, err := helpers.StringToFloat64(parsedValue[keys.responseTime])
	if err != nil {
		return nil, errSkipReadLine(strictMode, err)
	}

	bodyBytes, err := helpers.StringToFloat64(parsedValue[keys.bodyBytes])
	if err != nil {
		return nil, errSkipReadLine(strictMode, err)
	}

	status, err := helpers.StringToInt(parsedValue[keys.status])
	if err != nil {
		return nil, errSkipReadLine(strictMode, err)
	}

	method := parsedValue[keys.method]
	timestr := parsedValue[keys.time]

	return NewParsedHTTPStat(uri, method, timestr, resTime, bodyBytes, status), nil
}

func errSkipReadLine(strictMode bool, err error) error {
	if strictMode {
		return err
	}

	return errors.SkipReadLineErr
}
