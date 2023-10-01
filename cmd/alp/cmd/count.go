package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tkuchiki/alp/counter"
	"github.com/tkuchiki/alp/options"
	"github.com/tkuchiki/alp/parsers"
)

func newCountSubCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "count",
		Short: "Count by log entries",
		Long:  `Count by log entries`,
	}
}

func runCount(counter *counter.Counter, parser parsers.Parser, opts *options.Options) error {
	counter.SetParser(parser)
	return counter.CountAndPrint(opts.Count.Keys)
}
