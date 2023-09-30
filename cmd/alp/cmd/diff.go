package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tkuchiki/alp/parsers"
	"github.com/tkuchiki/alp/profiler"
	"github.com/tkuchiki/alp/stats"
)

func newDiffCmd(flags *flags) *cobra.Command {
	diffCmd := &cobra.Command{
		Use:   "diff <from> <to>",
		Args:  cobra.ExactArgs(2),
		Short: "Show the difference between the two profile results",
		Long:  `Show the difference between the two profile results`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := flags.createDiffOptions(cmd)
			if err != nil {
				return err
			}

			from := args[0]
			to := args[1]

			sts := stats.NewHTTPStats(true, false, false)

			err = sts.InitFilter(opts)
			if err != nil {
				return err
			}

			sts.SetOptions(opts)
			sts.SetSortOptions(flags.sortOptions)

			printOptions := stats.NewPrintOptions(opts.NoHeaders, opts.ShowFooters, opts.DecodeUri, opts.PaginationLimit)
			printer := stats.NewPrinter(os.Stdout, opts.Output, opts.Format, opts.Percentiles, printOptions)
			if err = printer.Validate(); err != nil {
				return err
			}

			fromf, err := os.Open(from)
			if err != nil {
				return err
			}
			err = sts.LoadStats(fromf)
			if err != nil {
				return err
			}
			defer fromf.Close()

			sts.SortWithOptions()

			toSts := stats.NewHTTPStats(true, false, false)
			err = toSts.InitFilter(opts)
			if err != nil {
				return err
			}

			toSts.SetOptions(opts)
			toSts.SetSortOptions(flags.sortOptions)

			tof, err := os.Open(to)
			if err != nil {
				return err
			}
			err = toSts.LoadStats(tof)
			if err != nil {
				return err
			}
			defer tof.Close()

			toSts.SortWithOptions()

			printer.Print(sts, toSts)

			return nil
		},
	}

	flags.defineDiffOptions(diffCmd)

	diffCmd.Flags().SortFlags = false
	diffCmd.PersistentFlags().SortFlags = false
	diffCmd.InheritedFlags().SortFlags = false

	return diffCmd
}

func newDiffSubCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "diff [<from>] <to>",
		Args:  cobra.RangeArgs(1, 2),
		Short: "Show the difference between the two access log_reader",
		Long:  `Show the difference between the two access log_reader`,
	}
}

func getFromTo(load string, args []string) (string, string) {
	if load != "" {
		return "", args[0]
	}

	return args[0], args[1]
}

func runDiff(sortOptions *stats.SortOptions,
	fromProf *profiler.Profiler, fromParser parsers.Parser,
	toProf *profiler.Profiler, toParser parsers.Parser) error {

	fromSts, err := fromProf.Profile(sortOptions, fromParser)
	if err != nil {
		return err
	}

	return toProf.Run(sortOptions, toParser, fromSts)
}
