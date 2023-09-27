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

	tempFile, err := testutil.CreateTempDirAndFile(t.TempDir(), "test_json_cmd_temp_file", jsonLog)
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

	command := NewCommand("test")
	command.setArgs(args)

	err = command.Execute()
	if err != nil {
		t.Fatal(err)
	}
}

func TestJSONDiffCmd(t *testing.T) {
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

	tempDir := t.TempDir()

	tempFromFile, err := testutil.CreateTempDirAndFile(tempDir, "test_json_diff_cmd_temp_from_file", jsonLog)
	if err != nil {
		t.Fatal(err)
	}

	tempToFile, err := testutil.CreateTempDirAndFile(tempDir, "test_json_diff_cmd_temp_to_file", jsonLog)
	if err != nil {
		t.Fatal(err)
	}

	tempDump, err := testutil.CreateTempDirAndFile(tempDir, "test_json_diff_cmd_temp_dump", "")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("alp json diff <from> <to>", func(t *testing.T) {
		args := []string{"json", "diff",
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

	t.Run("alp json diff --load <dumpfile> <to>", func(t *testing.T) {
		args := []string{"json", "diff",
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
