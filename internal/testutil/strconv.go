package testutil

import (
	"fmt"
	"strings"
)

func IntSliceToString(si []int) string {
	return strings.Trim(strings.Replace(fmt.Sprint(si), " ", ",", -1), "[]")
}
