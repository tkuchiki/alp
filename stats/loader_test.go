package stats

import (
	"bytes"
	"testing"
)

func TestLoadStats(t *testing.T) {
	data := bytes.NewBufferString(`- uri: /foo/bar
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

	stats := NewHTTPStats(true, false, false)

	err := stats.LoadStats(data)
	if err != nil {
		t.Fatal(err)
	}

	s := stats.Stats()

	count := 1
	if count != s[0].Cnt {
		t.Errorf(`count want: %d, got: %d`, count, s[0].Cnt)
	}

	method := "POST"
	if method != s[0].Method {
		t.Errorf(`method want: %s, got: %s`, method, s[0].Method)
	}

	uri := "/foo/bar"
	if uri != s[0].Uri {
		t.Errorf(`uri want: %s, got: %s`, uri, s[0].Uri)
	}

	bodyMax := float64(12)
	if bodyMax != s[0].RequestBodyBytes.Max {
		t.Errorf(`request body bytes max want: %f, got: %f`, bodyMax, s[0].RequestBodyBytes.Max)
	}

	restimeMax := float64(0.057)
	if restimeMax != s[0].ResponseTime.Max {
		t.Errorf(`response time max want: %f, got: %f`, restimeMax, s[0].ResponseTime.Max)
	}

	status1xx := 0
	if status1xx != s[0].Status1xx {
		t.Errorf(`status1xx want: %d, got: %d`, status1xx, s[0].Status1xx)
	}

	status2xx := 1
	if status2xx != s[0].Status2xx {
		t.Errorf(`status2xx want: %d, got: %d`, status2xx, s[0].Status2xx)
	}

	status3xx := 0
	if status3xx != s[0].Status3xx {
		t.Errorf(`status3xx want: %d, got: %d`, status3xx, s[0].Status3xx)
	}

	status4xx := 0
	if status4xx != s[0].Status4xx {
		t.Errorf(`status4xx want: %d, got: %d`, status4xx, s[0].Status4xx)
	}

	status5xx := 0
	if status5xx != s[0].Status5xx {
		t.Errorf(`status5xx want: %d, got: %d`, status5xx, s[0].Status5xx)
	}
}
