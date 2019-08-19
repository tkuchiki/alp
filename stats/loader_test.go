package stats

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadStats(t *testing.T) {
	data := bytes.NewBufferString(`- uri: /foo/bar
  cnt: 1
  status1xx: 0
  status2xx: 1
  status3xx: 0
  status4xx: 0
  status5xx: 0
  method: POST
  responsetime:
    max: 0.057
    min: 0.057
    sum: 0.057
    percentiles:
    - 0.057
  requestbodysize:
    max: 12
    min: 12
    sum: 12
  responsebodysize:
    max: 0
    min: 0
    sum: 0`)

	stats := NewHTTPStats(true, false, false)

	err := stats.LoadStats(data)
	assert.Nil(t, err)

	s := stats.Stats()

	assert.Equal(t, 1, s[0].Cnt)
	assert.Equal(t, "POST", s[0].Method)
	assert.Equal(t, "/foo/bar", s[0].Uri)
	assert.Equal(t, float64(12), s[0].RequestBodyBytes.Max)
	assert.Equal(t, 0.057, s[0].ResponseTime.Max)
	assert.Equal(t, 0, s[0].Status1xx)
	assert.Equal(t, 1, s[0].Status2xx)
	assert.Equal(t, 0, s[0].Status3xx)
	assert.Equal(t, 0, s[0].Status4xx)
	assert.Equal(t, 0, s[0].Status5xx)
}
