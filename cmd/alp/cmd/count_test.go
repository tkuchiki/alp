package cmd

import (
	"testing"
)

func TestCountCmd(t *testing.T) {
	file := "../../../example/logs/json_access.log"
	args := []string{"count",
		"--file", file,
		"--format", "json",
		"--reverse",
		"--keys", "ua",
	}

	rootCmd := NewRootCmd("test")
	rootCmd.SetArgs(args)

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal(err)
	}
}
