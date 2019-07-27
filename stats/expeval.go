package stats

import (
	"time"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/k0kubun/pp"
	"github.com/tkuchiki/alp/parsers"
	"github.com/tkuchiki/parsetime"
)

type ExpEval struct {
	program   *vm.Program
	parseTime parsetime.ParseTime
}

type ExpEvalEnv struct {
	Val                    *parsers.ParsedHTTPStat
	GreaterThanTime        func(l time.Time, r string) bool
	GreaterThanOrEqualTime func(l time.Time, r string) bool
	LessThanTime           func(l time.Time, r string) bool
	LessThanOrEqualTime    func(l time.Time, r string) bool
	Datetime               func(s string) time.Time
	TimeAgo                func(s string) time.Time
	BetweenTime            func(t, start, end string) bool
	IsStatus1xx            func(status int) bool
	IsStatus2xx            func(status int) bool
	IsStatus3xx            func(status int) bool
	IsStatus4xx            func(status int) bool
	IsStatus5xx            func(status int) bool
	IsNotStatus1xx         func(status int) bool
	IsNotStatus2xx         func(status int) bool
	IsNotStatus3xx         func(status int) bool
	IsNotStatus4xx         func(status int) bool
	IsNotStatus5xx         func(status int) bool
}

var parseTime parsetime.ParseTime

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

	pp.Println(l, t)

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

func IsStatus1xx(status int) bool {
	return 100 <= status && status <= 199
}

func IsStatus2xx(status int) bool {
	return 200 <= status && status <= 299
}

func IsStatus3xx(status int) bool {
	return 300 <= status && status <= 399
}

func IsStatus4xx(status int) bool {
	return 400 <= status && status <= 499
}

func IsStatus5xx(status int) bool {
	return 500 <= status && status <= 599
}

func IsNotStatus1xx(status int) bool {
	return !(100 <= status && status <= 199)
}

func IsNotStatus2xx(status int) bool {
	return !(200 <= status && status <= 299)
}

func IsNotStatus3xx(status int) bool {
	return !(300 <= status && status <= 399)
}

func IsNotStatus4xx(status int) bool {
	return !(400 <= status && status <= 499)
}

func IsNotStatus5xx(status int) bool {
	return !(500 <= status && status <= 599)
}

func NewExpEval(input string, pt parsetime.ParseTime) (*ExpEval, error) {
	program, err := expr.Compile(input, expr.Env(&ExpEvalEnv{}), expr.AsBool(),
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
		Val:                    stat,
		GreaterThanTime:        GreaterThanTime,
		GreaterThanOrEqualTime: GreaterThanOrEqualTime,
		LessThanTime:           LessThanTime,
		LessThanOrEqualTime:    LessThanOrEqualTime,
		Datetime:               Datetime,
		TimeAgo:                TimeAgo,
		BetweenTime:            BetweenTime,
		IsStatus1xx:            IsStatus1xx,
		IsStatus2xx:            IsStatus2xx,
		IsStatus3xx:            IsStatus3xx,
		IsStatus4xx:            IsStatus4xx,
		IsStatus5xx:            IsStatus5xx,
		IsNotStatus1xx:         IsNotStatus1xx,
		IsNotStatus2xx:         IsNotStatus2xx,
		IsNotStatus3xx:         IsNotStatus3xx,
		IsNotStatus4xx:         IsNotStatus4xx,
		IsNotStatus5xx:         IsNotStatus5xx,
	}

	output, err := expr.Run(ee.program, env)
	if err != nil {
		return false, err
	}

	return output.(bool), nil
}

func (ee *ExpEval) Now() time.Time {
	return time.Now()
}
