package cmd

import (
	"strings"
	"testing"

	"github.com/tkuchiki/alp/internal/testutil"
)

func TestCommonFlags(t *testing.T) {
	tempDir := t.TempDir()

	tempLog, err := testutil.CreateTempDirAndFile(tempDir, "test_common_flags_temp_log", testutil.JsonLog(testutil.NewJsonLogKeys()))
	if err != nil {
		t.Fatal(err)
	}

	tempConfig, err := testutil.CreateTempDirAndFile(tempDir, "test_common_flags_temp_config", testutil.ConfigFile())
	if err != nil {
		t.Fatal(err)
	}

	tempPos, err := testutil.CreateTempDirAndFile(tempDir, "test_common_flags_temp_pos", "")
	if err != nil {
		t.Fatal(err)
	}

	tempDump, err := testutil.CreateTempDirAndFile(tempDir, "test_common_flags_temp_dump", "")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		args []string
	}{
		{
			args: []string{"json",
				"--file", tempLog,
				"--noheaders",
				"--format", "tsv",
				"--config", tempConfig,
			},
		},
		{
			args: []string{"json",
				"--file", tempLog,
				"--decode-uri",
				"--filters", "Method == 'POST'",
				"--format", "markdown",
				"--limit", "5",
				"--location", "America/Adak",
				"--matching-groups", "/foo/bar/.+",
				"--output", "count,uri",
				"--page", "10",
				"--percentiles", "99",
				"--qs-ignore-values",
				"--query-string",
				"--reverse",
				"--show-footers",
				"--sort", "uri",
			},
		},
		{
			args: []string{"json",
				"--file", tempLog,
				"-f", "Method == 'POST'",
				"-m", "/foo/bar/.+",
				"-o", "count,uri",
				"-q",
				"-r",
			},
		},
		{
			args: []string{"json",
				"--file", tempLog,
				"--pos", tempPos,
			},
		},
		{
			args: []string{"json",
				"--file", tempLog,
				"--pos", tempPos,
				"--nosave-pos",
			},
		},
		// Do not change the order
		{
			args: []string{"json",
				"--file", tempLog,
				"--dump", tempDump,
			},
		},
		{
			args: []string{"json",
				"--load", tempDump,
			},
		},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			rootCmd := NewRootCmd("test")
			rootCmd.SetArgs(tt.args)

			err := rootCmd.Execute()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
