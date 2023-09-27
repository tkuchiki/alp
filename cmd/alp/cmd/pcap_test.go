package cmd

import (
	"testing"

	"github.com/tkuchiki/alp/internal/testutil"

	"github.com/tkuchiki/alp/options"
)

func TestPcapCmd(t *testing.T) {
	pcapFile := "../../../example/logs/http.cap"
	pcapServerPort := "18080"

	args := []string{"pcap",
		"--file", pcapFile,
		"--pcap-server-ip", options.DefaultPcapServerIPsOption[0],
		"--pcap-server-port", pcapServerPort,
	}

	command := NewCommand("test")
	command.setArgs(args)

	err := command.Execute()
	if err != nil {
		t.Fatal(err)
	}
}

func TestPcapDiffCmd(t *testing.T) {
	pcapFile := "../../../example/logs/http.cap"
	pcapServerPort := "18080"

	tempDir := t.TempDir()

	tempDump, err := testutil.CreateTempDirAndFile(tempDir, "test_pcap_diff_cmd_temp_dump", "")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("alp pcap diff <from> <to>", func(t *testing.T) {
		args := []string{"pcap", "diff",
			pcapFile,
			pcapFile,
			"--dump", tempDump,
			"--pcap-server-ip", options.DefaultPcapServerIPsOption[0],
			"--pcap-server-port", pcapServerPort,
		}

		command := NewCommand("test")
		command.setArgs(args)

		err = command.Execute()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("alp pcap diff --load <dumpfile> <to>", func(t *testing.T) {
		args := []string{"pcap", "diff",
			"--load", tempDump,
			pcapFile,
			"--pcap-server-ip", options.DefaultPcapServerIPsOption[0],
			"--pcap-server-port", pcapServerPort,
		}

		command := NewCommand("test")
		command.setArgs(args)

		err = command.Execute()
		if err != nil {
			t.Fatal(err)
		}
	})
}
