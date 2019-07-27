package stats

import (
	"fmt"
	"math"
	"regexp"
	"sync"

	"github.com/tkuchiki/alp/errors"

	"github.com/tkuchiki/alp/helpers"

	"github.com/tkuchiki/alp/parsers"

	"github.com/tkuchiki/alp/options"
)

type hints struct {
	values map[string]int
	len    int
	mu     sync.RWMutex
}

func newHints() *hints {
	return &hints{
		values: make(map[string]int),
	}
}

func (h *hints) loadOrStore(key string) int {
	h.mu.Lock()
	defer h.mu.Unlock()
	_, ok := h.values[key]
	if !ok {
		h.values[key] = h.len
		h.len++
	}

	return h.values[key]
}

type HTTPStats struct {
	hints                         *hints
	stats                         httpStats
	useResponseTimePercentile     bool
	useRequestBodySizePercentile  bool
	useResponseBodySizePercentile bool
	printOptions                  *PrintOptions
	filter                        *Filter
	options                       *options.Options
	uriCapturingGroups            []*regexp.Regexp
}

func NewHTTPStats(useResTimePercentile, useRequestBodySizePercentile, useResponseBodySizePercentile bool, po *PrintOptions) *HTTPStats {
	return &HTTPStats{
		hints:                         newHints(),
		stats:                         make([]*HTTPStat, 0),
		useResponseTimePercentile:     useResTimePercentile,
		useResponseBodySizePercentile: useResponseBodySizePercentile,
		printOptions:                  po,
	}
}

func (hs *HTTPStats) Set(uri, method string, status int, restime, resBodySize, reqBodySize float64) {
	if len(hs.uriCapturingGroups) > 0 {
		for _, re := range hs.uriCapturingGroups {
			if ok := re.Match([]byte(uri)); ok {
				pattern := re.String()
				uri = pattern
			}
		}
	}

	key := fmt.Sprintf("%s_%s", method, uri)

	idx := hs.hints.loadOrStore(key)

	if idx >= len(hs.stats) {
		hs.stats = append(hs.stats, newHTTPStat(uri, method, hs.useResponseTimePercentile, hs.useRequestBodySizePercentile, hs.useResponseBodySizePercentile))
	}

	hs.stats[idx].Set(status, restime, resBodySize, reqBodySize)
}

func (hs *HTTPStats) Stats() []*HTTPStat {
	return hs.stats
}

func (hs *HTTPStats) CountUris() int {
	return hs.hints.len
}

func (hs *HTTPStats) SetOptions(options *options.Options) {
	hs.options = options
}

func (hs *HTTPStats) SetURICapturingGroups(groups []string) error {
	uriGroups, err := helpers.CompileUriGroups(groups)
	if err != nil {
		return err
	}

	hs.uriCapturingGroups = uriGroups

	return nil
}

func (hs *HTTPStats) InitFilter(options *options.Options) error {
	hs.filter = NewFilter(options)
	return hs.filter.Init()
}

func (hs *HTTPStats) DoFilter(pstat *parsers.ParsedHTTPStat) (bool, error) {
	err := hs.filter.Do(pstat)
	if err == errors.SkipReadLineErr {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (hs *HTTPStats) SortWithOptions() {
	hs.Sort(hs.options.Sort, hs.options.Reverse)
}

type HTTPStat struct {
	Uri              string        `yaml:uri`
	Cnt              int           `yaml:count`
	Status1xx        int           `yaml:status1xx`
	Status2xx        int           `yaml:status2xx`
	Status3xx        int           `yaml:status3xx`
	Status4xx        int           `yaml:status4xx`
	Status5xx        int           `yaml:status5xx`
	Method           string        `yaml:method`
	ResponseTime     *responseTime `yaml:response_time`
	RequestBodySize  *bodySize     `yaml:request_body_size`
	ResponseBodySize *bodySize     `yaml:response_body_size`
	Time             string
}

type httpStats []*HTTPStat

func newHTTPStat(uri, method string, useResTimePercentile, useRequestBodySizePercentile, useResponseBodySizePercentile bool) *HTTPStat {
	return &HTTPStat{
		Uri:              uri,
		Method:           method,
		ResponseTime:     newResponseTime(useResTimePercentile),
		RequestBodySize:  newBodySize(useRequestBodySizePercentile),
		ResponseBodySize: newBodySize(useResponseBodySizePercentile),
	}
}

func (hs *HTTPStat) Set(status int, restime, reqBodySize, resBodySize float64) {
	hs.Cnt++
	hs.setStatus(status)
	hs.ResponseTime.Set(restime)
	hs.RequestBodySize.Set(reqBodySize)
	hs.ResponseBodySize.Set(resBodySize)
}

func (hs *HTTPStat) setStatus(status int) {
	if status >= 100 && status <= 199 {
		hs.Status1xx++
	} else if status >= 200 && status <= 299 {
		hs.Status2xx++
	} else if status >= 300 && status <= 399 {
		hs.Status3xx++
	} else if status >= 400 && status <= 499 {
		hs.Status4xx++
	} else if status >= 500 && status <= 599 {
		hs.Status5xx++
	}
}

func (hs *HTTPStat) StrStatus1xx() string {
	return fmt.Sprint(hs.Status1xx)
}

func (hs *HTTPStat) StrStatus2xx() string {
	return fmt.Sprint(hs.Status2xx)
}

func (hs *HTTPStat) StrStatus3xx() string {
	return fmt.Sprint(hs.Status3xx)
}

func (hs *HTTPStat) StrStatus4xx() string {
	return fmt.Sprint(hs.Status4xx)
}

func (hs *HTTPStat) StrStatus5xx() string {
	return fmt.Sprint(hs.Status5xx)
}

func (hs *HTTPStat) Count() int {
	return hs.Cnt
}

func (hs *HTTPStat) StrCount() string {
	return fmt.Sprint(hs.Cnt)
}

func (hs *HTTPStat) MaxResponseTime() float64 {
	return hs.ResponseTime.Max
}

func (hs *HTTPStat) MinResponseTime() float64 {
	return hs.ResponseTime.Min
}

func (hs *HTTPStat) SumResponseTime() float64 {
	return hs.ResponseTime.Sum
}

func (hs *HTTPStat) AvgResponseTime() float64 {
	return hs.ResponseTime.Avg(hs.Cnt)
}

func (hs *HTTPStat) P1ResponseTime() float64 {
	return hs.ResponseTime.P1(hs.Cnt)
}

func (hs *HTTPStat) P50ResponseTime() float64 {
	return hs.ResponseTime.P50(hs.Cnt)
}

func (hs *HTTPStat) P90ResponseTime() float64 {
	return hs.ResponseTime.P90(hs.Cnt)
}

func (hs *HTTPStat) P99ResponseTime() float64 {
	return hs.ResponseTime.P99(hs.Cnt)
}

func (hs *HTTPStat) StddevResponseTime() float64 {
	return hs.ResponseTime.Stddev(hs.Cnt)
}

// request
func (hs *HTTPStat) MaxRequestBodySize() float64 {
	return hs.RequestBodySize.Max
}

func (hs *HTTPStat) MinRequestBodySize() float64 {
	return hs.RequestBodySize.Min
}

func (hs *HTTPStat) SumRequestBodySize() float64 {
	return hs.RequestBodySize.Sum
}

func (hs *HTTPStat) AvgRequestBodySize() float64 {
	return hs.RequestBodySize.Avg(hs.Cnt)
}

func (hs *HTTPStat) P1RequestBodySize() float64 {
	return hs.RequestBodySize.P1(hs.Cnt)
}

func (hs *HTTPStat) P50RequestBodySize() float64 {
	return hs.RequestBodySize.P50(hs.Cnt)
}

func (hs *HTTPStat) P90RequestBodySize() float64 {
	return hs.RequestBodySize.P90(hs.Cnt)
}

func (hs *HTTPStat) P99RequestBodySize() float64 {
	return hs.RequestBodySize.P99(hs.Cnt)
}

func (hs *HTTPStat) StddevRequestBodySize() float64 {
	return hs.RequestBodySize.Stddev(hs.Cnt)
}

// response
func (hs *HTTPStat) MaxResponseBodySize() float64 {
	return hs.RequestBodySize.Max
}

func (hs *HTTPStat) MinResponseBodySize() float64 {
	return hs.RequestBodySize.Min
}

func (hs *HTTPStat) SumResponseBodySize() float64 {
	return hs.RequestBodySize.Sum
}

func (hs *HTTPStat) AvgResponseBodySize() float64 {
	return hs.RequestBodySize.Avg(hs.Cnt)
}

func (hs *HTTPStat) P1ResponseBodySize() float64 {
	return hs.RequestBodySize.P1(hs.Cnt)
}

func (hs *HTTPStat) P50ResponseBodySize() float64 {
	return hs.RequestBodySize.P50(hs.Cnt)
}

func (hs *HTTPStat) P90ResponseBodySize() float64 {
	return hs.RequestBodySize.P90(hs.Cnt)
}

func (hs *HTTPStat) P99ResponseBodySize() float64 {
	return hs.RequestBodySize.P99(hs.Cnt)
}

func (hs *HTTPStat) StddevResponseBodySize() float64 {
	return hs.RequestBodySize.Stddev(hs.Cnt)
}

func percentRank(l int, n int) int {
	pLen := (l * n / 100) - 1
	if pLen < 0 {
		pLen = 0
	}

	return pLen
}

type responseTime struct {
	Max           float64
	Min           float64
	Sum           float64
	usePercentile bool
	Percentiles   []float64
}

func newResponseTime(usePercentile bool) *responseTime {
	return &responseTime{
		usePercentile: usePercentile,
		Percentiles:   make([]float64, 0),
	}
}

func (res *responseTime) Set(val float64) {
	if res.Max < val {
		res.Max = val
	}

	if res.Min >= val || res.Min == 0 {
		res.Min = val
	}

	res.Sum += val

	if res.usePercentile {
		res.Percentiles = append(res.Percentiles, val)
	}
}

func (res *responseTime) Avg(cnt int) float64 {
	return res.Sum / float64(cnt)
}

func (res *responseTime) P1(cnt int) float64 {
	if !res.usePercentile {
		return 0.0
	}

	plen := percentRank(cnt, 1)
	return res.Percentiles[plen]
}

func (res *responseTime) P50(cnt int) float64 {
	if !res.usePercentile {
		return 0.0
	}

	plen := percentRank(cnt, 50)
	return res.Percentiles[plen]
}

func (res *responseTime) P90(cnt int) float64 {
	if !res.usePercentile {
		return 0.0
	}

	plen := percentRank(cnt, 90)
	return res.Percentiles[plen]
}

func (res *responseTime) P99(cnt int) float64 {
	if !res.usePercentile {
		return 0.0
	}

	plen := percentRank(cnt, 99)
	return res.Percentiles[plen]
}

func (res *responseTime) Stddev(cnt int) float64 {
	if !res.usePercentile {
		return 0.0
	}

	var stdd float64
	avg := res.Avg(cnt)
	n := float64(cnt)

	for _, v := range res.Percentiles {
		stdd += (v - avg) * (v - avg)
	}

	return math.Sqrt(stdd / n)
}

type bodySize struct {
	Max           float64
	Min           float64
	Sum           float64
	usePercentile bool
	percentiles   []float64
}

func newBodySize(usePercentile bool) *bodySize {
	return &bodySize{
		usePercentile: usePercentile,
		percentiles:   make([]float64, 0),
	}
}

func (body *bodySize) Set(val float64) {
	if body.Max < val {
		body.Max = val
	}

	if body.Min >= val || body.Min == 0.0 {
		body.Min = val
	}

	body.Sum += val

	if body.usePercentile {
		body.percentiles = append(body.percentiles, val)
	}
}

func (body *bodySize) Avg(cnt int) float64 {
	return body.Sum / float64(cnt)
}

func (body *bodySize) P1(cnt int) float64 {
	if !body.usePercentile {
		return 0.0
	}

	plen := percentRank(cnt, 1)
	return body.percentiles[plen]
}

func (body *bodySize) P50(cnt int) float64 {
	if !body.usePercentile {
		return 0.0
	}

	plen := percentRank(cnt, 50)
	return body.percentiles[plen]
}

func (body *bodySize) P90(cnt int) float64 {
	if !body.usePercentile {
		return 0.0
	}

	plen := percentRank(cnt, 90)
	return body.percentiles[plen]
}

func (body *bodySize) P99(cnt int) float64 {
	if !body.usePercentile {
		return 0.0
	}

	plen := percentRank(cnt, 99)
	return body.percentiles[plen]
}

func (body *bodySize) Stddev(cnt int) float64 {
	if !body.usePercentile {
		return 0.0
	}

	var stdd float64
	avg := body.Avg(cnt)
	n := float64(cnt)

	for _, v := range body.percentiles {
		stdd += (v - avg) * (v - avg)
	}

	return math.Sqrt(stdd / n)
}
