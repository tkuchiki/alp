package cmd

import "testing"

func TestDiffCmd(t *testing.T) {
	from := "../../../example/logs/dump1.yaml"
	to := "../../../example/logs/dump2.yaml"
	args := []string{"diff",
		from, to,
	}

	command := NewCommand("test")
	command.setArgs(args)

	err := command.Execute()
	if err != nil {
		t.Fatal(err)
	}
}
