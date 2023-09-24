package cmd

import (
	"os"

	"github.com/tkuchiki/alp/parsers"
	"github.com/tkuchiki/alp/profiler"

	"github.com/spf13/cobra"
)

func NewJSONCmd(commandFlags *flags) *cobra.Command {
	var jsonCmd = &cobra.Command{
		Use:   "json",
		Short: "Profile the logs for JSON",
		Long:  `Profile the logs for JSON`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := commandFlags.createJSONOptions(cmd)
			if err != nil {
				return err
			}

			prof := profiler.NewProfiler(os.Stdout, os.Stderr, opts)

			f, err := prof.Open(opts.File)
			if err != nil {
				return err
			}
			defer f.Close()

			keys := parsers.NewJSONKeys(opts.JSON.UriKey, opts.JSON.MethodKey, opts.JSON.TimeKey,
				opts.JSON.ResponseTimeKey, opts.JSON.RequestTimeKey, opts.JSON.BodyBytesKey, opts.JSON.StatusKey)
			parser := parsers.NewJSONParser(f, keys, opts.QueryString, opts.QueryStringIgnoreValues)

			err = prof.Run(commandFlags.sortOptions, parser)

			return err
		},
	}

	commandFlags.defineProfileOptions(jsonCmd)
	commandFlags.defineJSONOptions(jsonCmd)

	jsonCmd.Flags().SortFlags = false
	jsonCmd.PersistentFlags().SortFlags = false
	jsonCmd.InheritedFlags().SortFlags = false

	return jsonCmd
}
