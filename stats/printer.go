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

type Printer struct {
	keywords    []string
	format      string
	percentiles []int
	noHeaders   bool
	showFooters bool
	headers     []string
	headersMap  map[string]string
	writer      io.Writer
	all         bool
}

func NewPrinter(w io.Writer, val, format string, percentiles []int, noHeaders, showFooters bool) *Printer {
	p := &Printer{
		format:      format,
		percentiles: percentiles,
		headersMap:  headersMap(percentiles),
		writer:      w,
		showFooters: showFooters,
		noHeaders:   noHeaders,
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

func generateAllLine(s *HTTPStat, percentiles []int, quoteUri bool) []string {
	uri := s.Uri
	if quoteUri && strings.Contains(s.Uri, ",") {
		uri = fmt.Sprintf(`"%s"`, s.Uri)
	}

	l1 := []string{
		s.StrCount(),
		s.StrStatus1xx(),
		s.StrStatus2xx(),
		s.StrStatus3xx(),
		s.StrStatus4xx(),
		s.StrStatus5xx(),
		s.Method,
		uri,
		round(s.MinResponseTime()),
		round(s.MaxResponseTime()),
		round(s.SumResponseTime()),
		round(s.AvgResponseTime()),
	}

	l2 := []string{
		round(s.StddevResponseTime()),
		round(s.MinResponseBodyBytes()),
		round(s.MaxResponseBodyBytes()),
		round(s.SumResponseBodyBytes()),
		round(s.AvgResponseBodyBytes()),
	}

	lp := make([]string, 0, len(percentiles))
	for _, p := range percentiles {
		lp = append(lp, round(s.PNResponseTime(p)))
	}

	l := make([]string, 0, len(l1)+len(l2)+len(lp))
	l = append(l, l1...)
	l = append(l, lp...)
	l = append(l, l2...)

	return l
}

func (p *Printer) GenerateLine(s *HTTPStat, quoteUri bool) []string {
	if p.all {
		return generateAllLine(s, p.percentiles, quoteUri)
	}

	keyLen := len(p.keywords)
	line := make([]string, 0, keyLen)

	for i := 0; i < keyLen; i++ {
		switch p.keywords[i] {
		case "count":
			line = append(line, s.StrCount())
		case "method":
			line = append(line, s.Method)
		case "uri":
			uri := s.Uri
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

func (p *Printer) SetFormat(format string) {
	p.format = format
}

func (p *Printer) SetHeaders(headers []string) {
	p.headers = headers
}

func (p *Printer) SetWriter(w io.Writer) {
	p.writer = w
}

func (p *Printer) Print(hs *HTTPStats) {
	switch p.format {
	case "table":
		p.printTable(hs)
	case "md", "markdown":
		p.printMarkdown(hs)
	case "tsv":
		p.printTSV(hs)
	case "csv":
		p.printCSV(hs)
	}
}

func round(num float64) string {
	return fmt.Sprintf("%.3f", num)
}

func (p *Printer) printTable(hs *HTTPStats) {
	table := tablewriter.NewWriter(p.writer)
	table.SetHeader(p.headers)
	for _, s := range hs.stats {
		data := p.GenerateLine(s, false)
		table.Append(data)
	}

	if p.showFooters {
		footer := p.GenerateFooter(hs.CountAll())
		table.SetFooter(footer)
		table.SetFooterAlignment(tablewriter.ALIGN_RIGHT)
	}

	table.Render()
}

func (p *Printer) printMarkdown(hs *HTTPStats) {
	table := tablewriter.NewWriter(p.writer)
	table.SetHeader(p.headers)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	for _, s := range hs.stats {
		data := p.GenerateLine(s, false)
		table.Append(data)
	}

	if p.showFooters {
		footer := p.GenerateFooter(hs.CountAll())
		table.Append(footer)
	}

	table.Render()
}

func (p *Printer) printTSV(hs *HTTPStats) {
	if !p.noHeaders {
		fmt.Println(strings.Join(p.headers, "\t"))
	}
	for _, s := range hs.stats {
		data := p.GenerateLine(s, false)
		fmt.Println(strings.Join(data, "\t"))
	}
}

func (p *Printer) printCSV(hs *HTTPStats) {
	if !p.noHeaders {
		fmt.Println(strings.Join(p.headers, ","))
	}
	for _, s := range hs.stats {
		data := p.GenerateLine(s, true)
		fmt.Println(strings.Join(data, ","))
	}
}
