package stats

import (
	"bytes"
	"testing"

	godiff "github.com/kylelemons/godebug/diff"
)

func TestDumpStats(t *testing.T) {
	got := new(bytes.Buffer)
	stats := NewHTTPStats(true, false, false)
	stats.Set("/foo/bar", "POST", 200, 0.057, 12, 0)

	err := stats.DumpStats(got)
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
  request_body_bytes:
    max: 12
    min: 12
    sum: 12
    usepercentile: false
    percentiles: []
  response_body_bytes:
    max: 0
    min: 0
    sum: 0
    usepercentile: false
    percentiles: []
  time: ""
`)

	if diff := godiff.Diff(got.String(), want.String()); diff != "" {
		t.Errorf("diff\n%s", diff)
	}
}
