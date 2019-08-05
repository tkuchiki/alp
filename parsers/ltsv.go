package parsers

import (
	"io"

	"github.com/najeira/ltsv"
)

type LTSVParser struct {
	reader      *ltsv.Reader
	label       *statKeys
	strictMode  bool
	queryString bool
}

func NewLTSVLabel(uri, method, time, responseTime, size, status string) *statKeys {
	return newStatKeys(
		uriKey(uri),
		methodKey(method),
		timeKey(time),
		responseTimeKey(responseTime),
		bodyBytesKey(size),
		statusKey(status),
	)
}

func NewLTSVParser(r io.Reader, l *statKeys, query bool) Parser {
	return &LTSVParser{
		reader:      ltsv.NewReader(r),
		label:       l,
		queryString: query,
	}
}

func (l *LTSVParser) Parse() (*ParsedHTTPStat, error) {
	parsedValue, err := l.reader.Read()
	if err != nil && l.strictMode {
		return nil, err
	} else if err == io.EOF {
		return nil, err
	}

	return toStats(parsedValue, l.label, l.strictMode, l.queryString)
}
