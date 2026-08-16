package parsers

import (
	"bufio"
	"io"

	"github.com/Songmu/go-ltsv"
	httpv1 "github.com/tkuchiki/logschema/http/v1"
	"github.com/tkuchiki/parsetime"
)

type LTSVParser struct {
	reader         *bufio.Reader
	label          *statKeys
	strictMode     bool
	queryString    bool
	qsIgnoreValues bool
	readBytes      int
	parseTime      parsetime.ParseTime
}

func NewLTSVLabel(uri, method, time, responseTime, requestTime, size, status string) *statKeys {
	return newStatKeys(
		uriKey(uri),
		methodKey(method),
		timeKey(time),
		responseTimeKey(responseTime),
		requestTimeKey(requestTime),
		bodyBytesKey(size),
		statusKey(status),
	)
}

func NewLTSVParser(r io.Reader, l *statKeys, query, qsIgnoreValues bool, location string) (Parser, error) {
	parseTime, err := parsetime.NewParseTime(location)
	if err != nil {
		return nil, err
	}

	return &LTSVParser{
		reader:         bufio.NewReader(r),
		label:          l,
		queryString:    query,
		qsIgnoreValues: qsIgnoreValues,
		parseTime:      parseTime,
	}, nil
}

func (l *LTSVParser) Parse() (*httpv1.Request, error) {
	b, i, err := readline(l.reader)
	if len(b) == 0 && err != nil {
		return nil, err
	}
	l.readBytes += i

	parsedValue := make(map[string]string, 0)
	if err := ltsv.Unmarshal(b, &parsedValue); err != nil && l.strictMode {
		return nil, err
	}

	attributes := attributesFromStrings(parsedValue)

	return toHTTPRecord(parsedValue, attributes, l.label, l.strictMode, l.queryString, l.qsIgnoreValues, l.parseTime)
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
