package stats

import (
	"bytes"
	"testing"

	godiff "github.com/kylelemons/godebug/diff"
	corev1 "github.com/tkuchiki/logschema/core/v1"
	httpv1 "github.com/tkuchiki/logschema/http/v1"
)

func TestDumpStats(t *testing.T) {
	got := new(bytes.Buffer)
	stats := NewHTTPStats(true, false, false)
	status := 200
	bodySize := corev1.DecimalUint64(12)
	record, err := httpv1.NewRequest(57_000_000, corev1.Source{Kind: corev1.SourceOther}, httpv1.RequestData{
		Method:                "POST",
		URLPath:               "/foo/bar",
		StatusCode:            &status,
		ResponseBodySizeBytes: &bodySize,
	})
	if err != nil {
		t.Fatal(err)
	}
	stats.Observe(&record)

	err = stats.DumpStats(got)
	if err != nil {
		t.Fatal(err)
	}

	want := bytes.NewBufferString(`- uri: /foo/bar
  count: 1
  status1xx: 0
  status2xx: 1
  status3xx: 0
  status4xx: 0
  status5xx: 0
  method: POST
  response_time:
    max: 0.057
    min: 0.057
    sum: 0.057
    usepercentile: true
    percentiles:
    - 0.057
    sample_count: 1
  request_body_bytes:
    max: 0
    min: 0
    sum: 0
    usepercentile: false
    percentiles: []
    sample_count: 0
  response_body_bytes:
    max: 12
    min: 12
    sum: 12
    usepercentile: false
    percentiles: []
    sample_count: 1
  time: ""
`)

	if diff := godiff.Diff(got.String(), want.String()); diff != "" {
		t.Errorf("diff\n%s", diff)
	}
}
