package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tkuchiki/alp/stats"
)

func NewDiffCmd(rootCmd *cobra.Command) *cobra.Command {
	diffCmd := &cobra.Command{
		Use:   "diff <from> <to>",
		Args:  cobra.RangeArgs(0, 2),
		Short: "Show the difference between the two profile results",
		Long:  `Show the difference between the two profile results`,
		RunE: func(cmd *cobra.Command, args []string) error {
			sortOptions := stats.NewSortOptions()
			opts, err := createOptions(cmd, sortOptions)
			if err != nil {
				return err
			}

			from, err := cmd.PersistentFlags().GetString("from")
			if err != nil {
				return err
			}
			if from == "" {
				if len(args) < 1 {
					return fmt.Errorf("from is required")
				}
				from = args[0]
			}

			to, err := cmd.PersistentFlags().GetString("to")
			if err != nil {
				return err
			}
			if to == "" {
				if len(args) < 2 {
					return fmt.Errorf("to is required")
				}
				to = args[1]
			}

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

	// app.Arg("from", "").Required().StringVar(&f.From)
	// app.Arg("to", "").Required().StringVar(&f.To)

	defineOptions(diffCmd)

	diffCmd.PersistentFlags().StringP("from", "", "", "The comparison source file")
	diffCmd.PersistentFlags().StringP("to", "", "", "The comparison target file")

	return diffCmd
}
