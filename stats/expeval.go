package stats

import (
	"time"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/tkuchiki/alp/parsers"
	"github.com/tkuchiki/parsetime"
)

type ExpEval struct {
	program   *vm.Program
	parseTime parsetime.ParseTime
}

type ExpEvalEnv struct {
	Uri                              string
	Method                           string
	Time                             string
	ResponseTime                     float64
	BodyBytes                        float64
	Status                           int
	TimeStringEqualTime              func(l time.Time, r string) bool
	TimeStringNotEqualTime           func(l time.Time, r string) bool
	TimeStringGreaterThanTime        func(l time.Time, r string) bool
	TimeStringGreaterThanOrEqualTime func(l time.Time, r string) bool
	TimeStringLessThanTime           func(l time.Time, r string) bool
	TimeStringLessThanOrEqualTime    func(l time.Time, r string) bool
	StringTimeEqualTime              func(l string, r time.Time) bool
	StringTimeNotEqualTime           func(l string, r time.Time) bool
	StringTimeGreaterThanTime        func(l string, r time.Time) bool
	StringTimeGreaterThanOrEqualTime func(l string, r time.Time) bool
	StringTimeLessThanTime           func(l string, r time.Time) bool
	StringTimeLessThanOrEqualTime    func(l string, r time.Time) bool
	TimeAgo                          func(s string) time.Time
	BetweenTime                      func(t, start, end string) bool
}

var parseTime parsetime.ParseTime

// =
func TimeStringEqualTime(l time.Time, r string) bool {
	t, err := parseTime.Parse(r)
	if err != nil {
		panic(err)
	}

	return l.Equal(t)
}

// !=
func TimeStringNotEqualTime(l time.Time, r string) bool {
	t, err := parseTime.Parse(r)
	if err != nil {
		panic(err)
	}

	return !l.Equal(t)
}

// >
func TimeStringGreaterThanTime(l time.Time, r string) bool {
	t, err := parseTime.Parse(r)
	if err != nil {
		panic(err)
	}

	return l.After(t)
}

// >=
func TimeStringGreaterThanOrEqualTime(l time.Time, r string) bool {
	t, err := parseTime.Parse(r)
	if err != nil {
		panic(err)
	}

	return l.After(t) || l.Equal(t)
}

// <
func TimeStringLessThanTime(l time.Time, r string) bool {
	t, err := parseTime.Parse(r)
	if err != nil {
		panic(err)
	}

	return l.Before(t)
}

// <=
func TimeStringLessThanOrEqualTime(l time.Time, r string) bool {
	t, err := parseTime.Parse(r)
	if err != nil {
		panic(err)
	}

	return l.Before(t) || l.Equal(t)
}

// =
func StringTimeEqualTime(l string, r time.Time) bool {
	t, err := parseTime.Parse(l)
	if err != nil {
		panic(err)
	}

	return t.Equal(r)
}

// !=
func StringTimeNotEqualTime(l string, r time.Time) bool {
	t, err := parseTime.Parse(l)
	if err != nil {
		panic(err)
	}

	return !t.Equal(r)
}

// >
func StringTimeGreaterThanTime(l string, r time.Time) bool {
	t, err := parseTime.Parse(l)
	if err != nil {
		panic(err)
	}

	return t.After(r)
}

// >=
func StringTimeGreaterThanOrEqualTime(l string, r time.Time) bool {
	t, err := parseTime.Parse(l)
	if err != nil {
		panic(err)
	}

	return t.After(r) || t.Equal(r)
}

// <
func StringTimeLessThanTime(l string, r time.Time) bool {
	t, err := parseTime.Parse(l)
	if err != nil {
		panic(err)
	}

	return t.Before(r)
}

// <=
func StringTimeLessThanOrEqualTime(l string, r time.Time) bool {
	t, err := parseTime.Parse(l)
	if err != nil {
		panic(err)
	}

	return t.Before(r) || t.Equal(r)
}

func TimeAgo(s string) time.Time {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}

	return time.Now().Add(-1 * d)
}

func BetweenTime(t, start, end string) bool {
	val, err := parseTime.Parse(t)
	if err != nil {
		panic(err)
	}

	st, err := parseTime.Parse(start)
	if err != nil {
		panic(err)
	}

	et, err := parseTime.Parse(end)
	if err != nil {
		panic(err)
	}

	return st.UnixNano() <= val.UnixNano() && val.UnixNano() <= et.UnixNano()
}

func NewExpEval(input string, pt parsetime.ParseTime) (*ExpEval, error) {
	program, err := expr.Compile(input, expr.Env(&ExpEvalEnv{}), expr.AsBool(),
		expr.Operator("=", "TimeStringEqualTime"),
		expr.Operator("!=", "TimeStringNotEqualTime"),
		expr.Operator(">", "TimeStringGreaterThanTime"),
		expr.Operator(">=", "TimeStringGreaterThanOrEqualTime"),
		expr.Operator("<", "TimeStringLessThanTime"),
		expr.Operator("<=", "TimeStringLessThanOrEqualTime"),
		expr.Operator("=", "StringTimeEqualTime"),
		expr.Operator("!=", "StringTimeNotEqualTime"),
		expr.Operator(">", "StringTimeGreaterThanTime"),
		expr.Operator(">=", "StringTimeGreaterThanOrEqualTime"),
		expr.Operator("<", "StringTimeLessThanTime"),
		expr.Operator("<=", "StringTimeLessThanOrEqualTime"),
	)
	if err != nil {
		return nil, err
	}

	parseTime = pt

	return &ExpEval{
		program: program,
	}, nil
}

func (ee *ExpEval) Run(stat *parsers.ParsedHTTPStat) (bool, error) {
	env := &ExpEvalEnv{
		Uri:                              stat.Uri,
		Method:                           stat.Method,
		Time:                             stat.Time,
		ResponseTime:                     stat.ResponseTime,
		BodyBytes:                        stat.BodyBytes,
		Status:                           stat.Status,
		TimeStringEqualTime:              TimeStringEqualTime,
		TimeStringNotEqualTime:           TimeStringNotEqualTime,
		TimeStringGreaterThanTime:        TimeStringGreaterThanTime,
		TimeStringGreaterThanOrEqualTime: TimeStringGreaterThanOrEqualTime,
		TimeStringLessThanTime:           TimeStringLessThanTime,
		TimeStringLessThanOrEqualTime:    TimeStringLessThanOrEqualTime,
		StringTimeEqualTime:              StringTimeEqualTime,
		StringTimeNotEqualTime:           StringTimeNotEqualTime,
		StringTimeGreaterThanTime:        StringTimeGreaterThanTime,
		StringTimeGreaterThanOrEqualTime: StringTimeGreaterThanOrEqualTime,
		StringTimeLessThanTime:           StringTimeLessThanTime,
		StringTimeLessThanOrEqualTime:    StringTimeLessThanOrEqualTime,
		TimeAgo:                          TimeAgo,
		BetweenTime:                      BetweenTime,
	}

	output, err := expr.Run(ee.program, env)
	if err != nil {
		return false, err
	}

	return output.(bool), nil
}
