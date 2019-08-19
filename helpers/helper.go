package helpers

import (
	"regexp"
	"strconv"
	"strings"
)

func CompileUriMatchingGroups(groups []string) ([]*regexp.Regexp, error) {
	uriMatchingGroups := make([]*regexp.Regexp, 0, len(groups))
	for _, pattern := range groups {
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
