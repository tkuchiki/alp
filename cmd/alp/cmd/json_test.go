package cmd

import (
	"testing"

	"github.com/tkuchiki/alp/internal/testutil"
)

func TestJSONCmd(t *testing.T) {
	keys := testutil.LogKeys{
		Uri:          "u",
		Method:       "m",
		Time:         "t",
		ResponseTime: "r",
		RequestTime:  "r2",
		BodyBytes:    "b",
		Status:       "s",
	}

	jsonLog := testutil.JsonLog(keys)

	tempFile, err := testutil.CreateTempDirAndFile(t.TempDir(), jsonLog)
	if err != nil {
		t.Fatal(err)
	}

	args := []string{"json",
		"--file", tempFile,
		"--uri-key", keys.Uri,
		"--method-key", keys.Method,
		"--time-key", keys.Time,
		"--restime-key", keys.ResponseTime,
		"--reqtime-key", keys.RequestTime,
		"--body-bytes-key", keys.BodyBytes,
		"--status-key", keys.Status,
	}

	rootCmd := NewRootCmd("test")
	rootCmd.SetArgs(args)

	err = rootCmd.Execute()
	if err != nil {
		t.Fatal(err)
	}
}
