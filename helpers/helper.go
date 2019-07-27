package helpers

import (
	"regexp"
	"strconv"
)

func CompileUriGroups(groups []string) ([]*regexp.Regexp, error) {
	uriGroups := make([]*regexp.Regexp, 0, len(groups))
	for _, pattern := range groups {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return []*regexp.Regexp{}, err
		}
		uriGroups = append(uriGroups, re)
	}

	return uriGroups, nil
}

func StringToFloat64(val string) (float64, error) {
	return strconv.ParseFloat(val, 64)
}

func StringToInt(val string) (int, error) {
	return strconv.Atoi(val)
}
