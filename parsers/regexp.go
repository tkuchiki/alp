package parsers

import (
	"bufio"
	"errors"
	"io"
	"regexp"
)

type RegexpParser struct {
	scanner     *bufio.Scanner
	subexpNames *statKeys
	strictMode  bool
	queryString bool
	re          *regexp.Regexp
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
		scanner:     bufio.NewScanner(r),
		re:          re,
		subexpNames: names,
		queryString: query,
	}, nil
}

func (rp *RegexpParser) Parse() (*ParsedHTTPStat, error) {
	if !rp.scanner.Scan() {
		return nil, io.EOF
	}

	groups := rp.re.FindStringSubmatch(rp.scanner.Text())
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
