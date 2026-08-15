package parsers

import (
	"encoding/json"
	"math"
	"net/url"
	"strings"
	"testing"
	"time"

	corev1 "github.com/tkuchiki/logschema/core/v1"
	httpv1 "github.com/tkuchiki/logschema/http/v1"
)

func TestJSONParserReturnsLogSchemaRecord(t *testing.T) {
	input := `{"time":"2015-09-06T05:58:05+09:00","method":"POST","uri":"/search?b=2&a=1","status":200,"body_bytes":12,"response_time":0.057,"ua":"UA1","context":{"b":2,"a":1},"nullable":null}`
	parser, err := NewJSONParser(strings.NewReader(input), NewJSONKeys("", "", "", "", "", "", ""), true, false, "Local")
	if err != nil {
		t.Fatal(err)
	}

	record, err := parser.Parse()
	if err != nil {
		t.Fatal(err)
	}
	if record.SchemaVersion != httpv1.SchemaVersion || record.Kind != httpv1.Kind {
		t.Fatalf("schema metadata = %q, %q", record.SchemaVersion, record.Kind)
	}
	if err := record.Validate(); err != nil {
		t.Fatalf("parsed record is invalid: %v", err)
	}

	if got, want := record.Data.URI(), "/search?a=1&b=2"; got != want {
		t.Fatalf("URI() = %q, want %q", got, want)
	}
	if got, want := record.Data.Method, "POST"; got != want {
		t.Fatalf("method = %q, want %q", got, want)
	}
	wantTime, err := time.Parse(time.RFC3339, "2015-09-06T05:58:05+09:00")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := int64(*record.TimeUnixNano), wantTime.UnixNano(); got != want {
		t.Fatalf("event time = %d, want %d", got, want)
	}
	if got, want := float64(record.DurationNano)/float64(time.Second), 0.057; got != want {
		t.Fatalf("response time = %v, want %v", got, want)
	}
	if got, want := uint64(*record.Data.ResponseBodySizeBytes), uint64(12); got != want {
		t.Fatalf("response body size = %v, want %v", got, want)
	}
	if got, want := *record.Data.StatusCode, 200; got != want {
		t.Fatalf("status = %d, want %d", got, want)
	}
	if got, want := OriginalTime(record), "2015-09-06T05:58:05+09:00"; got != want {
		t.Fatalf("original time = %q, want %q", got, want)
	}
	if got, ok := record.Data.Attributes["ua"]; !ok || got != "UA1" {
		t.Fatalf("ua attribute = %q, %v", got, ok)
	}
	if got, ok := record.Data.Attributes["body_bytes"].(json.Number); !ok || got != "12" {
		t.Fatalf("body_bytes attribute = %v (%T)", got, record.Data.Attributes["body_bytes"])
	}
	if got, want := record.Data.Attributes["context"], `{"a":1,"b":2}`; got != want {
		t.Fatalf("context attribute = %q, want %q", got, want)
	}
	if _, ok := record.Data.Attributes["nullable"]; ok {
		t.Fatal("null attribute was not omitted")
	}
}

func TestOriginalTimePreservesUnusualInput(t *testing.T) {
	const input = `{"time":"not-a-time","method":"GET","uri":"/","status":200,"body_bytes":1,"response_time":0.001}`
	parser, err := NewJSONParser(strings.NewReader(input), NewJSONKeys("", "", "", "", "", "", ""), false, false, "Local")
	if err != nil {
		t.Fatal(err)
	}

	record, err := parser.Parse()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := OriginalTime(record), "not-a-time"; got != want {
		t.Fatalf("original time = %q, want %q", got, want)
	}
}

func TestHTTPRequestNetworkPathReference(t *testing.T) {
	u, err := url.Parse("//example.com/assets/app.js")
	if err != nil {
		t.Fatal(err)
	}

	record, err := newHTTPRequest(u, "GET", nil, 0.01, 512, 200, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := record.Data.URLReference(), "//example.com/assets/app.js"; got != want {
		t.Fatalf("URLReference() = %q, want %q", got, want)
	}
	if got, want := *record.Data.URLAuthority, "example.com"; got != want {
		t.Fatalf("URLAuthority = %q, want %q", got, want)
	}
}

func TestHTTPRequestAbsoluteURL(t *testing.T) {
	u, err := url.Parse("https://example.com:8443/users/42?view=compact")
	if err != nil {
		t.Fatal(err)
	}

	record, err := newHTTPRequest(u, "GET", nil, 0.125, 1536, 200, nil)
	if err != nil {
		t.Fatal(err)
	}
	got, ok := record.Data.URLFull()
	if want := "https://example.com:8443/users/42?view=compact"; !ok || got != want {
		t.Fatalf("URLFull() = %q, %v, want %q", got, ok, want)
	}
}

func TestNormalizeURLRedactsAndSortsQuery(t *testing.T) {
	u, err := url.Parse("/login?token=secret&redirect=%2Fhome")
	if err != nil {
		t.Fatal(err)
	}

	normalized := normalizeURL(u, true, true)
	if got, want := normalized.String(), "/login?redirect=xxx&token=xxx"; got != want {
		t.Fatalf("normalized URL = %q, want %q", got, want)
	}
}

func TestRequestDataFromOpaqueTarget(t *testing.T) {
	u, err := url.Parse("example.com:443")
	if err != nil {
		t.Fatal(err)
	}

	data := requestDataFromURL(u)
	if got, want := data.URI(), "example.com:443"; got != want {
		t.Fatalf("URI() = %q, want %q", got, want)
	}
	if data.URLScheme != nil {
		t.Fatalf("opaque request target has URL scheme %q", *data.URLScheme)
	}
}

func TestNumericConversions(t *testing.T) {
	duration, err := secondsToDurationNano(0.057)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := duration, corev1.DecimalUint64(57_000_000); got != want {
		t.Fatalf("duration = %d, want %d", got, want)
	}

	for _, value := range []float64{-1, math.Inf(1), math.NaN()} {
		if _, err := secondsToDurationNano(value); err == nil {
			t.Fatalf("invalid duration %v was accepted", value)
		}
	}

	if _, err := floatToUint64(1.5, "size"); err == nil {
		t.Fatal("fractional byte size was accepted")
	}
}

func TestUnknownResponseBodySize(t *testing.T) {
	u, err := url.Parse("/stream")
	if err != nil {
		t.Fatal(err)
	}

	record, err := newHTTPRequest(u, "GET", nil, 0.1, -1, 200, nil)
	if err != nil {
		t.Fatal(err)
	}
	if record.Data.ResponseBodySizeBytes != nil {
		t.Fatalf("unknown response body size = %d, want nil", *record.Data.ResponseBodySizeBytes)
	}
}

func BenchmarkJSONParser(b *testing.B) {
	const input = `{"time":"2015-09-06T05:58:05+09:00","method":"POST","uri":"/search?b=2&a=1","status":200,"body_bytes":12,"response_time":0.057,"ua":"UA1"}`
	keys := NewJSONKeys("", "", "", "", "", "", "")

	b.ReportAllocs()
	for b.Loop() {
		parser, err := NewJSONParser(strings.NewReader(input), keys, true, false, "Local")
		if err != nil {
			b.Fatal(err)
		}
		if _, err := parser.Parse(); err != nil {
			b.Fatal(err)
		}
	}
}
