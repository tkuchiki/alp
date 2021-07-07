package helpers

import (
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func CompileUriMatchingGroups(groups []string) ([]*regexp.Regexp, error) {
	uriMatchingGroups := make([]*regexp.Regexp, 0, len(groups))
	for _, pattern := range groups {
		u, err := url.Parse(pattern)
		if err == nil {
			if u.RawQuery != "" {
				queries := make(map[string][]string)
				for _, q := range strings.Split(u.RawQuery, "&") {
					item := strings.SplitN(q, "=", 2)
					if len(item) > 0 {
						if len(item) == 2 {
							queries[item[0]] = append(queries[item[0]], item[1])
						} else {
							queries[item[0]] = append(queries[item[0]], "")
						}
					}
				}

				keys := make([]string, 0, len(queries))
				for k, _ := range queries {
					keys = append(keys, k)
				}
				sort.Strings(keys)

				sortedQueries := make([]string, 0)
				for _, k := range keys {
					values := make([]string, 0, len(queries[k]))
					for _, q := range queries[k] {
						values = append(values, fmt.Sprintf("%s=%s", k, q))
					}
					sortedQueries = append(sortedQueries, strings.Join(values, "="))
				}
				sortedRawQuery := strings.Join(sortedQueries, "&")
				pattern = fmt.Sprintf("%s?%s", u.Path, sortedRawQuery)
			}
		}

		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		uriMatchingGroups = append(uriMatchingGroups, re)
	}

	return uriMatchingGroups, nil
}

func StringToFloat64(val string) (float64, error) {
	return strconv.ParseFloat(val, 64)
}

func StringToInt(val string) (int, error) {
	return strconv.Atoi(val)
}

func SplitCSV(val string) []string {
	strs := strings.Split(val, ",")
	if len(strs) == 1 && strs[0] == "" {
		return []string{}
	}

	trimedStrs := make([]string, 0, len(strs))

	for _, s := range strs {
		trimedStrs = append(trimedStrs, strings.Trim(s, " "))
	}

	return trimedStrs
}

func SplitCSVIntoInts(val string) ([]int, error) {
	strs := strings.Split(val, ",")
	if len(strs) == 1 && strs[0] == "" {
		return []int{}, nil
	}

	trimedInts := make([]int, 0, len(strs))

	for _, s := range strs {
		i, err := strconv.Atoi(strings.Trim(s, " "))
		if err != nil {
			return []int{}, err
		}
		trimedInts = append(trimedInts, i)
	}

	for _, i := range trimedInts {
		if i < 1 && i > 100 {
			return []int{}, fmt.Errorf(``)
		}
	}

	return trimedInts, nil
}

func ValidatePercentiles(percentiles []int) error {
	if len(percentiles) == 0 {
		return nil
	}

	for _, i := range percentiles {
		if i < 0 && i > 100 {
			return fmt.Errorf(`percentiles allowed 0 to 100`)
		}
	}

	return nil
}
