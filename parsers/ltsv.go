package parsers

import (
	"bufio"
	"io"

	"github.com/Songmu/go-ltsv"
)

type LTSVParser struct {
	reader      *bufio.Reader
	label       *statKeys
	strictMode  bool
	queryString bool
	readBytes   int
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
		reader:      bufio.NewReader(r),
		label:       l,
		queryString: query,
	}
}

func (l *LTSVParser) Parse() (*ParsedHTTPStat, error) {
	b, i, err := readline(l.reader)
	if len(b) == 0 && err != nil {
		return nil, err
	}
	l.readBytes += i

	parsedValue := make(map[string]string, 0)
	err2 := ltsv.Unmarshal(b, &parsedValue)
	if err2 != nil && l.strictMode {
		return nil, err
	}

	return toStats(parsedValue, l.label, l.strictMode, l.queryString)
}

func (l *LTSVParser) ReadBytes() int {
	return l.readBytes
}

func (l *LTSVParser) SetReadBytes(n int) {
	l.readBytes = n
}

func (l *LTSVParser) Seek(n int) error {
	_, err := l.reader.Discard(n)
	return err
}
