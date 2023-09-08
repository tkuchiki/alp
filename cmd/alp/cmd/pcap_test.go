package cmd

import (
	"testing"

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

	rootCmd := NewRootCmd("test")
	rootCmd.SetArgs(args)

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal(err)
	}
}
