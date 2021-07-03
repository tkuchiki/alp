package stats

import (
	"fmt"
	"sort"
)

const (
	SortCount                   = "Count"
	SortUri                     = "Uri"
	SortMethod                  = "Method"
	SortMaxResponseTime         = "MaxResponseTime"
	SortMinResponseTime         = "MinResponseTime"
	SortSumResponseTime         = "SumResponseTime"
	SortAvgResponseTime         = "AvgResponseTime"
	SortPNResponseTime          = "PNResponseTime"
	SortStddevResponseTime      = "StddevResponseTime"
	SortMaxRequestBodyBytes     = "MaxRequestBodyBytes"
	SortMinRequestBodyBytes     = "MinRequestBodyBytes"
	SortSumRequestBodyBytes     = "SumRequestBodyBytes"
	SortAvgRequestBodyBytes     = "AvgRequestBodyBytes"
	SortPNRequestBodyBytes      = "PNRequestBodyBytes"
	SortStddevRequestBodyBytes  = "StddevRequestBodyBytes"
	SortMaxResponseBodyBytes    = "MaxResponseBodyBytes"
	SortMinResponseBodyBytes    = "MinResponseBodyBytes"
	SortSumResponseBodyBytes    = "SumResponseBodyBytes"
	SortAvgResponseBodyBytes    = "AvgResponseBodyBytes"
	SortPNResponseBodyBytes     = "PNResponseBodyBytes"
	SortStddevResponseBodyBytes = "StddevResponseBodyBytes"
)

type SortOptions struct {
	options    map[string]string
	sortType   string
	percentile int
}

func NewSortOptions() *SortOptions {
	options := map[string]string{
		"max":      SortMaxResponseTime,
		"min":      SortMinResponseTime,
		"avg":      SortAvgResponseTime,
		"sum":      SortSumResponseTime,
		"count":    SortCount,
		"uri":      SortUri,
		"method":   SortMethod,
		"max-body": SortMaxResponseBodyBytes,
		"min-body": SortMinResponseBodyBytes,
		"avg-body": SortAvgResponseBodyBytes,
		"sum-body": SortSumResponseBodyBytes,
		"stddev":   SortStddevResponseTime,
		"pn":       SortPNResponseTime,
	}

	return &SortOptions{
		options: options,
	}
}

func (so *SortOptions) SetAndValidate(opt string) error {
	_, ok := so.options[opt]
	if ok {
		so.sortType = so.options[opt]
		return nil
	}

	var n int
	_, err := fmt.Sscanf(opt, "p%d", &n)
	if err != nil {
		return err
	}

	if n < 0 && n > 100 {
		return fmt.Errorf("enum value must be one of max,min,avg,sum,count,uri,method,max-body,min-body,avg-body,sum-body,pN(N = 0 ~ 100),stddev, got '%s'", opt)
	}

	so.sortType = so.options["pn"]
	so.percentile = n

	return nil
}

func (so *SortOptions) SortType() string {
	return so.sortType
}

func (so *SortOptions) Percentile() int {
	return so.percentile
}

func (hs *HTTPStats) Sort(sortOptions *SortOptions, reverse bool) {
	switch sortOptions.sortType {
	case SortCount:
		hs.SortCount(reverse)
	case SortUri:
		hs.SortUri(reverse)
	case SortMethod:
		hs.SortMethod(reverse)
	// response time
	case SortMaxResponseTime:
		hs.SortMaxResponseTime(reverse)
	case SortMinResponseTime:
		hs.SortMinResponseTime(reverse)
	case SortSumResponseTime:
		hs.SortSumResponseTime(reverse)
	case SortAvgResponseTime:
		hs.SortAvgResponseTime(reverse)
	case SortPNResponseTime:
		hs.SortPNResponseTime(reverse)
	case SortStddevResponseTime:
		hs.SortStddevResponseTime(reverse)
	// request body bytes
	case SortMaxRequestBodyBytes:
		hs.SortMaxRequestBodyBytes(reverse)
	case SortMinRequestBodyBytes:
		hs.SortMinRequestBodyBytes(reverse)
	case SortSumRequestBodyBytes:
		hs.SortSumRequestBodyBytes(reverse)
	case SortAvgRequestBodyBytes:
		hs.SortAvgRequestBodyBytes(reverse)
	case SortPNRequestBodyBytes:
		hs.SortPNRequestBodyBytes(reverse)
	case SortStddevRequestBodyBytes:
		hs.SortStddevRequestBodyBytes(reverse)
	// response body bytes
	case SortMaxResponseBodyBytes:
		hs.SortMaxResponseBodyBytes(reverse)
	case SortMinResponseBodyBytes:
		hs.SortMinResponseBodyBytes(reverse)
	case SortSumResponseBodyBytes:
		hs.SortSumResponseBodyBytes(reverse)
	case SortAvgResponseBodyBytes:
		hs.SortAvgResponseBodyBytes(reverse)
	case SortPNResponseBodyBytes:
		hs.SortPNResponseBodyBytes(reverse)
	case SortStddevResponseBodyBytes:
		hs.SortStddevResponseBodyBytes(reverse)
	default:
		hs.SortCount(reverse)
	}
}

func (hs *HTTPStats) SortCount(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].Count() > hs.stats[j].Count()
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].Count() < hs.stats[j].Count()
		})
	}
}

func (hs *HTTPStats) SortUri(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].Uri > hs.stats[j].Uri
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].Uri < hs.stats[j].Uri
		})
	}
}

func (hs *HTTPStats) SortMethod(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].Method > hs.stats[j].Method
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].Method < hs.stats[j].Method
		})
	}
}

func (hs *HTTPStats) SortMaxResponseTime(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].MaxResponseTime() > hs.stats[j].MaxResponseTime()
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].MaxResponseTime() < hs.stats[j].MaxResponseTime()
		})
	}
}

func (hs *HTTPStats) SortMinResponseTime(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].MinResponseTime() > hs.stats[j].MinResponseTime()
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].MinResponseTime() < hs.stats[j].MinResponseTime()
		})
	}
}

func (hs *HTTPStats) SortSumResponseTime(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].SumResponseTime() > hs.stats[j].SumResponseTime()
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].SumResponseTime() < hs.stats[j].SumResponseTime()
		})
	}
}

func (hs *HTTPStats) SortAvgResponseTime(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].AvgResponseTime() > hs.stats[j].AvgResponseTime()
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].AvgResponseTime() < hs.stats[j].AvgResponseTime()
		})
	}
}

func (hs *HTTPStats) SortPNResponseTime(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].PNResponseTime(hs.sortOptions.percentile) > hs.stats[j].PNResponseTime(hs.sortOptions.percentile)
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].PNResponseTime(hs.sortOptions.percentile) < hs.stats[j].PNResponseTime(hs.sortOptions.percentile)
		})
	}
}

func (hs *HTTPStats) SortStddevResponseTime(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].StddevResponseTime() > hs.stats[j].StddevResponseTime()
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].StddevResponseTime() < hs.stats[j].StddevResponseTime()
		})
	}
}

// request
func (hs *HTTPStats) SortMaxRequestBodyBytes(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].MaxRequestBodyBytes() > hs.stats[j].MaxRequestBodyBytes()
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].MaxRequestBodyBytes() < hs.stats[j].MaxRequestBodyBytes()
		})
	}
}

func (hs *HTTPStats) SortMinRequestBodyBytes(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].MinRequestBodyBytes() > hs.stats[j].MinRequestBodyBytes()
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].MinRequestBodyBytes() < hs.stats[j].MinRequestBodyBytes()
		})
	}
}

func (hs *HTTPStats) SortSumRequestBodyBytes(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].SumRequestBodyBytes() > hs.stats[j].SumRequestBodyBytes()
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].SumRequestBodyBytes() < hs.stats[j].SumRequestBodyBytes()
		})
	}
}

func (hs *HTTPStats) SortAvgRequestBodyBytes(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].AvgRequestBodyBytes() > hs.stats[j].AvgRequestBodyBytes()
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].AvgRequestBodyBytes() < hs.stats[j].AvgRequestBodyBytes()
		})
	}
}

func (hs *HTTPStats) SortPNRequestBodyBytes(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].PNRequestBodyBytes(hs.sortOptions.percentile) > hs.stats[j].PNRequestBodyBytes(hs.sortOptions.percentile)
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].PNRequestBodyBytes(hs.sortOptions.percentile) < hs.stats[j].PNRequestBodyBytes(hs.sortOptions.percentile)
		})
	}
}

func (hs *HTTPStats) SortStddevRequestBodyBytes(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].StddevRequestBodyBytes() > hs.stats[j].StddevRequestBodyBytes()
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].StddevRequestBodyBytes() < hs.stats[j].StddevRequestBodyBytes()
		})
	}
}

// response
func (hs *HTTPStats) SortMaxResponseBodyBytes(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].MaxResponseBodyBytes() > hs.stats[j].MaxResponseBodyBytes()
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].MaxResponseBodyBytes() < hs.stats[j].MaxResponseBodyBytes()
		})
	}
}

func (hs *HTTPStats) SortMinResponseBodyBytes(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].MinResponseBodyBytes() > hs.stats[j].MinResponseBodyBytes()
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].MinResponseBodyBytes() < hs.stats[j].MinResponseBodyBytes()
		})
	}
}

func (hs *HTTPStats) SortSumResponseBodyBytes(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].SumResponseBodyBytes() > hs.stats[j].SumResponseBodyBytes()
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].SumResponseBodyBytes() < hs.stats[j].SumResponseBodyBytes()
		})
	}
}

func (hs *HTTPStats) SortAvgResponseBodyBytes(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].AvgResponseBodyBytes() > hs.stats[j].AvgResponseBodyBytes()
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].AvgResponseBodyBytes() < hs.stats[j].AvgResponseBodyBytes()
		})
	}
}

func (hs *HTTPStats) SortPNResponseBodyBytes(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].PNResponseBodyBytes(hs.sortOptions.percentile) > hs.stats[j].PNResponseBodyBytes(hs.sortOptions.percentile)
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].PNResponseBodyBytes(hs.sortOptions.percentile) < hs.stats[j].PNResponseBodyBytes(hs.sortOptions.percentile)
		})
	}
}

func (hs *HTTPStats) SortStddevResponseBodyBytes(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].StddevResponseBodyBytes() > hs.stats[j].StddevResponseBodyBytes()
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].StddevResponseBodyBytes() < hs.stats[j].StddevResponseBodyBytes()
		})
	}
}
