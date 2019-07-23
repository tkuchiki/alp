package stats

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

func IsIncludedInTime(start, end, val int64) bool {
	if start > 0 && end == 0 {
		return start <= val
	} else if end > 0 && start == 0 {
		return end >= val
	} else if start > 0 && end > 0 {
		return start <= val && end >= val
	}

	return false
}

func StringToFloat64(val string) (float64, error) {
	return strconv.ParseFloat(val, 64)
}

func StringToInt(val string) (int, error) {
	return strconv.Atoi(val)
}
