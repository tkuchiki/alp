package stats

import (
	"strings"
	"testing"
	"time"

	"github.com/tkuchiki/alp/parsers"
	corev1 "github.com/tkuchiki/logschema/core/v1"
	httpv1 "github.com/tkuchiki/logschema/http/v1"
	"github.com/tkuchiki/parsetime"
)

func TestExpEvalUsesLogSchemaRecord(t *testing.T) {
	const input = `{"time":"2015-09-06T05:58:05+09:00","method":"POST","uri":"/foo","status":200,"body_bytes":12,"response_time":0.057}`
	parser, err := parsers.NewJSONParser(
		strings.NewReader(input),
		parsers.NewJSONKeys("", "", "", "", "", "", ""),
		false,
		false,
		"Local",
	)
	if err != nil {
		t.Fatal(err)
	}
	record, err := parser.Parse()
	if err != nil {
		t.Fatal(err)
	}

	parseTime, err := parsetime.NewParseTime("Local")
	if err != nil {
		t.Fatal(err)
	}
	evaluator, err := NewExpEval(
		`Uri == "/foo" && Method == "POST" && Time == "2015-09-06T05:58:05+09:00" && ResponseTime == 0.057 && BodyBytes == 12 && Status == 200`,
		parseTime,
	)
	if err != nil {
		t.Fatal(err)
	}

	matched, err := evaluator.Run(record)
	if err != nil {
		t.Fatal(err)
	}
	if !matched {
		t.Fatal("legacy filter values did not match the LogSchema-backed record")
	}
}

func TestExpEvalComparesCanonicalTimeToTimeAgo(t *testing.T) {
	parseTime, err := parsetime.NewParseTime("Local")
	if err != nil {
		t.Fatal(err)
	}
	evaluator, err := NewExpEval(`Time < TimeAgo("1h")`, parseTime)
	if err != nil {
		t.Fatal(err)
	}

	eventTime := corev1.DecimalInt64(time.Now().Add(-2 * time.Hour).UnixNano())
	record, err := httpv1.NewRequest(0, corev1.Source{Kind: corev1.SourceOther}, httpv1.RequestData{
		Method:  "GET",
		URLPath: "/",
	})
	if err != nil {
		t.Fatal(err)
	}
	record.TimeUnixNano = &eventTime

	matched, err := evaluator.Run(&record)
	if err != nil {
		t.Fatal(err)
	}
	if !matched {
		t.Fatal("canonical event time did not match TimeAgo filter")
	}
}

func TestExpEvalTimeOperators(t *testing.T) {
	parseTime, err := parsetime.NewParseTime("Local")
	if err != nil {
		t.Fatal(err)
	}

	eventTime := corev1.DecimalInt64(time.Date(2015, 9, 6, 5, 58, 5, 0, time.FixedZone("JST", 9*60*60)).UnixNano())
	record, err := httpv1.NewRequest(0, corev1.Source{Kind: corev1.SourceOther}, httpv1.RequestData{
		Method:  "GET",
		URLPath: "/",
	})
	if err != nil {
		t.Fatal(err)
	}
	record.TimeUnixNano = &eventTime

	tests := []string{
		`Time == "2015-09-06T05:58:05+09:00"`,
		`Time != "2015-09-06T05:58:06+09:00"`,
		`Time > "2015-09-06T05:58:04+09:00"`,
		`Time >= "2015-09-06T05:58:05+09:00"`,
		`Time < "2015-09-06T05:58:06+09:00"`,
		`Time <= "2015-09-06T05:58:05+09:00"`,
		`"2015-09-06T05:58:05+09:00" == Time`,
		`"2015-09-06T05:58:06+09:00" != Time`,
		`"2015-09-06T05:58:06+09:00" > Time`,
		`"2015-09-06T05:58:05+09:00" >= Time`,
		`"2015-09-06T05:58:04+09:00" < Time`,
		`"2015-09-06T05:58:05+09:00" <= Time`,
		`BetweenTime(Time, "2015-09-06T05:58:04+09:00", "2015-09-06T05:58:06+09:00")`,
	}

	for _, expression := range tests {
		t.Run(expression, func(t *testing.T) {
			evaluator, err := NewExpEval(expression, parseTime)
			if err != nil {
				t.Fatal(err)
			}

			matched, err := evaluator.Run(&record)
			if err != nil {
				t.Fatal(err)
			}
			if !matched {
				t.Fatalf("time expression did not match: %s", expression)
			}
		})
	}
}
