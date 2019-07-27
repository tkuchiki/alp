package parsers

import (
	"github.com/tkuchiki/alp/errors"
)

type Parser interface {
	Parse() *ParsedHTTPStat
}

type ParsedHTTPStat struct {
	Uri          string
	Method       string
	Time         string
	ResponseTime float64
	BodySize     float64
	Status       int
}

func NewParsedHTTPStat(uri, method, time string, resTime, bodySize float64, status int) *ParsedHTTPStat {
	return &ParsedHTTPStat{
		Uri:          uri,
		Method:       method,
		Time:         time,
		ResponseTime: resTime,
		BodySize:     bodySize,
		Status:       status,
	}
}

func errSkipReadLine(strictMode bool, err error) error {
	if strictMode {
		return err
	}

	return errors.SkipReadLineErr
}
