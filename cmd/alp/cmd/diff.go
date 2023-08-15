package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tkuchiki/alp/stats"
)

func NewDiffCmd(rootCmd *cobra.Command) *cobra.Command {
	diffCmd := &cobra.Command{
		Use:   "diff <from> <to>",
		Args:  cobra.ExactArgs(2),
		Short: "Show the difference between the two profile results",
		Long:  `Show the difference between the two profile results`,
		RunE: func(cmd *cobra.Command, args []string) error {
			sortOptions := stats.NewSortOptions()
			opts, err := createOptions(cmd, sortOptions)
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
			sts.SetSortOptions(sortOptions)

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
			toSts.SetSortOptions(sortOptions)

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

	defineOptions(diffCmd)

	return diffCmd
}
