package stats

import (
	"fmt"
	"io"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/tkuchiki/alp/helpers"
)

func keywords(percentiles []int) []string {
	s1 := []string{
		"count",
		"1xx",
		"2xx",
		"3xx",
		"4xx",
		"5xx",
		"method",
		"uri",
		"min",
		"max",
		"sum",
		"avg",
	}

	s2 := []string{
		"stddev",
		"min_body",
		"max_body",
		"sum_body",
		"avg_body",
	}

	sp := make([]string, 0, len(percentiles))
	for _, p := range percentiles {
		sp = append(sp, fmt.Sprintf("p%d", p))
	}

	s := make([]string, 0, len(s1)+len(s2)+len(sp))
	s = append(s, s1...)
	s = append(s, sp...)
	s = append(s, s2...)

	return s
}

func defaultHeaders(percentiles []int) []string {
	s1 := []string{
		"Count",
		"1xx",
		"2xx",
		"3xx",
		"4xx",
		"5xx",
		"Method",
		"Uri",
		"Min",
		"Max",
		"Sum",
		"Avg",
	}

	s2 := []string{
		"Stddev",
		"Min(Body)",
		"Max(Body)",
		"Sum(Body)",
		"Avg(Body)",
	}

	sp := make([]string, 0, len(percentiles))
	for _, p := range percentiles {
		sp = append(sp, fmt.Sprintf("P%d", p))
	}

	s := make([]string, 0, len(s1)+len(s2)+len(sp))
	s = append(s, s1...)
	s = append(s, sp...)
	s = append(s, s2...)

	return s
}

func headersMap(percentiles []int) map[string]string {
	headers := map[string]string{
		"count":    "Count",
		"1xx":      "1xx",
		"2xx":      "2xx",
		"3xx":      "3xx",
		"4xx":      "4xx",
		"5xx":      "5xx",
		"method":   "Method",
		"uri":      "Uri",
		"min":      "Min",
		"max":      "Max",
		"sum":      "Sum",
		"avg":      "Avg",
		"stddev":   "Stddev",
		"min_body": "Min(Body)",
		"max_body": "Max(Body)",
		"sum_body": "Sum(Body)",
		"avg_body": "Avg(Body)",
	}

	for _, p := range percentiles {
		key := fmt.Sprintf("p%d", p)
		val := fmt.Sprintf("P%d", p)
		headers[key] = val
	}

	return headers
}

type PrintOptions struct {
	noHeaders   bool
	showFooters bool
	decodeUri   bool
}

func NewPrintOptions(noHeaders, showFooters, decodeUri bool) *PrintOptions {
	return &PrintOptions{
		noHeaders:   noHeaders,
		showFooters: showFooters,
		decodeUri:   decodeUri,
	}
}

type Printer struct {
	keywords     []string
	format       string
	percentiles  []int
	printOptions *PrintOptions
	headers      []string
	headersMap   map[string]string
	writer       io.Writer
	all          bool
}

func NewPrinter(w io.Writer, val, format string, percentiles []int, printOptions *PrintOptions) *Printer {
	p := &Printer{
		format:       format,
		percentiles:  percentiles,
		headersMap:   headersMap(percentiles),
		writer:       w,
		printOptions: printOptions,
	}

	if val == "all" {
		p.keywords = keywords(percentiles)
		p.headers = defaultHeaders(percentiles)
		p.all = true
	} else {
		p.keywords = helpers.SplitCSV(val)
		for _, key := range p.keywords {
			p.headers = append(p.headers, p.headersMap[key])
			if key == "all" {
				p.keywords = keywords(percentiles)
				p.headers = defaultHeaders(percentiles)
				p.all = true
				break
			}
		}
	}

	return p
}

func (p *Printer) Validate() error {
	if p.all {
		return nil
	}

	invalids := make([]string, 0)
	for _, key := range p.keywords {
		if _, ok := p.headersMap[key]; !ok {
			invalids = append(invalids, key)
		}
	}

	if len(invalids) > 0 {
		return fmt.Errorf("invalid keywords: %s", strings.Join(invalids, ","))
	}

	return nil
}

func (p *Printer) GenerateLine(s *HTTPStat, quoteUri bool) []string {
	keyLen := len(p.keywords)
	line := make([]string, 0, keyLen)

	for i := 0; i < keyLen; i++ {
		switch p.keywords[i] {
		case "count":
			line = append(line, s.StrCount())
		case "method":
			line = append(line, s.Method)
		case "uri":
			uri := s.UriWithOptions(p.printOptions.decodeUri)
			if quoteUri && strings.Contains(s.Uri, ",") {
				uri = fmt.Sprintf(`"%s"`, s.Uri)
			}
			line = append(line, uri)
		case "1xx":
			line = append(line, s.StrStatus1xx())
		case "2xx":
			line = append(line, s.StrStatus2xx())
		case "3xx":
			line = append(line, s.StrStatus3xx())
		case "4xx":
			line = append(line, s.StrStatus4xx())
		case "5xx":
			line = append(line, s.StrStatus5xx())
		case "min":
			line = append(line, round(s.MinResponseTime()))
		case "max":
			line = append(line, round(s.MaxResponseTime()))
		case "sum":
			line = append(line, round(s.SumResponseTime()))
		case "avg":
			line = append(line, round(s.AvgResponseTime()))
		case "stddev":
			line = append(line, round(s.StddevResponseTime()))
		case "min_body":
			line = append(line, round(s.MinResponseBodyBytes()))
		case "max_body":
			line = append(line, round(s.MaxResponseBodyBytes()))
		case "sum_body":
			line = append(line, round(s.SumResponseBodyBytes()))
		case "avg_body":
			line = append(line, round(s.AvgResponseBodyBytes()))
		default: // percentile
			var n int
			_, err := fmt.Sscanf(p.keywords[i], "p%d", &n)
			if err != nil {
				continue
			}
			line = append(line, round(s.PNResponseTime(n)))
		}
	}

	return line
}

func formattedLineWithDiff(val, diff string) string {
	if diff == "+0" || diff == "+0.000" {
		return val
	}
	return fmt.Sprintf("%s (%s)", val, diff)
}

func (p *Printer) GenerateLineWithDiff(from, to *HTTPStat, quoteUri bool) []string {
	keyLen := len(p.keywords)
	line := make([]string, 0, keyLen)

	differ := NewDiffer(from, to)

	for i := 0; i < keyLen; i++ {
		switch p.keywords[i] {
		case "count":
			line = append(line, formattedLineWithDiff(to.StrCount(), differ.DiffCnt()))
		case "method":
			line = append(line, to.Method)
		case "uri":
			uri := to.UriWithOptions(p.printOptions.decodeUri)
			if quoteUri && strings.Contains(to.Uri, ",") {
				uri = fmt.Sprintf(`"%s"`, to.Uri)
			}
			line = append(line, uri)
		case "1xx":
			line = append(line, formattedLineWithDiff(to.StrStatus1xx(), differ.DiffStatus1xx()))
		case "2xx":
			line = append(line, formattedLineWithDiff(to.StrStatus2xx(), differ.DiffStatus2xx()))
		case "3xx":
			line = append(line, formattedLineWithDiff(to.StrStatus3xx(), differ.DiffStatus3xx()))
		case "4xx":
			line = append(line, formattedLineWithDiff(to.StrStatus4xx(), differ.DiffStatus4xx()))
		case "5xx":
			line = append(line, formattedLineWithDiff(to.StrStatus5xx(), differ.DiffStatus5xx()))
		case "min":
			line = append(line, formattedLineWithDiff(round(to.MinResponseTime()), differ.DiffMinResponseTime()))
		case "max":
			line = append(line, formattedLineWithDiff(round(to.MaxResponseTime()), differ.DiffMaxResponseTime()))
		case "sum":
			line = append(line, formattedLineWithDiff(round(to.SumResponseTime()), differ.DiffSumResponseTime()))
		case "avg":
			line = append(line, formattedLineWithDiff(round(to.AvgResponseTime()), differ.DiffAvgResponseTime()))
		case "stddev":
			line = append(line, formattedLineWithDiff(round(to.StddevResponseTime()), differ.DiffStddevResponseTime()))
		case "min_body":
			line = append(line, formattedLineWithDiff(round(to.MinResponseBodyBytes()), differ.DiffMinResponseBodyBytes()))
		case "max_body":
			line = append(line, formattedLineWithDiff(round(to.MaxResponseBodyBytes()), differ.DiffMaxResponseBodyBytes()))
		case "sum_body":
			line = append(line, formattedLineWithDiff(round(to.SumResponseBodyBytes()), differ.DiffSumResponseBodyBytes()))
		case "avg_body":
			line = append(line, formattedLineWithDiff(round(to.AvgResponseBodyBytes()), differ.DiffAvgResponseBodyBytes()))
		default: // percentile
			var n int
			_, err := fmt.Sscanf(p.keywords[i], "p%d", &n)
			if err != nil {
				continue
			}
			line = append(line, formattedLineWithDiff(round(to.PNResponseTime(n)), differ.DiffPNResponseTime(n)))
		}
	}

	return line
}

func (p *Printer) GenerateFooter(counts map[string]int) []string {
	keyLen := len(p.keywords)
	line := make([]string, 0, keyLen)

	for i := 0; i < keyLen; i++ {
		switch p.keywords[i] {
		case "count":
			line = append(line, fmt.Sprint(counts["count"]))
		case "1xx":
			line = append(line, fmt.Sprint(counts["1xx"]))
		case "2xx":
			line = append(line, fmt.Sprint(counts["2xx"]))
		case "3xx":
			line = append(line, fmt.Sprint(counts["3xx"]))
		case "4xx":
			line = append(line, fmt.Sprint(counts["4xx"]))
		case "5xx":
			line = append(line, fmt.Sprint(counts["5xx"]))
		default:
			line = append(line, "")
		}
	}

	return line
}

func (p *Printer) GenerateFooterWithDiff(countsFrom, countsTo map[string]int) []string {
	keyLen := len(p.keywords)
	line := make([]string, 0, keyLen)
	counts := DiffCountAll(countsFrom, countsTo)

	for i := 0; i < keyLen; i++ {
		switch p.keywords[i] {
		case "count":
			line = append(line, formattedLineWithDiff(fmt.Sprint(countsTo["count"]), counts["count"]))
		case "1xx":
			line = append(line, formattedLineWithDiff(fmt.Sprint(countsTo["1xx"]), counts["1xx"]))
		case "2xx":
			line = append(line, formattedLineWithDiff(fmt.Sprint(countsTo["2xx"]), counts["2xx"]))
		case "3xx":
			line = append(line, formattedLineWithDiff(fmt.Sprint(countsTo["3xx"]), counts["3xx"]))
		case "4xx":
			line = append(line, formattedLineWithDiff(fmt.Sprint(countsTo["4xx"]), counts["4xx"]))
		case "5xx":
			line = append(line, formattedLineWithDiff(fmt.Sprint(countsTo["5xx"]), counts["5xx"]))
		default:
			line = append(line, "")
		}
	}

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

func (p *Printer) Print(hs, hsTo *HTTPStats) {
	switch p.format {
	case "table":
		p.printTable(hs, hsTo)
	case "md", "markdown":
		p.printMarkdown(hs, hsTo)
	case "tsv":
		p.printTSV(hs, hsTo)
	case "csv":
		p.printCSV(hs, hsTo)
	}
}

func round(num float64) string {
	return fmt.Sprintf("%.3f", num)
}

func findHTTPStatFrom(hsFrom *HTTPStats, hsTo *HTTPStat) *HTTPStat {
	for _, sFrom := range hsFrom.stats {
		if sFrom.Uri == hsTo.Uri && sFrom.Method == hsTo.Method {
			return sFrom
		}
	}
	return nil
}

func (p *Printer) printTable(hsFrom, hsTo *HTTPStats) {
	table := tablewriter.NewWriter(p.writer)
	table.SetHeader(p.headers)
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
		table.SetFooter(footer)
		table.SetFooterAlignment(tablewriter.ALIGN_LEFT)
	}

	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
}

func (p *Printer) printMarkdown(hsFrom, hsTo *HTTPStats) {
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

func (p *Printer) printTSV(hsFrom, hsTo *HTTPStats) {
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

func (p *Printer) printCSV(hsFrom, hsTo *HTTPStats) {
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
