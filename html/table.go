package html

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"text/template"
)

const tplTableWithGridJS = `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>{{.Title}}</title>
<link href="https://unpkg.com/gridjs/dist/theme/mermaid.min.css" rel="stylesheet" />
<script src="https://unpkg.com/gridjs/dist/gridjs.umd.js"></script>
</head>
<body>
<div id="gridContainer"></div>
<script>
new gridjs.Grid({
  {{ .Pagination }}
  columns: [{{range .Columns}}"{{ . }}",{{end}}],
  autoWidth: true,
  style: {
    table: {
      'width': '100%'
    }
  },
  data: {{ .Rows }},
  resizable: true,
  sort: {
    multiColumn: true
  },
  search: true,
}).render(document.getElementById("gridContainer"));
</script>
</body>
</html>`

type ColumnValue struct {
	ColumnType string
	Value      string
}

func RenderTableWithGridJS(title string, columns []string, rows [][]string, paginationLimit int) (string, error) {
	var pagination string
	if paginationLimit > 0 {
		pagination = fmt.Sprintf(`pagination: { limit: %d },`, paginationLimit)
	}
	funcMap := template.FuncMap{
		"escape": func(val string) string { return strings.Replace(val, `"`, `\"`, -1) },
	}

	t, err := template.New("query stats").Funcs(funcMap).Parse(tplTableWithGridJS)
	if err != nil {
		return "", err
	}

	var columnValues [][]interface{}
	for _, row := range rows {
		var values []interface{}
		for _, col := range row {
			vf, err := strconv.ParseFloat(col, 64)
			if err == nil {
				values = append(values, vf)
				continue
			}

			vi, err := strconv.ParseInt(col, 10, 64)
			if err == nil {
				values = append(values, vi)
				continue
			}

			values = append(values, col)
		}
		columnValues = append(columnValues, values)
	}

	s, err := json.Marshal(columnValues)
	if err != nil {
		return "", err
	}

	data := struct {
		Title      string
		Columns    []string
		Rows       string
		Pagination string
	}{
		Title:      title,
		Columns:    columns,
		Rows:       string(s),
		Pagination: pagination,
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	return buf.String(), err
}
