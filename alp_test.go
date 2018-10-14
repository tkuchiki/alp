package alp

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tkuchiki/gohttpstats"
)

func TestRun(t *testing.T) {
	outw, errw := new(bytes.Buffer), new(bytes.Buffer)
	p := NewProfiler(outw, errw)

	data := []string{
		"time:2018-10-14T05:58:05+09:00",
		"method:POST",
		"uri:/foo/bar?token=xxx&uuid=1234",
		"status:200",
		"size:12",
		"request_time:0.247",
		"apptime:0.057",
	}

	ltsvData := strings.Join(data, "\t")

	tmpfile, err := ioutil.TempFile(os.TempDir(), "alp-mock-stdin")
	assert.Nil(t, err)
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()
	if _, err := tmpfile.Write([]byte(ltsvData)); err != nil {
		assert.Nil(t, err)
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		assert.Nil(t, err)
	}

	flags := httpstats.NewEmptyFlags()

	p.SetFlags(flags)
	p.SetInReader(tmpfile)
	err = p.Run()
	assert.Nil(t, err)
}
