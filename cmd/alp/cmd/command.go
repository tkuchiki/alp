package cmd

import (
	"github.com/spf13/cobra"
)

type Command struct {
	// alp
	rootCmd *cobra.Command

	// alp diff
	diffCmd *cobra.Command

	// alp count
	countCmd *cobra.Command

	// alp json
	jsonCmd *cobra.Command
	// alp json diff
	jsonDiffCmd *cobra.Command

	// alp ltsv
	ltsvCmd *cobra.Command
	// alp ltsv diff
	ltsvDiffCmd *cobra.Command

	// alp regexp
	regexpCmd *cobra.Command
	// alp regexp diff
	regexpDiffCmd *cobra.Command

	// alp pcap
	pcapCmd *cobra.Command
	// alp pcap diff
	pcapDiffCmd *cobra.Command

	flags *flags
}

func NewCommand(version string) *Command {
	command := &Command{}
	command.flags = newFlags()

	command.rootCmd = newRootCmd(version)

	command.flags.defineGlobalOptions(command.rootCmd)

	// alp ltsv
	command.ltsvCmd = newLTSVCmd(command.flags)
	command.rootCmd.AddCommand(command.ltsvCmd)
	// alp ltsv diff
	command.ltsvDiffCmd = newLTSVDiffCmd(command.flags)
	command.ltsvCmd.AddCommand(command.ltsvDiffCmd)

	// alp json
	command.jsonCmd = newJSONCmd(command.flags)
	command.rootCmd.AddCommand(command.jsonCmd)
	// alp json diff
	command.jsonDiffCmd = newJsonDiffCmd(command.flags)
	command.jsonCmd.AddCommand(command.jsonDiffCmd)

	// alp regexp
	command.regexpCmd = newRegexpCmd(command.flags)
	command.rootCmd.AddCommand(command.regexpCmd)
	// alp regexp diff
	command.regexpDiffCmd = newRegexpDiffCmd(command.flags)
	command.regexpCmd.AddCommand(command.regexpDiffCmd)

	// alp pcap
	command.pcapCmd = newPcapCmd(command.flags)
	command.rootCmd.AddCommand(command.pcapCmd)
	// alp pcap diff
	command.pcapDiffCmd = newPcapDiffCmd(command.flags)
	command.pcapCmd.AddCommand(command.pcapDiffCmd)

	// alp diff
	command.diffCmd = newDiffCmd(command.flags)
	command.rootCmd.AddCommand(command.diffCmd)

	// alp count
	command.countCmd = newCountCmd(command.flags)
	command.rootCmd.AddCommand(command.countCmd)

	return command
}

func (c *Command) Execute() error {
	return c.rootCmd.Execute()
}

func (c *Command) setArgs(args []string) {
	c.rootCmd.SetArgs(args)
}
