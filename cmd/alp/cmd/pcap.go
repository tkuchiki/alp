package cmd

import (
	"os"

	"github.com/tkuchiki/alp/options"
	"github.com/tkuchiki/alp/parsers"
	"github.com/tkuchiki/alp/profiler"

	"github.com/spf13/cobra"
	"github.com/tkuchiki/alp/stats"
)

func NewPcapCmd(rootCmd *cobra.Command) *cobra.Command {
	var pcapCmd = &cobra.Command{
		Use:   "pcap",
		Short: "Profile the HTTP requests for captured packets",
		Long:  `Profile the HTTP requests for captured packets`,
		RunE: func(cmd *cobra.Command, args []string) error {
			sortOptions := stats.NewSortOptions()
			opts, err := createOptions(cmd, sortOptions)
			if err != nil {
				return err
			}

			serverIPs, err := cmd.PersistentFlags().GetStringSlice("pcap-server-ip")
			if err != nil {
				return err
			}

			serverPort, err := cmd.PersistentFlags().GetUint16("pcap-server-port")
			if err != nil {
				return err
			}

			opts = options.SetOptions(opts,
				options.PcapServerIPs(serverIPs),
				options.PcapServerPort(serverPort),
			)

			prof := profiler.NewProfiler(os.Stdout, os.Stderr, opts)

			f, err := prof.Open(opts.File)
			if err != nil {
				return err
			}
			defer f.Close()

			parser, err := parsers.NewPcapParser(f, opts.Pcap.ServerIPs, opts.Pcap.ServerPort, opts.QueryString, opts.QueryStringIgnoreValues)
			if err != nil {
				return err
			}

			err = prof.Run(sortOptions, parser)

			return err
		},
	}

	defineOptions(pcapCmd)

	pcapCmd.PersistentFlags().StringSliceP("pcap-server-ip", "", []string{options.DefaultPcapServerIPsOption[0]}, "HTTP server IP address of the captured packets")
	pcapCmd.PersistentFlags().Uint16P("pcap-server-port", "", options.DefaultPcapServerPortOption, "HTTP server TCP port of the captured packets")

	return pcapCmd
}
