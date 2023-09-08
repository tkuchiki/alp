package cmd

import "testing"

func TestNewRootCmd(t *testing.T) {
	rootCmd := NewRootCmd("test")

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal(err)
	}
}
