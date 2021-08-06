package parsers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
)

type JSONParser struct {
	reader         *bufio.Reader
	keys           *statKeys
	strictMode     bool
	queryString    bool
	qsIgnoreValues bool
	readBytes      int
}

func NewJSONKeys(uri, method, time, responseTime, requestTime, size, status string) *statKeys {
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

func NewJSONParser(r io.Reader, keys *statKeys, query, qsIgnoreValues bool) Parser {
	return &JSONParser{
		reader:         bufio.NewReader(r),
		keys:           keys,
		queryString:    query,
		qsIgnoreValues: qsIgnoreValues,
	}
}

func (j *JSONParser) Parse() (*ParsedHTTPStat, error) {
	b, i, err := readline(j.reader)
	if len(b) == 0 && err != nil {
		return nil, err
	}
	j.readBytes += i

	var tmp map[string]interface{}
	err = json.Unmarshal(b, &tmp)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 6)
	keys = []string{
		j.keys.uri,
		j.keys.method,
		j.keys.time,
		j.keys.responseTime,
		j.keys.requestTime,
		j.keys.bodyBytes,
		j.keys.status,
	}
	parsedValue := make(map[string]string, 6)
	for _, key := range keys {
		val, ok := tmp[key]
		if !ok {
			continue
		}

		parsedValue[key] = fmt.Sprintf("%v", val)
	}

	return toStats(parsedValue, j.keys, j.strictMode, j.queryString, j.qsIgnoreValues)
}

func (j *JSONParser) ReadBytes() int {
	return j.readBytes
}

func (j *JSONParser) SetReadBytes(n int) {
	j.readBytes = n
}

func (j *JSONParser) Seek(n int) error {
	_, err := j.reader.Discard(n)
	return err
}
