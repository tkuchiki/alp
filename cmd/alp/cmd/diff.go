package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tkuchiki/alp/stats"
)

func NewDiffCmd(rootCmd *cobra.Command) *cobra.Command {
	var diffCmd = &cobra.Command{
		Use:   "diff <from> <to>",
		Args:  cobra.ExactArgs(2),
		Short: "Show the difference between the two profile results",
		Long:  `Show the difference between the two profile results`,
		RunE: func(cmd *cobra.Command, args []string) error {
			sortOptions := stats.NewSortOptions()
			opts, err := createOptions(rootCmd, sortOptions)
			if err != nil {
				return err
			}

			from, err := cmd.PersistentFlags().GetString("from")
			if err != nil {
				return err
			}

			to, err := cmd.PersistentFlags().GetString("to")
			if err != nil {
				return err
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

	//app.Arg("from", "").Required().StringVar(&f.From)
	//app.Arg("to", "").Required().StringVar(&f.To)

	diffCmd.PersistentFlags().StringP("from", "", "", "The comparison source file")
	diffCmd.PersistentFlags().StringP("to", "", "", "The comparison target file")

	return diffCmd
}
