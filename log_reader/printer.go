package log_reader

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/tkuchiki/alp/convert"
	"github.com/tkuchiki/alp/html"
)

var (
	headerKeys = []string{
		"rank",
		"uri",
		"method",
		"status",
		"restime",
		"bytes",
		"time",
	}

	headers = []string{
		"",
		"Uri",
		"Method",
		"Status",
		"Response Time",
		"Body Bytes",
		"Time",
	}

	headersMap = map[string]string{
		"rank":    "",
		"uri":     "Uri",
		"method":  "Method",
		"status":  "Status",
		"restime": "Response Time",
		"bytes":   "Body Bytes",
		"time":    "Time",
	}
)

type PrintOptions struct {
	noHeaders       bool
	decodeUri       bool
	paginationLimit int
}

func NewPrintOptions(noHeaders, decodeUri bool, paginationLimit int) *PrintOptions {
	return &PrintOptions{
		noHeaders:       noHeaders,
		decodeUri:       decodeUri,
		paginationLimit: paginationLimit,
	}
}

type Printer struct {
	format       string
	printOptions *PrintOptions
	headerKeys   []string
	headers      []string
	headersMap   map[string]string
	writer       io.Writer
	all          bool
	rank         int64
}

func NewPrinter(w io.Writer, format string, printOptions *PrintOptions) *Printer {
	return &Printer{
		format:       format,
		headers:      headers,
		headerKeys:   headerKeys,
		headersMap:   headersMap,
		writer:       w,
		printOptions: printOptions,
		rank:         1,
	}
}

func (p *Printer) Validate() error {
	if p.all {
		return nil
	}

	invalids := make([]string, 0)
	for _, key := range p.headerKeys {
		if _, ok := p.headersMap[key]; !ok {
			invalids = append(invalids, key)
		}
	}

	if len(invalids) > 0 {
		return fmt.Errorf("invalid keywords: %s", strings.Join(invalids, ","))
	}

	return nil
}

func (p *Printer) GenerateLine(l *AccessLog, quoteUri bool) []string {
	keyLen := len(p.headerKeys)
	line := make([]string, 0, keyLen)

	line = append(line, fmt.Sprint(p.rank))

	for i := 0; i < keyLen; i++ {
		switch p.headerKeys[i] {
		case "uri":
			uri := l.UriWithOptions(p.printOptions.decodeUri)
			if quoteUri && strings.Contains(l.Uri, ",") {
				uri = fmt.Sprintf(`"%s"`, l.Uri)
			}
			line = append(line, uri)
		case "method":
			line = append(line, l.Method)
		case "status":
			line = append(line, fmt.Sprint(l.Status))
		case "restime":
			line = append(line, round(l.ResponseTime))
		case "bytes":
			line = append(line, round(l.BodyBytes))
		case "time":
			line = append(line, l.TimeStr)
		}
	}

	p.rank++

	return line
}

func (p *Printer) SetFormat(format string) {
	p.format = format
}

func (p *Printer) SetHeaders(headers []string) {
	p.headers = headers
}

func (p *Printer) SetWriter(w io.Writer) {
	p.writer = w
}

func (p *Printer) Print(logs []*AccessLog) {
	switch p.format {
	case "table":
		p.printTable(logs)
	case "md", "markdown":
		p.printMarkdown(logs)
	case "tsv":
		p.printTSV(logs)
	case "csv":
		p.printCSV(logs)
	case "html":
		p.printHTML(logs)
	case "json":
		p.printJSON(logs)
	}
}

func round(num float64) string {
	return fmt.Sprintf("%.3f", num)
}

func (p *Printer) printTable(logs []*AccessLog) {
	table := tablewriter.NewWriter(p.writer)
	table.SetHeader(p.headers)

	for _, l := range logs {
		data := p.GenerateLine(l, false)
		table.Append(data)
	}

	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
}

func (p *Printer) printMarkdown(logs []*AccessLog) {
	table := tablewriter.NewWriter(p.writer)
	table.SetHeader(p.headers)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	for _, l := range logs {
		data := p.GenerateLine(l, false)
		table.Append(data)
	}

	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
}

func (p *Printer) printTSV(logs []*AccessLog) {
	if !p.printOptions.noHeaders {
		fmt.Println(strings.Join(p.headers, "\t"))
	}

	for _, l := range logs {
		data := p.GenerateLine(l, false)
		fmt.Println(strings.Join(data, "\t"))
	}

}

func (p *Printer) printCSV(logs []*AccessLog) {
	if !p.printOptions.noHeaders {
		fmt.Println(strings.Join(p.headers, ","))
	}

	for _, l := range logs {
		data := p.GenerateLine(l, true)
		fmt.Println(strings.Join(data, ","))
	}
}

func (p *Printer) printHTML(logs []*AccessLog) {
	var data [][]string

	for _, l := range logs {
		data = append(data, p.GenerateLine(l, true))
	}

	content, _ := html.RenderTableWithGridJS("alp", p.headers, data, p.printOptions.paginationLimit)
	fmt.Println(content)
}

func (p *Printer) printJSON(logs []*AccessLog) {
	var data [][]string
	data = append(data, headerKeys)

	for _, l := range logs {
		data = append(data, p.GenerateLine(l, true))
	}

	i := convert.ToJSONValues(data)
	b, _ := json.Marshal(i)

	fmt.Println(string(b))
}
