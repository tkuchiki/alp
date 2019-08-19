package stats

import "sort"

const (
	SortCount                   = "Count"
	SortUri                     = "Uri"
	SortMethod                  = "Method"
	SortMaxResponseTime         = "MaxResponseTime"
	SortMinResponseTime         = "MinResponseTime"
	SortSumResponseTime         = "SumResponseTime"
	SortAvgResponseTime         = "AvgResponseTime"
	SortP1ResponseTime          = "P1ResponseTime"
	SortP50ResponseTime         = "P50ResponseTime"
	SortP90ResponseTime         = "P90ResponseTime"
	SortP99ResponseTime         = "P99ResponseTime"
	SortStddevResponseTime      = "StddevResponseTime"
	SortMaxRequestBodyBytes     = "MaxRequestBodyBytes"
	SortMinRequestBodyBytes     = "MinRequestBodyBytes"
	SortSumRequestBodyBytes     = "SumRequestBodyBytes"
	SortAvgRequestBodyBytes     = "AvgRequestBodyBytes"
	SortP1RequestBodyBytes      = "P1RequestBodyBytes"
	SortP50RequestBodyBytes     = "P50RequestBodyBytes"
	SortP90RequestBodyBytes     = "P90RequestBodyBytes"
	SortP99RequestBodyBytes     = "P99RequestBodyBytes"
	SortStddevRequestBodyBytes  = "StddevRequestBodyBytes"
	SortMaxResponseBodyBytes    = "MaxResponseBodyBytes"
	SortMinResponseBodyBytes    = "MinResponseBodyBytes"
	SortSumResponseBodyBytes    = "SumResponseBodyBytes"
	SortAvgResponseBodyBytes    = "AvgResponseBodyBytes"
	SortP1ResponseBodyBytes     = "P1ResponseBodyBytes"
	SortP50ResponseBodyBytes    = "P50ResponseBodyBytes"
	SortP90ResponseBodyBytes    = "P90ResponseBodyBytes"
	SortP99ResponseBodyBytes    = "P99ResponseBodyBytes"
	SortStddevResponseBodyBytes = "StddevResponseBodyBytes"
)

func (hs *HTTPStats) Sort(sortType string, reverse bool) {
	switch sortType {
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
	case SortP1ResponseTime:
		hs.SortP1ResponseTime(reverse)
	case SortP50ResponseTime:
		hs.SortP50ResponseTime(reverse)
	case SortP90ResponseTime:
		hs.SortP90ResponseTime(reverse)
	case SortP99ResponseTime:
		hs.SortP99ResponseTime(reverse)
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
	case SortP1RequestBodyBytes:
		hs.SortP1RequestBodyBytes(reverse)
	case SortP50RequestBodyBytes:
		hs.SortP50RequestBodyBytes(reverse)
	case SortP90RequestBodyBytes:
		hs.SortP90RequestBodyBytes(reverse)
	case SortP99RequestBodyBytes:
		hs.SortP99RequestBodyBytes(reverse)
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
	case SortP1ResponseBodyBytes:
		hs.SortP1ResponseBodyBytes(reverse)
	case SortP50ResponseBodyBytes:
		hs.SortP50ResponseBodyBytes(reverse)
	case SortP90ResponseBodyBytes:
		hs.SortP90ResponseBodyBytes(reverse)
	case SortP99ResponseBodyBytes:
		hs.SortP99ResponseBodyBytes(reverse)
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

func (hs *HTTPStats) SortP1ResponseTime(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].P1ResponseTime() > hs.stats[j].P1ResponseTime()
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].P1ResponseTime() < hs.stats[j].P1ResponseTime()
		})
	}
}

func (hs *HTTPStats) SortP50ResponseTime(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].P50ResponseTime() > hs.stats[j].P50ResponseTime()
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].P50ResponseTime() < hs.stats[j].P50ResponseTime()
		})
	}
}

func (hs *HTTPStats) SortP90ResponseTime(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].P90ResponseTime() > hs.stats[j].P90ResponseTime()
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].P90ResponseTime() < hs.stats[j].P90ResponseTime()
		})
	}
}

func (hs *HTTPStats) SortP99ResponseTime(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].P99ResponseTime() > hs.stats[j].P99ResponseTime()
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].P99ResponseTime() < hs.stats[j].P99ResponseTime()
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

func (hs *HTTPStats) SortP1RequestBodyBytes(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].P1RequestBodyBytes() > hs.stats[j].P1RequestBodyBytes()
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].P1RequestBodyBytes() < hs.stats[j].P1RequestBodyBytes()
		})
	}
}

func (hs *HTTPStats) SortP50RequestBodyBytes(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].P50RequestBodyBytes() > hs.stats[j].P50RequestBodyBytes()
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].P50RequestBodyBytes() < hs.stats[j].P50RequestBodyBytes()
		})
	}
}

func (hs *HTTPStats) SortP90RequestBodyBytes(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].P90RequestBodyBytes() > hs.stats[j].P90RequestBodyBytes()
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].P90RequestBodyBytes() < hs.stats[j].P90RequestBodyBytes()
		})
	}
}

func (hs *HTTPStats) SortP99RequestBodyBytes(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].P99RequestBodyBytes() > hs.stats[j].P99RequestBodyBytes()
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].P99RequestBodyBytes() < hs.stats[j].P99RequestBodyBytes()
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

func (hs *HTTPStats) SortP1ResponseBodyBytes(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].P1ResponseBodyBytes() > hs.stats[j].P1ResponseBodyBytes()
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].P1ResponseBodyBytes() < hs.stats[j].P1ResponseBodyBytes()
		})
	}
}

func (hs *HTTPStats) SortP50ResponseBodyBytes(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].P50ResponseBodyBytes() > hs.stats[j].P50ResponseBodyBytes()
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].P50ResponseBodyBytes() < hs.stats[j].P50ResponseBodyBytes()
		})
	}
}

func (hs *HTTPStats) SortP90ResponseBodyBytes(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].P90ResponseBodyBytes() > hs.stats[j].P90ResponseBodyBytes()
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].P90ResponseBodyBytes() < hs.stats[j].P90ResponseBodyBytes()
		})
	}
}

func (hs *HTTPStats) SortP99ResponseBodyBytes(reverse bool) {
	if reverse {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].P99ResponseBodyBytes() > hs.stats[j].P99ResponseBodyBytes()
		})
	} else {
		sort.Slice(hs.stats, func(i, j int) bool {
			return hs.stats[i].P99ResponseBodyBytes() < hs.stats[j].P99ResponseBodyBytes()
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
