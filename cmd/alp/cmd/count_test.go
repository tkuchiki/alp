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

	command := NewCommand("test")
	command.setArgs(args)

	err := command.Execute()
	if err != nil {
		t.Fatal(err)
	}
}
