package stats

import (
	"time"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	httpv1 "github.com/tkuchiki/logschema/http/v1"
	"github.com/tkuchiki/parsetime"
)

type ExpEval struct {
	program   *vm.Program
	parseTime parsetime.ParseTime
}

type ExpEvalEnv struct {
	Uri                              string
	Method                           string
	Time                             time.Time
	ResponseTime                     float64
	BodyBytes                        float64
	Status                           int
	TimeStringEqualTime              func(l time.Time, r string) bool
	TimeStringNotEqualTime           func(l time.Time, r string) bool
	TimeStringGreaterThanTime        func(l time.Time, r string) bool
	TimeStringGreaterThanOrEqualTime func(l time.Time, r string) bool
	TimeStringLessThanTime           func(l time.Time, r string) bool
	TimeStringLessThanOrEqualTime    func(l time.Time, r string) bool
	TimeEqualTime                    func(l, r time.Time) bool
	TimeNotEqualTime                 func(l, r time.Time) bool
	TimeGreaterThanTime              func(l, r time.Time) bool
	TimeGreaterThanOrEqualTime       func(l, r time.Time) bool
	TimeLessThanTime                 func(l, r time.Time) bool
	TimeLessThanOrEqualTime          func(l, r time.Time) bool
	TimeAgo                          func(s string) time.Time
	BetweenTime                      func(t time.Time, start, end string) bool
}

func NewExpEval(input string, parseTime parsetime.ParseTime) (*ExpEval, error) {
	program, err := expr.Compile(input, expr.Env(&ExpEvalEnv{}), expr.AsBool(),
		expr.Operator("==", "TimeStringEqualTime", "TimeEqualTime"),
		expr.Operator("!=", "TimeStringNotEqualTime", "TimeNotEqualTime"),
		expr.Operator(">", "TimeStringGreaterThanTime", "TimeGreaterThanTime"),
		expr.Operator(">=", "TimeStringGreaterThanOrEqualTime", "TimeGreaterThanOrEqualTime"),
		expr.Operator("<", "TimeStringLessThanTime", "TimeLessThanTime"),
		expr.Operator("<=", "TimeStringLessThanOrEqualTime", "TimeLessThanOrEqualTime"),
	)
	if err != nil {
		return nil, err
	}

	return &ExpEval{
		program:   program,
		parseTime: parseTime,
	}, nil
}

func (ee *ExpEval) Run(request *httpv1.Request) (bool, error) {
	env := &ExpEvalEnv{
		Uri:                              request.Data.URLReference(),
		Method:                           request.Data.Method,
		Time:                             eventTime(request),
		ResponseTime:                     float64(request.DurationNano) / float64(time.Second),
		BodyBytes:                        responseBodySize(request),
		Status:                           statusCode(request),
		TimeStringEqualTime:              ee.timeStringEqual,
		TimeStringNotEqualTime:           ee.timeStringNotEqual,
		TimeStringGreaterThanTime:        ee.timeStringGreaterThan,
		TimeStringGreaterThanOrEqualTime: ee.timeStringGreaterThanOrEqual,
		TimeStringLessThanTime:           ee.timeStringLessThan,
		TimeStringLessThanOrEqualTime:    ee.timeStringLessThanOrEqual,
		TimeEqualTime:                    timeEqual,
		TimeNotEqualTime:                 timeNotEqual,
		TimeGreaterThanTime:              timeGreaterThan,
		TimeGreaterThanOrEqualTime:       timeGreaterThanOrEqual,
		TimeLessThanTime:                 timeLessThan,
		TimeLessThanOrEqualTime:          timeLessThanOrEqual,
		TimeAgo:                          timeAgo,
		BetweenTime:                      ee.betweenTime,
	}

	output, err := expr.Run(ee.program, env)
	if err != nil {
		return false, err
	}

	return output.(bool), nil
}

func eventTime(request *httpv1.Request) time.Time {
	if request.TimeUnixNano == nil {
		return time.Time{}
	}

	return time.Unix(0, int64(*request.TimeUnixNano))
}

func responseBodySize(request *httpv1.Request) float64 {
	if request.Data.ResponseBodySizeBytes == nil {
		return -1
	}

	return float64(*request.Data.ResponseBodySizeBytes)
}

func statusCode(request *httpv1.Request) int {
	if request.Data.StatusCode == nil {
		return 0
	}

	return *request.Data.StatusCode
}

func (ee *ExpEval) parse(value string) time.Time {
	parsed, err := ee.parseTime.Parse(value)
	if err != nil {
		panic(err)
	}

	return parsed
}

func (ee *ExpEval) timeStringEqual(left time.Time, right string) bool {
	return left.Equal(ee.parse(right))
}

func (ee *ExpEval) timeStringNotEqual(left time.Time, right string) bool {
	return !left.Equal(ee.parse(right))
}

func (ee *ExpEval) timeStringGreaterThan(left time.Time, right string) bool {
	return left.After(ee.parse(right))
}

func (ee *ExpEval) timeStringGreaterThanOrEqual(left time.Time, right string) bool {
	parsed := ee.parse(right)

	return left.After(parsed) || left.Equal(parsed)
}

func (ee *ExpEval) timeStringLessThan(left time.Time, right string) bool {
	return left.Before(ee.parse(right))
}

func (ee *ExpEval) timeStringLessThanOrEqual(left time.Time, right string) bool {
	parsed := ee.parse(right)

	return left.Before(parsed) || left.Equal(parsed)
}

func timeEqual(left, right time.Time) bool {
	return left.Equal(right)
}

func timeNotEqual(left, right time.Time) bool {
	return !left.Equal(right)
}

func timeGreaterThan(left, right time.Time) bool {
	return left.After(right)
}

func timeGreaterThanOrEqual(left, right time.Time) bool {
	return left.After(right) || left.Equal(right)
}

func timeLessThan(left, right time.Time) bool {
	return left.Before(right)
}

func timeLessThanOrEqual(left, right time.Time) bool {
	return left.Before(right) || left.Equal(right)
}

func timeAgo(value string) time.Time {
	duration, err := time.ParseDuration(value)
	if err != nil {
		panic(err)
	}

	return time.Now().Add(-duration)
}

func (ee *ExpEval) betweenTime(value time.Time, start, end string) bool {
	startTime := ee.parse(start)
	endTime := ee.parse(end)

	return !value.Before(startTime) && !value.After(endTime)
}
