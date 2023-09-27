package cmd

import "testing"

func TestNewRootCmd(t *testing.T) {
	command := NewCommand("test")

	err := command.Execute()
	if err != nil {
		t.Fatal(err)
	}
}
