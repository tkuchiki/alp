package cmd

import (
	"os"

	"github.com/tkuchiki/alp/parsers"
	"github.com/tkuchiki/alp/profiler"

	"github.com/spf13/cobra"
)

func NewLTSVCmd(commandFlags *flags) *cobra.Command {
	var ltsvCmd = &cobra.Command{
		Use:   "ltsv",
		Short: "Profile the logs for LTSV",
		Long:  `Profile the logs for LTSV`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := commandFlags.createLTSVOptions(cmd)
			if err != nil {
				return err
			}

			prof := profiler.NewProfiler(os.Stdout, os.Stderr, opts)

			f, err := prof.Open(opts.File)
			if err != nil {
				return err
			}
			defer f.Close()

			label := parsers.NewLTSVLabel(opts.LTSV.UriLabel, opts.LTSV.MethodLabel, opts.LTSV.TimeLabel,
				opts.LTSV.ApptimeLabel, opts.LTSV.ReqtimeLabel, opts.LTSV.SizeLabel, opts.LTSV.StatusLabel,
			)
			parser := parsers.NewLTSVParser(f, label, opts.QueryString, opts.QueryStringIgnoreValues)

			err = prof.Run(commandFlags.sortOptions, parser)

			return err
		},
	}

	commandFlags.defineProfileOptions(ltsvCmd)
	commandFlags.defineLTSVOptions(ltsvCmd)

	ltsvCmd.Flags().SortFlags = false
	ltsvCmd.PersistentFlags().SortFlags = false
	ltsvCmd.InheritedFlags().SortFlags = false

	return ltsvCmd
}
