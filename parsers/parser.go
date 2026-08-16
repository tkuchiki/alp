package parsers

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"net/url"
	"time"

	"github.com/tkuchiki/alp/errors"
	"github.com/tkuchiki/alp/helpers"
	corev1 "github.com/tkuchiki/logschema/core/v1"
	httpv1 "github.com/tkuchiki/logschema/http/v1"
	"github.com/tkuchiki/parsetime"
)

const originalTimeAttribute = "alp.time.original"

type Parser interface {
	Parse() (*httpv1.Request, error)
	ReadBytes() int
	SetReadBytes(n int)
	Seek(n int) error
}

type statKeys struct {
	uri          string
	method       string
	time         string
	responseTime string
	requestTime  string
	bodyBytes    string
	status       string
}

type statKey func(*statKeys)

func uriKey(s string) statKey {
	return func(sk *statKeys) {
		if s != "" {
			sk.uri = s
		}
	}
}

func methodKey(s string) statKey {
	return func(sk *statKeys) {
		if s != "" {
			sk.method = s
		}
	}
}

func timeKey(s string) statKey {
	return func(sk *statKeys) {
		if s != "" {
			sk.time = s
		}
	}
}

func responseTimeKey(s string) statKey {
	return func(sk *statKeys) {
		if s != "" {
			sk.responseTime = s
		}
	}
}

func requestTimeKey(s string) statKey {
	return func(sk *statKeys) {
		if s != "" {
			sk.requestTime = s
		}
	}
}

func bodyBytesKey(s string) statKey {
	return func(sk *statKeys) {
		if s != "" {
			sk.bodyBytes = s
		}
	}
}

func statusKey(s string) statKey {
	return func(sk *statKeys) {
		if s != "" {
			sk.status = s
		}
	}
}

func newStatKeys(sk ...statKey) *statKeys {
	sks := &statKeys{
		uri:          "uri",
		method:       "method",
		time:         "time",
		responseTime: "response_time",
		requestTime:  "request_time",
		bodyBytes:    "body_bytes",
		status:       "status",
	}

	for _, s := range sk {
		s(sks)
	}

	return sks
}

func readline(reader *bufio.Reader) ([]byte, int, error) {
	var b []byte
	var i int
	var err error
	for {
		line, _err := reader.ReadBytes('\n')
		if _err == io.EOF && len(line) == 0 {
			err = io.EOF
			break
		}

		if _err != io.EOF && _err != nil {
			return []byte{}, 0, err
		}
		trimmedLine := bytes.TrimRight(line, "\r\n")
		if len(trimmedLine) > 0 {
			b = append(b, trimmedLine...)
		} else {
			err = errors.SkipReadLineErr
		}

		size := len(line)
		i += size

		if line[size-1] == byte('\n') {
			break
		}
	}

	return b, i, err
}

func toHTTPRecord(
	parsedValue map[string]string,
	attributes corev1.Attributes,
	keys *statKeys,
	strictMode bool,
	queryString bool,
	qsIgnoreValues bool,
	parseTime parsetime.ParseTime,
) (*httpv1.Request, error) {
	u, err := url.Parse(parsedValue[keys.uri])
	if err != nil {
		return nil, errSkipReadLine(strictMode, err)
	}

	normalizedURL := normalizeURL(u, queryString, qsIgnoreValues)
	if normalizedURL.String() == "" {
		return nil, errSkipReadLine(strictMode, fmt.Errorf("URI must not be empty"))
	}

	resTime, err := helpers.StringToFloat64(parsedValue[keys.responseTime])
	if err != nil {
		resTime, err = helpers.StringToFloat64(parsedValue[keys.requestTime])
		if err != nil {
			return nil, errSkipReadLine(strictMode, err)
		}
	}

	bodyBytes, err := helpers.StringToFloat64(parsedValue[keys.bodyBytes])
	if err != nil {
		return nil, errSkipReadLine(strictMode, err)
	}

	status, err := helpers.StringToInt(parsedValue[keys.status])
	if err != nil {
		return nil, errSkipReadLine(strictMode, err)
	}

	method := parsedValue[keys.method]

	rawTime := parsedValue[keys.time]
	var eventTime *corev1.DecimalInt64
	if rawTime != "" {
		parsedTime, err := parseTime.Parse(rawTime)
		if err != nil {
			return nil, errSkipReadLine(strictMode, err)
		}

		unixNano := corev1.DecimalInt64(parsedTime.UnixNano())
		eventTime = &unixNano
	}

	record, err := newHTTPRequest(normalizedURL, method, eventTime, resTime, bodyBytes, status, attributes)
	if err != nil {
		return nil, errSkipReadLine(strictMode, err)
	}
	setOriginalTime(record, rawTime)

	return record, nil
}

// OriginalTime returns the producer-preserved timestamp representation used by
// ALP's raw-entry views.
func OriginalTime(request *httpv1.Request) string {
	value, ok := request.Data.Attributes[originalTimeAttribute].(string)
	if !ok {
		return ""
	}

	return value
}

func setOriginalTime(request *httpv1.Request, value string) {
	if value == "" {
		return
	}
	if request.Data.Attributes == nil {
		request.Data.Attributes = make(corev1.Attributes)
	}

	request.Data.Attributes[originalTimeAttribute] = value
}

func normalizeURL(src *url.URL, queryString, qsIgnoreValues bool) *url.URL {
	u := *src
	if src.RawQuery == "" {
		return &u
	}

	if !queryString {
		u.RawQuery = ""
		return &u
	}

	if qsIgnoreValues {
		values := u.Query()
		for q := range values {
			values.Set(q, "xxx")
		}
		u.RawQuery = values.Encode()
	} else {
		u.RawQuery = u.Query().Encode() // re-encode to sort queries
	}

	return &u
}

func newHTTPRequest(
	u *url.URL,
	method string,
	eventTime *corev1.DecimalInt64,
	responseTimeSeconds float64,
	responseBodySize float64,
	statusCode int,
	attributes corev1.Attributes,
) (*httpv1.Request, error) {
	durationNano, err := secondsToDurationNano(responseTimeSeconds)
	if err != nil {
		return nil, err
	}

	var bodySizeBytes *corev1.DecimalUint64
	if responseBodySize != -1 {
		converted, err := floatToUint64(responseBodySize, "response body size")
		if err != nil {
			return nil, err
		}

		bodySizeBytes = &converted
	}

	data := requestDataFromURL(u)
	data.Method = method
	data.StatusCode = &statusCode
	data.ResponseBodySizeBytes = bodySizeBytes
	data.Attributes = attributes

	request, err := httpv1.NewRequest(
		durationNano,
		corev1.Source{Kind: corev1.SourceOther},
		data,
	)
	if err != nil {
		return nil, err
	}
	request.TimeUnixNano = eventTime

	return &request, nil
}

func requestDataFromURL(u *url.URL) httpv1.RequestData {
	path := u.EscapedPath()
	if u.Opaque != "" && u.Host == "" {
		path = u.Scheme + ":" + u.Opaque
	}
	if path == "" && u.Host != "" {
		path = "/"
	}

	data := httpv1.RequestData{URLPath: path}
	if u.Host != "" {
		data.URLAuthority = stringPointer(u.Host)
		if u.Scheme != "" {
			data.URLScheme = stringPointer(u.Scheme)
		}
	}
	if u.RawQuery != "" || u.ForceQuery {
		data.URLQuery = stringPointer(u.RawQuery)
	}

	return data
}

func secondsToDurationNano(seconds float64) (corev1.DecimalUint64, error) {
	if math.IsNaN(seconds) || math.IsInf(seconds, 0) || seconds < 0 {
		return 0, fmt.Errorf("response time must be a finite non-negative number")
	}

	nanoseconds := math.Round(seconds * float64(time.Second))
	if nanoseconds >= math.Exp2(64) {
		return 0, fmt.Errorf("response time exceeds the LogSchema duration range")
	}

	return corev1.DecimalUint64(uint64(nanoseconds)), nil
}

func floatToUint64(value float64, field string) (corev1.DecimalUint64, error) {
	if math.IsNaN(value) || math.IsInf(value, 0) || value < 0 || math.Trunc(value) != value {
		return 0, fmt.Errorf("%s must be a finite non-negative integer", field)
	}
	if value >= math.Exp2(64) {
		return 0, fmt.Errorf("%s exceeds the LogSchema uint64 range", field)
	}

	return corev1.DecimalUint64(uint64(value)), nil
}

func stringPointer(value string) *string {
	return &value
}

func attributesFromStrings(values map[string]string) corev1.Attributes {
	attributes := make(corev1.Attributes, len(values))
	for key, value := range values {
		attributes[key] = value
	}

	return attributes
}

func errSkipReadLine(strictMode bool, err error) error {
	if strictMode {
		return err
	}

	return errors.SkipReadLineErr
}
