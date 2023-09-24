package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewRootCmd(version string) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:           "alp",
		Short:         "Access Log Profiler",
		Long:          `alp is the access log profiler for LTSV, JSON, Pcap, and others.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       version,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Help()
			}

			return nil
		},
	}

	commandFlags := newFlags()

	commandFlags.defineGlobalOptions(rootCmd)

	rootCmd.AddCommand(NewLTSVCmd(commandFlags))
	rootCmd.AddCommand(NewJSONCmd(commandFlags))
	rootCmd.AddCommand(NewRegexpCmd(commandFlags))
	rootCmd.AddCommand(NewPcapCmd(commandFlags))
	rootCmd.AddCommand(NewDiffCmd(commandFlags))
	rootCmd.AddCommand(NewCountCmd(commandFlags))
	rootCmd.SetVersionTemplate(fmt.Sprintln(version))

	return rootCmd
}

func Execute(version string) error {
	rootCmd := NewRootCmd(version)
	return rootCmd.Execute()
}
