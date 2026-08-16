package parsers

import (
	"bufio"
	"errors"
	"io"
	"regexp"

	httpv1 "github.com/tkuchiki/logschema/http/v1"
	"github.com/tkuchiki/parsetime"
)

type RegexpParser struct {
	reader         *bufio.Reader
	subexpNames    *statKeys
	strictMode     bool
	queryString    bool
	qsIgnoreValues bool
	re             *regexp.Regexp
	readBytes      int
	parseTime      parsetime.ParseTime
}

var errPatternNotMatched = errors.New("pattern not matched")

func NewSubexpNames(uri, method, time, responseTime, requestTime, size, status string) *statKeys {
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

func NewRegexpParser(r io.Reader, expr string, names *statKeys, query, qsIgnoreValues bool, location string) (Parser, error) {
	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	parseTime, err := parsetime.NewParseTime(location)
	if err != nil {
		return nil, err
	}

	return &RegexpParser{
		reader:         bufio.NewReader(r),
		re:             re,
		subexpNames:    names,
		queryString:    query,
		qsIgnoreValues: qsIgnoreValues,
		parseTime:      parseTime,
	}, nil
}

func (rp *RegexpParser) Parse() (*httpv1.Request, error) {
	b, i, err := readline(rp.reader)
	if len(b) == 0 && err != nil {
		return nil, err
	}
	rp.readBytes += i

	groups := rp.re.FindStringSubmatch(string(b))
	if len(groups) == 0 {
		return nil, errSkipReadLine(rp.strictMode, errPatternNotMatched)
	}

	parsedValue := make(map[string]string, len(groups))
	names := rp.re.SubexpNames()
	for i := 1; i < len(groups); i++ {
		if names[i] == "" {
			continue
		}
		parsedValue[names[i]] = groups[i]
	}

	attributes := attributesFromStrings(parsedValue)

	return toHTTPRecord(parsedValue, attributes, rp.subexpNames, rp.strictMode, rp.queryString, rp.qsIgnoreValues, rp.parseTime)
}

func (rp *RegexpParser) ReadBytes() int {
	return rp.readBytes
}

func (rp *RegexpParser) SetReadBytes(n int) {
	rp.readBytes = n
}

func (rp *RegexpParser) Seek(n int) error {
	_, err := rp.reader.Discard(n)
	return err
}
