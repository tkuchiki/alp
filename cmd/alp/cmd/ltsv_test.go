package cmd

import (
	"testing"

	"github.com/tkuchiki/alp/internal/testutil"
)

func TestLTSVCmd(t *testing.T) {
	keys := testutil.LogKeys{
		Uri:          "u",
		Method:       "m",
		Time:         "t",
		ResponseTime: "r",
		RequestTime:  "r2",
		BodyBytes:    "b",
		Status:       "s",
	}

	ltsvLog := testutil.LTSVLog(keys)

	tempFile, err := testutil.CreateTempDirAndFile(t.TempDir(), ltsvLog)
	if err != nil {
		t.Fatal(err)
	}

	args := []string{"ltsv",
		"--file", tempFile,
		"--uri-label", keys.Uri,
		"--method-label", keys.Method,
		"--time-label", keys.Time,
		"--apptime-label", keys.ResponseTime,
		"--reqtime-label", keys.RequestTime,
		"--size-label", keys.BodyBytes,
		"--status-label", keys.Status,
	}

	rootCmd := NewRootCmd("test")
	rootCmd.SetArgs(args)

	err = rootCmd.Execute()
	if err != nil {
		t.Fatal(err)
	}
}
