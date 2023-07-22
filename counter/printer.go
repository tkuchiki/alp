package counter

import (
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
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

func NewPrinter(w io.Writer) *Printer {
	return &Printer{
		format: "table",
		writer: w,
	}
}

func (p *Printer) Print(groups *groups) {

	switch p.format {
	case "table":
		p.printTable(groups)
		/*case "md", "markdown":
			p.printMarkdown(groups)
		case "tsv":
			p.printTSV(groups)
		case "csv":
			p.printCSV(groups)
		case "html":
			p.printHTML(groups)*/
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

	/*
		if p.printOptions.showFooters {
			var footer []string
			if hsTo == nil {
				footer = p.GenerateFooter(hsFrom.CountAll())
			} else {
				footer = p.GenerateFooterWithDiff(hsFrom.CountAll(), hsTo.CountAll())
			}
			table.SetFooter(footer)
			table.SetFooterAlignment(tablewriter.ALIGN_LEFT)
		}*/

	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
}

/*
func (p *Printer) printMarkdown(groups) {
	table := tablewriter.NewWriter(p.writer)
	table.SetHeader(p.headers)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	if hsTo == nil {
		for _, s := range hsFrom.stats {
			data := p.GenerateLine(s, false)
			table.Append(data)
		}
	} else {
		for _, to := range hsTo.stats {
			from := findHTTPStatFrom(hsFrom, to)

			var data []string
			if from == nil {
				data = p.GenerateLine(to, false)
			} else {
				data = p.GenerateLineWithDiff(from, to, false)
			}
			table.Append(data)
		}
	}

	if p.printOptions.showFooters {
		var footer []string
		if hsTo == nil {
			footer = p.GenerateFooter(hsFrom.CountAll())
		} else {
			footer = p.GenerateFooterWithDiff(hsFrom.CountAll(), hsTo.CountAll())
		}
		table.Append(footer)
	}

	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
}

func (p *Printer) printTSV(groups) {
	if !p.printOptions.noHeaders {
		fmt.Println(strings.Join(p.headers, "\t"))
	}

	var data []string
	if hsTo == nil {
		for _, s := range hsFrom.stats {
			data = p.GenerateLine(s, false)
			fmt.Println(strings.Join(data, "\t"))
		}
	} else {
		for _, to := range hsTo.stats {
			from := findHTTPStatFrom(hsFrom, to)

			if from == nil {
				data = p.GenerateLine(to, false)
			} else {
				data = p.GenerateLineWithDiff(from, to, false)
			}
			fmt.Println(strings.Join(data, "\t"))
		}
	}
}

func (p *Printer) printCSV(groups) {
	if !p.printOptions.noHeaders {
		fmt.Println(strings.Join(p.headers, ","))
	}

	var data []string
	if hsTo == nil {
		for _, s := range hsFrom.stats {
			data = p.GenerateLine(s, true)
			fmt.Println(strings.Join(data, ","))
		}
	} else {
		for _, to := range hsTo.stats {
			from := findHTTPStatFrom(hsFrom, to)

			if from == nil {
				data = p.GenerateLine(to, false)
			} else {
				data = p.GenerateLineWithDiff(from, to, false)
			}
			fmt.Println(strings.Join(data, ","))
		}
	}
}

func (p *Printer) printHTML(groups) {
	var data [][]string

	if hsTo == nil {
		for _, s := range hsFrom.stats {
			data = append(data, p.GenerateLine(s, true))
		}
	} else {
		for _, to := range hsTo.stats {
			from := findHTTPStatFrom(hsFrom, to)

			if from == nil {
				data = append(data, p.GenerateLine(to, false))
			} else {
				data = append(data, p.GenerateLineWithDiff(from, to, false))
			}
		}
	}
	content, _ := html.RenderTableWithGridJS("alp", p.headers, data, p.printOptions.paginationLimit)
	fmt.Println(content)
}
*/
