package cmd

import (
	"testing"

	"github.com/tkuchiki/alp/log_reader"

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

	tempFile, err := testutil.CreateTempDirAndFile(t.TempDir(), "test_ltsv_cmd_temp_file", ltsvLog)
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

	command := NewCommand("test")
	command.setArgs(args)

	err = command.Execute()
	if err != nil {
		t.Fatal(err)
	}
}

func TestLTSVDiffCmd(t *testing.T) {
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

	tempDir := t.TempDir()

	tempFromFile, err := testutil.CreateTempDirAndFile(tempDir, "test_ltsv_diff_cmd_temp_from_file", ltsvLog)
	if err != nil {
		t.Fatal(err)
	}

	tempToFile, err := testutil.CreateTempDirAndFile(tempDir, "test_ltsv_diff_cmd_temp_to_file", ltsvLog)
	if err != nil {
		t.Fatal(err)
	}

	tempDump, err := testutil.CreateTempDirAndFile(tempDir, "test_ltsv_diff_cmd_temp_dump", "")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("alp ltsv diff <from> <to>", func(t *testing.T) {
		args := []string{"ltsv", "diff",
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

	t.Run("alp ltsv diff --load <dumpfile> <to>", func(t *testing.T) {
		args := []string{"ltsv", "diff",
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

func TestLTSVTopNCmd(t *testing.T) {
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

	tempFile, err := testutil.CreateTempDirAndFile(t.TempDir(), "test_ltsv_topN_cmd_temp_file", ltsvLog)
	if err != nil {
		t.Fatal(err)
	}

	args := []string{"ltsv", "topN", "10",
		"--file", tempFile,
		"--uri-label", keys.Uri,
		"--method-label", keys.Method,
		"--time-label", keys.Time,
		"--apptime-label", keys.ResponseTime,
		"--reqtime-label", keys.RequestTime,
		"--size-label", keys.BodyBytes,
		"--status-label", keys.Status,
		"--reverse",
		"--sort", log_reader.SortBodyBytes,
	}

	command := NewCommand("test")
	command.setArgs(args)

	err = command.Execute()
	if err != nil {
		t.Fatal(err)
	}
}

func TestLTSVCountCmd(t *testing.T) {
	keys := testutil.NewLTSVLogKeys()

	ltsvLog := testutil.LTSVLog(keys)

	tempFile, err := testutil.CreateTempDirAndFile(t.TempDir(), "test_ltsv_count_cmd_temp_file", ltsvLog)
	if err != nil {
		t.Fatal(err)
	}

	args := []string{"ltsv", "count",
		"--file", tempFile,
		"--reverse",
		"--keys", "ua",
	}

	command := NewCommand("test")
	command.setArgs(args)

	err = command.Execute()
	if err != nil {
		t.Fatal(err)
	}
}
