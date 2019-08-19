package stats

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDumpStats(t *testing.T) {
	outw := new(bytes.Buffer)
	po := NewPrintOptions()
	stats := NewHTTPStats(true, false, false, po)
	stats.Set("/foo/bar", "POST", 200, 0.057, 12, 0)

	err := stats.DumpStats(outw)
	assert.Nil(t, err)

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
    sum: 0
`)

	assert.Equal(t, data, outw)
}
