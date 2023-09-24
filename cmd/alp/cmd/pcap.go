package cmd

import (
	"os"

	"github.com/tkuchiki/alp/parsers"
	"github.com/tkuchiki/alp/profiler"

	"github.com/spf13/cobra"
)

func NewPcapCmd(commandFlags *flags) *cobra.Command {
	var pcapCmd = &cobra.Command{
		Use:   "pcap",
		Short: "Profile the HTTP requests for captured packets",
		Long:  `Profile the HTTP requests for captured packets`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := commandFlags.createPcapOptions(cmd)
			if err != nil {
				return err
			}

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

			err = prof.Run(commandFlags.sortOptions, parser)

			return err
		},
	}

	commandFlags.defineProfileOptions(pcapCmd)
	commandFlags.definePcapOptions(pcapCmd)

	pcapCmd.Flags().SortFlags = false
	pcapCmd.PersistentFlags().SortFlags = false
	pcapCmd.InheritedFlags().SortFlags = false

	return pcapCmd
}
