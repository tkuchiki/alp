package parsers

import (
	"bufio"
	"errors"
	"io"
	"regexp"
)

type RegexpParser struct {
	reader      *bufio.Reader
	subexpNames *statKeys
	strictMode  bool
	queryString bool
	re          *regexp.Regexp
	readBytes   int
}

var errPatternNotMatched = errors.New("pattern not matched")

func NewSubexpNames(uri, method, time, responseTime, size, status string) *statKeys {
	return newStatKeys(
		uriKey(uri),
		methodKey(method),
		timeKey(time),
		responseTimeKey(responseTime),
		bodyBytesKey(size),
		statusKey(status),
	)
}

func NewRegexpParser(r io.Reader, expr string, names *statKeys, query bool) (Parser, error) {
	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}

	return &RegexpParser{
		reader:      bufio.NewReader(r),
		re:          re,
		subexpNames: names,
		queryString: query,
	}, nil
}

func (rp *RegexpParser) Parse() (*ParsedHTTPStat, error) {
	b, i, err := readline(rp.reader)
	if len(b) == 0 && err != nil {
		return nil, err
	}
	rp.readBytes += i

	groups := rp.re.FindStringSubmatch(string(b))
	if len(groups) == 0 {
		return nil, errSkipReadLine(rp.strictMode, errPatternNotMatched)
	}

	parsedValue := make(map[string]string, 6)
	names := rp.re.SubexpNames()
	for i := 1; i < len(groups); i++ {
		if names[i] == "" {
			continue
		}
		parsedValue[names[i]] = groups[i]
	}

	return toStats(parsedValue, rp.subexpNames, rp.strictMode, rp.queryString)
}

func (rp *RegexpParser) ReadBytes() int {
	return rp.readBytes
}

func (rp *RegexpParser) Seek(n int) error {
	_, err := rp.reader.Discard(n)
	return err
}
