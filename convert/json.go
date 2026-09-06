package convert

import "encoding/json"

func ToJSONValues(data [][]string) [][]any {
	i := [][]any{}
	for _, values := range data {
		l := []any{}
		for _, val := range values {
			n := json.Number(val)

			i64, err := n.Int64()
			if err == nil {
				l = append(l, i64)
				continue
			}

			f64, err := n.Float64()
			if err == nil {
				l = append(l, f64)
				continue
			}

			l = append(l, val)
		}
		i = append(i, l)
	}

	return i
}
