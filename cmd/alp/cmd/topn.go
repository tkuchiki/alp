package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tkuchiki/alp/helpers"
	"github.com/tkuchiki/alp/log_reader"
	"github.com/tkuchiki/alp/parsers"
)

const (
	defaultTopN = 50
)

func newTopNSubCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "topN [<number>]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Show the top N entries of the access log",
		Long:  `Show the top N entries of the access log`,
	}
}

func getN(args []string) (int, error) {
	if len(args) == 0 {
		return defaultTopN, nil
	}

	return helpers.StringToInt(args[0])
}

func runTopN(logReader *log_reader.AccessLogReader, parser parsers.Parser) error {
	err := logReader.ReadAll(parser)
	if err != nil {
		return err
	}

	logReader.Print()

	return nil
}
