package cmd

import (
	"os"

	"github.com/tkuchiki/alp/options"
	"github.com/tkuchiki/alp/parsers"
	"github.com/tkuchiki/alp/profiler"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tkuchiki/alp/stats"
)

func NewPcapCmd(rootCmd *cobra.Command) *cobra.Command {
	var pcapCmd = &cobra.Command{
		Use:   "pcap",
		Short: "Profile the HTTP requests for captured packets",
		Long:  `Profile the HTTP requests for captured packets`,
		RunE: func(cmd *cobra.Command, args []string) error {
			sortOptions := stats.NewSortOptions()
			opts, err := createPcapOptions(cmd, sortOptions)
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

			err = prof.Run(sortOptions, parser)

			return err
		},
	}

	defineOptions(pcapCmd)

	pcapCmd.PersistentFlags().StringSliceP("pcap-server-ip", "", []string{options.DefaultPcapServerIPsOption[0]}, "HTTP server IP address of the captured packets")
	pcapCmd.PersistentFlags().Uint16P("pcap-server-port", "", options.DefaultPcapServerPortOption, "HTTP server TCP port of the captured packets")

	return pcapCmd
}

func createPcapOptions(cmd *cobra.Command, sortOptions *stats.SortOptions) (*options.Options, error) {
	config, err := cmd.PersistentFlags().GetString("config")
	if err != nil {
		return nil, err
	}
	if config != "" {
		bindCommonFlags(cmd)
		bindPcapFlags(cmd)
		return createOptionsFromConfig(cmd, sortOptions, config)
	}

	opts, err := createCommonOptionsFromFlags(cmd, sortOptions)
	if err != nil {
		return nil, err
	}

	serverIPs, err := cmd.PersistentFlags().GetStringSlice("pcap-server-ip")
	if err != nil {
		return nil, err
	}

	serverPort, err := cmd.PersistentFlags().GetUint16("pcap-server-port")
	if err != nil {
		return nil, err
	}

	return options.SetOptions(opts,
		options.PcapServerIPs(serverIPs),
		options.PcapServerPort(serverPort),
	), nil
}

func bindPcapFlags(cmd *cobra.Command) {
	viper.BindPFlag("pcap.server_ips", cmd.PersistentFlags().Lookup("pcap-server-ip"))
	viper.BindPFlag("pcap.server_port", cmd.PersistentFlags().Lookup("pcap-server-port"))

}
