package cmd

import (
	"testing"

	"github.com/tkuchiki/alp/internal/testutil"
)

func TestRegexpCmd(t *testing.T) {
	keys := testutil.LogKeys{
		Uri:          "u",
		Method:       "m",
		Time:         "t",
		ResponseTime: "r",
		RequestTime:  "r2",
		BodyBytes:    "b",
		Status:       "s",
	}

	regexpLog := testutil.RegexpLog()

	tempFile, err := testutil.CreateTempDirAndFile(t.TempDir(), "test_regexp_cmd_temp_file", regexpLog)
	if err != nil {
		t.Fatal(err)
	}

	args := []string{"regexp",
		"--file", tempFile,
		"--pattern", testutil.RegexpPattern(keys),
		"--uri-subexp", keys.Uri,
		"--method-subexp", keys.Method,
		"--time-subexp", keys.Time,
		"--restime-subexp", keys.ResponseTime,
		"--reqtime-subexp", keys.RequestTime,
		"--body-bytes-subexp", keys.BodyBytes,
		"--status-subexp", keys.Status,
	}

	rootCmd := NewRootCmd("test")
	rootCmd.SetArgs(args)

	err = rootCmd.Execute()
	if err != nil {
		t.Fatal(err)
	}
}
