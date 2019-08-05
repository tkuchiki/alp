package parsers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
)

type JSONParser struct {
	scanner     *bufio.Scanner
	keys        *statKeys
	strictMode  bool
	queryString bool
}

func NewJSONKeys(uri, method, time, responseTime, size, status string) *statKeys {
	return newStatKeys(
		uriKey(uri),
		methodKey(method),
		timeKey(time),
		responseTimeKey(responseTime),
		bodyBytesKey(size),
		statusKey(status),
	)
}

func NewJSONParser(r io.Reader, keys *statKeys, query bool) Parser {
	return &JSONParser{
		scanner:     bufio.NewScanner(r),
		keys:        keys,
		queryString: query,
	}
}

func (j *JSONParser) Parse() (*ParsedHTTPStat, error) {
	if !j.scanner.Scan() {
		return nil, io.EOF
	}

	var tmp map[string]interface{}
	err := json.Unmarshal([]byte(j.scanner.Text()), &tmp)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 6)
	keys = []string{
		j.keys.uri,
		j.keys.method,
		j.keys.time,
		j.keys.responseTime,
		j.keys.bodyBytes,
		j.keys.status,
	}
	parsedValue := make(map[string]string, 6)
	for _, key := range keys {
		parsedValue[key] = fmt.Sprintf("%v", tmp[key])
	}

	return toStats(parsedValue, j.keys, j.strictMode, j.queryString)
}
