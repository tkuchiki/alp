package log_reader

import (
	"fmt"
	"sort"
)

const (
	SortResponseTime = "restime"
	SortBodyBytes    = "bytes"
)

func (a *AccessLogReader) Sort(sortType string, reverse bool) error {
	switch sortType {
	case SortResponseTime:
		a.SortResponseTime(reverse)
	case SortBodyBytes:
		a.SortBodyBytes(reverse)
	default:
		return fmt.Errorf("%s is invalid sort type", sortType)
	}

	return nil
}

func (a *AccessLogReader) SortResponseTime(reverse bool) {
	if reverse {
		sort.Slice(a.logs, func(i, j int) bool {
			return a.logs[i].ResponseTime > a.logs[j].ResponseTime
		})
	} else {
		sort.Slice(a.logs, func(i, j int) bool {
			return a.logs[i].ResponseTime < a.logs[j].ResponseTime
		})
	}
}

func (a *AccessLogReader) SortBodyBytes(reverse bool) {
	if reverse {
		sort.Slice(a.logs, func(i, j int) bool {
			return a.logs[i].BodyBytes > a.logs[j].BodyBytes
		})
	} else {
		sort.Slice(a.logs, func(i, j int) bool {
			return a.logs[i].BodyBytes < a.logs[j].BodyBytes
		})
	}
}
