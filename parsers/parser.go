package parsers

import (
	"github.com/tkuchiki/alp/stats"
)

type Parser interface {
	Parse() *HTTPStat
}

type HTTPStat struct {
	Uri          string
	Method       string
	Time         string
	ResponseTime float64
	BodySize     float64
	Status       int
}

func NewHTTPStat(uri, method, time string, resTime, bodySize float64, status int) *HTTPStat {
	return &HTTPStat{
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

	return stats.SkipReadLineErr
}
