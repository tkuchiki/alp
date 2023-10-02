package counter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/tkuchiki/alp/convert"
	"github.com/tkuchiki/alp/html"
)

const (
	defaultSumHeader = "sum"
)

type Printer struct {
	format       string
	printOptions *PrintOptions
	writer       io.Writer
}

type PrintOptions struct {
	noHeaders       bool
	showFooters     bool
	paginationLimit int
}

func NewPrintOptions(noHeaders, showFooters bool, paginationLimit int) *PrintOptions {
	return &PrintOptions{
		noHeaders:       noHeaders,
		showFooters:     showFooters,
		paginationLimit: paginationLimit,
	}
}

func NewPrinter(w io.Writer, format string, printOptions *PrintOptions) *Printer {
	return &Printer{
		format:       format,
		writer:       w,
		printOptions: printOptions,
	}
}

func (p *Printer) Print(groups *groups) {

	switch p.format {
	case "table":
		p.printTable(groups)
	case "md", "markdown":
		p.printMarkdown(groups)
	case "tsv":
		p.printTSV(groups)
	case "csv":
		p.printCSV(groups)
	case "html":
		p.printHTML(groups)
	case "json":
		p.printJSON(groups)
	}
}

func (p *Printer) generateLine(keys []string, group *group) []string {
	var data []string
	data = append(data, fmt.Sprint(group.getCount()))
	for _, key := range keys {
		data = append(data, group.values[key])
	}

	return data
}

func (p *Printer) printTable(groups *groups) {
	table := tablewriter.NewWriter(p.writer)
	var headers []string
	headers = append(headers, defaultSumHeader)
	headers = append(headers, groups.keys...)
	table.SetHeader(headers)

	for _, group := range groups.groups {
		data := p.generateLine(groups.keys, group)
		table.Append(data)
	}

	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
}

func (p *Printer) printMarkdown(groups *groups) {
	table := tablewriter.NewWriter(p.writer)
	var headers []string
	headers = append(headers, defaultSumHeader)
	headers = append(headers, groups.keys...)

	table.SetHeader(headers)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	for _, group := range groups.groups {
		data := p.generateLine(groups.keys, group)
		table.Append(data)
	}

	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
}

func (p *Printer) printTSV(groups *groups) {
	var headers []string
	headers = append(headers, defaultSumHeader)
	headers = append(headers, groups.keys...)

	if !p.printOptions.noHeaders {
		fmt.Println(strings.Join(headers, "\t"))
	}

	for _, group := range groups.groups {
		data := p.generateLine(groups.keys, group)
		fmt.Println(strings.Join(data, "\t"))
	}
}

func (p *Printer) printCSV(groups *groups) {
	var headers []string
	headers = append(headers, defaultSumHeader)
	headers = append(headers, groups.keys...)

	if !p.printOptions.noHeaders {
		fmt.Println(strings.Join(headers, ","))
	}

	for _, group := range groups.groups {
		data := p.generateLine(groups.keys, group)
		fmt.Println(strings.Join(data, ","))
	}
}

func (p *Printer) printHTML(groups *groups) {
	var headers []string
	headers = append(headers, defaultSumHeader)
	headers = append(headers, groups.keys...)

	var data [][]string
	for _, group := range groups.groups {
		data = append(data, p.generateLine(groups.keys, group))
	}

	content, _ := html.RenderTableWithGridJS("alp", headers, data, p.printOptions.paginationLimit)
	fmt.Println(content)
}

func (p *Printer) printJSON(groups *groups) {
	var headers []string
	headers = append(headers, defaultSumHeader)
	headers = append(headers, groups.keys...)

	var data [][]string
	data = append(data, headers)

	for _, group := range groups.groups {
		data = append(data, p.generateLine(groups.keys, group))
	}

	i := convert.ToJSONValues(data)
	b, _ := json.Marshal(i)

	fmt.Println(string(b))
}
