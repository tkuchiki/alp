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
	Uri                    string
	Method                 string
	Time                   string
	ResponseTime           float64
	BodyBytes              float64
	Status                 int
	EqualTime              func(l time.Time, r string) bool
	NotEqualTime           func(l time.Time, r string) bool
	GreaterThanTime        func(l time.Time, r string) bool
	GreaterThanOrEqualTime func(l time.Time, r string) bool
	LessThanTime           func(l time.Time, r string) bool
	LessThanOrEqualTime    func(l time.Time, r string) bool
	TimeAgo                func(s string) time.Time
	BetweenTime            func(t, start, end string) bool
}

var parseTime parsetime.ParseTime

// =
func EqualTime(l time.Time, r string) bool {
	t, err := parseTime.Parse(r)
	if err != nil {
		panic(err)
	}

	return l.Equal(t)
}

// !=
func NotEqualTime(l time.Time, r string) bool {
	t, err := parseTime.Parse(r)
	if err != nil {
		panic(err)
	}

	return !l.Equal(t)
}

// >
func GreaterThanTime(l time.Time, r string) bool {
	t, err := parseTime.Parse(r)
	if err != nil {
		panic(err)
	}

	return l.After(t)
}

// >=
func GreaterThanOrEqualTime(l time.Time, r string) bool {
	t, err := parseTime.Parse(r)
	if err != nil {
		panic(err)
	}

	return l.After(t) || l.Equal(t)
}

// <
func LessThanTime(l time.Time, r string) bool {
	t, err := parseTime.Parse(r)
	if err != nil {
		panic(err)
	}

	return l.Before(t)
}

// <=
func LessThanOrEqualTime(l time.Time, r string) bool {
	t, err := parseTime.Parse(r)
	if err != nil {
		panic(err)
	}

	return l.Before(t) || l.Equal(t)
}

func Datetime(s string) time.Time {
	datetime, err := parseTime.Parse(s)
	if err != nil {
		panic(err)
	}

	return datetime
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
		expr.Operator("=", "EqualTime"),
		expr.Operator("!=", "NotEqualTime"),
		expr.Operator(">", "GreaterThanTime"),
		expr.Operator(">=", "GreaterThanOrEqualTime"),
		expr.Operator("<", "LessThanTime"),
		expr.Operator("<=", "LessThanOrEqualTime"),
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
		Uri:                    stat.Uri,
		Method:                 stat.Method,
		Time:                   stat.Time,
		ResponseTime:           stat.ResponseTime,
		BodyBytes:              stat.BodyBytes,
		Status:                 stat.Status,
		EqualTime:              EqualTime,
		NotEqualTime:           NotEqualTime,
		GreaterThanTime:        GreaterThanTime,
		GreaterThanOrEqualTime: GreaterThanOrEqualTime,
		LessThanTime:           LessThanTime,
		LessThanOrEqualTime:    LessThanOrEqualTime,
		TimeAgo:                TimeAgo,
		BetweenTime:            BetweenTime,
	}

	output, err := expr.Run(ee.program, env)
	if err != nil {
		return false, err
	}

	return output.(bool), nil
}
