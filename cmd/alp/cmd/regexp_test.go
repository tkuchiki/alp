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

	command := NewCommand("test")
	command.setArgs(args)

	err = command.Execute()
	if err != nil {
		t.Fatal(err)
	}
}

func TestRegexpDiffCmd(t *testing.T) {
	regexpLog := testutil.RegexpLog()

	tempDir := t.TempDir()

	tempFromFile, err := testutil.CreateTempDirAndFile(tempDir, "test_regexp_diff_cmd_temp_from_file", regexpLog)
	if err != nil {
		t.Fatal(err)
	}

	tempToFile, err := testutil.CreateTempDirAndFile(tempDir, "test_regexp_diff_cmd_temp_to_file", regexpLog)
	if err != nil {
		t.Fatal(err)
	}

	tempDump, err := testutil.CreateTempDirAndFile(tempDir, "test_regexp_diff_cmd_temp_dump", "")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("alp regexp diff <from> <to>", func(t *testing.T) {
		args := []string{"regexp", "diff",
			tempFromFile,
			tempToFile,
			"--dump", tempDump,
		}

		command := NewCommand("test")
		command.setArgs(args)

		err = command.Execute()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("alp regexp diff --load <dumpfile> <to>", func(t *testing.T) {
		args := []string{"regexp", "diff",
			"--load", tempDump,
			tempToFile,
		}

		command := NewCommand("test")
		command.setArgs(args)

		err = command.Execute()
		if err != nil {
			t.Fatal(err)
		}
	})
}
