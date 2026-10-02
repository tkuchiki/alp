package parsers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	corev1 "github.com/tkuchiki/logschema/core/v1"
	httpv1 "github.com/tkuchiki/logschema/http/v1"
	"github.com/tkuchiki/parsetime"
)

type JSONParser struct {
	reader         *bufio.Reader
	keys           *statKeys
	strictMode     bool
	queryString    bool
	qsIgnoreValues bool
	readBytes      int
	parseTime      parsetime.ParseTime
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

func NewJSONParser(r io.Reader, keys *statKeys, query, qsIgnoreValues bool, location string) (Parser, error) {
	parseTime, err := parsetime.NewParseTime(location)
	if err != nil {
		return nil, err
	}

	return &JSONParser{
		reader:         bufio.NewReader(r),
		keys:           keys,
		queryString:    query,
		qsIgnoreValues: qsIgnoreValues,
		parseTime:      parseTime,
	}, nil
}

func (j *JSONParser) Parse() (*httpv1.Request, error) {
	b, i, err := readline(j.reader)
	if len(b) == 0 && err != nil {
		return nil, err
	}
	j.readBytes += i

	var tmp map[string]any
	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.UseNumber()
	err = decoder.Decode(&tmp)
	if err != nil {
		return nil, err
	}

	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("multiple JSON values are not allowed")
		}

		return nil, fmt.Errorf("decode trailing JSON: %w", err)
	}

	keys := []string{
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

	attributes := make(corev1.Attributes, len(tmp))
	for key, val := range tmp {
		if val == nil {
			continue
		}

		attributes[key] = jsonAttributeValue(val)
	}

	return toHTTPRecord(parsedValue, attributes, j.keys, j.strictMode, j.queryString, j.qsIgnoreValues, j.parseTime)
}

func jsonAttributeValue(value any) any {
	switch value := value.(type) {
	case map[string]any:
		return jsonString(value)
	case []any:
		for _, item := range value {
			switch item.(type) {
			case nil, map[string]any, []any:
				return jsonString(value)
			}
		}

		return value
	default:
		return value
	}
}

func jsonString(value any) string {
	encoded, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprint(value)
	}

	return string(encoded)
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
