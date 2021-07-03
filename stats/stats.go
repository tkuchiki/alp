package stats

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"sync"

	"github.com/tkuchiki/alp/errors"
	"github.com/tkuchiki/alp/helpers"
	"github.com/tkuchiki/alp/options"
	"github.com/tkuchiki/alp/parsers"
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
	hints                          *hints
	stats                          httpStats
	useResponseTimePercentile      bool
	useRequestBodyBytesPercentile  bool
	useResponseBodyBytesPercentile bool
	filter                         *Filter
	options                        *options.Options
	sortOptions                    *SortOptions
	uriMatchingGroups              []*regexp.Regexp
}

func NewHTTPStats(useResTimePercentile, useRequestBodyBytesPercentile, useResponseBodyBytesPercentile bool) *HTTPStats {
	return &HTTPStats{
		hints:                          newHints(),
		stats:                          make([]*HTTPStat, 0),
		useResponseTimePercentile:      useResTimePercentile,
		useResponseBodyBytesPercentile: useResponseBodyBytesPercentile,
	}
}

func (hs *HTTPStats) Set(uri, method string, status int, restime, resBodyBytes, reqBodyBytes float64) {
	if len(hs.uriMatchingGroups) > 0 {
		for _, re := range hs.uriMatchingGroups {
			if ok := re.Match([]byte(uri)); ok {
				pattern := re.String()
				uri = pattern
			}
		}
	}

	key := fmt.Sprintf("%s_%s", method, uri)

	idx := hs.hints.loadOrStore(key)

	if idx >= len(hs.stats) {
		hs.stats = append(hs.stats, newHTTPStat(uri, method, hs.useResponseTimePercentile, hs.useRequestBodyBytesPercentile, hs.useResponseBodyBytesPercentile))
	}

	hs.stats[idx].Set(status, restime, resBodyBytes, reqBodyBytes)
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

func (hs *HTTPStats) SetSortOptions(options *SortOptions) {
	hs.sortOptions = options
}

func (hs *HTTPStats) SetURIMatchingGroups(groups []string) error {
	uriGroups, err := helpers.CompileUriMatchingGroups(groups)
	if err != nil {
		return err
	}

	hs.uriMatchingGroups = uriGroups

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

func (hs *HTTPStats) CountAll() map[string]int {
	counts := make(map[string]int, 6)

	for _, s := range hs.stats {
		counts["count"] += s.Cnt
		counts["1xx"] += s.Status1xx
		counts["2xx"] += s.Status2xx
		counts["3xx"] += s.Status3xx
		counts["4xx"] += s.Status4xx
		counts["5xx"] += s.Status5xx
	}

	return counts
}

func (hs *HTTPStats) SortWithOptions() {
	hs.Sort(hs.sortOptions, hs.options.Reverse)
}

type HTTPStat struct {
	Uri               string        `yaml:"uri"`
	Cnt               int           `yaml:"count"`
	Status1xx         int           `yaml:"status1xx"`
	Status2xx         int           `yaml:"status2xx"`
	Status3xx         int           `yaml:"status3xx"`
	Status4xx         int           `yaml:"status4xx"`
	Status5xx         int           `yaml:"status5xx"`
	Method            string        `yaml:"method"`
	ResponseTime      *responseTime `yaml:"response_time"`
	RequestBodyBytes  *bodyBytes    `yaml:"request_body_bytes"`
	ResponseBodyBytes *bodyBytes    `yaml:"response_body_bytes"`
	Time              string
}

type httpStats []*HTTPStat

func newHTTPStat(uri, method string, useResTimePercentile, useRequestBodyBytesPercentile, useResponseBodyBytesPercentile bool) *HTTPStat {
	return &HTTPStat{
		Uri:               uri,
		Method:            method,
		ResponseTime:      newResponseTime(useResTimePercentile),
		RequestBodyBytes:  newBodyBytes(useRequestBodyBytesPercentile),
		ResponseBodyBytes: newBodyBytes(useResponseBodyBytesPercentile),
	}
}

func (hs *HTTPStat) Set(status int, restime, reqBodyBytes, resBodyBytes float64) {
	hs.Cnt++
	hs.setStatus(status)
	hs.ResponseTime.Set(restime)
	hs.RequestBodyBytes.Set(reqBodyBytes)
	hs.ResponseBodyBytes.Set(resBodyBytes)
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

func (hs *HTTPStat) PNResponseTime(n int) float64 {
	return hs.ResponseTime.PN(hs.Cnt, n)
}

func (hs *HTTPStat) StddevResponseTime() float64 {
	return hs.ResponseTime.Stddev(hs.Cnt)
}

// request
func (hs *HTTPStat) MaxRequestBodyBytes() float64 {
	return hs.RequestBodyBytes.Max
}

func (hs *HTTPStat) MinRequestBodyBytes() float64 {
	return hs.RequestBodyBytes.Min
}

func (hs *HTTPStat) SumRequestBodyBytes() float64 {
	return hs.RequestBodyBytes.Sum
}

func (hs *HTTPStat) AvgRequestBodyBytes() float64 {
	return hs.RequestBodyBytes.Avg(hs.Cnt)
}

func (hs *HTTPStat) PNRequestBodyBytes(n int) float64 {
	return hs.RequestBodyBytes.PN(hs.Cnt, n)
}

func (hs *HTTPStat) StddevRequestBodyBytes() float64 {
	return hs.RequestBodyBytes.Stddev(hs.Cnt)
}

// response
func (hs *HTTPStat) MaxResponseBodyBytes() float64 {
	return hs.RequestBodyBytes.Max
}

func (hs *HTTPStat) MinResponseBodyBytes() float64 {
	return hs.RequestBodyBytes.Min
}

func (hs *HTTPStat) SumResponseBodyBytes() float64 {
	return hs.RequestBodyBytes.Sum
}

func (hs *HTTPStat) AvgResponseBodyBytes() float64 {
	return hs.RequestBodyBytes.Avg(hs.Cnt)
}

func (hs *HTTPStat) PNResponseBodyBytes(n int) float64 {
	return hs.RequestBodyBytes.PN(hs.Cnt, n)
}

func (hs *HTTPStat) StddevResponseBodyBytes() float64 {
	return hs.RequestBodyBytes.Stddev(hs.Cnt)
}

func percentRank(n int, pi int) int {
	if pi == 0 {
		return 0
	} else if pi == 100 {
		return n - 1
	}

	p := float64(pi) / 100.0
	pos := int(float64(n+1) * p)
	if pos < 0 {
		pos = 0
	}

	return pos - 1
}

type responseTime struct {
	Max           float64 `yaml:"max"`
	Min           float64 `yaml:"min"`
	Sum           float64 `yaml:"sum"`
	UsePercentile bool
	Percentiles   []float64 `yaml:"percentiles"`
}

func newResponseTime(usePercentile bool) *responseTime {
	return &responseTime{
		UsePercentile: usePercentile,
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

	if res.UsePercentile {
		res.Percentiles = append(res.Percentiles, val)
	}
}

func (res *responseTime) Avg(cnt int) float64 {
	return res.Sum / float64(cnt)
}

func (res *responseTime) PN(cnt, n int) float64 {
	if !res.UsePercentile {
		return 0.0
	}

	plen := percentRank(cnt, n)
	res.Sort()
	return res.Percentiles[plen]
}

func (res *responseTime) Stddev(cnt int) float64 {
	if !res.UsePercentile {
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

func (res *responseTime) Sort() {
	sort.Slice(res.Percentiles, func(i, j int) bool {
		return res.Percentiles[i] < res.Percentiles[j]
	})
}

type bodyBytes struct {
	Max           float64 `yaml:"max"`
	Min           float64 `yaml:"min"`
	Sum           float64 `yaml:"sum"`
	UsePercentile bool
	Percentiles   []float64 `yaml:"percentiles"`
}

func newBodyBytes(usePercentile bool) *bodyBytes {
	return &bodyBytes{
		UsePercentile: usePercentile,
		Percentiles:   make([]float64, 0),
	}
}

func (body *bodyBytes) Set(val float64) {
	if body.Max < val {
		body.Max = val
	}

	if body.Min >= val || body.Min == 0.0 {
		body.Min = val
	}

	body.Sum += val

	if body.UsePercentile {
		body.Percentiles = append(body.Percentiles, val)
	}
}

func (body *bodyBytes) Avg(cnt int) float64 {
	return body.Sum / float64(cnt)
}

func (body *bodyBytes) PN(cnt, n int) float64 {
	if !body.UsePercentile {
		return 0.0
	}

	plen := percentRank(cnt, n)
	body.Sort()
	return body.Percentiles[plen]
}

func (body *bodyBytes) Stddev(cnt int) float64 {
	if !body.UsePercentile {
		return 0.0
	}

	var stdd float64
	avg := body.Avg(cnt)
	n := float64(cnt)

	for _, v := range body.Percentiles {
		stdd += (v - avg) * (v - avg)
	}

	return math.Sqrt(stdd / n)
}

func (body *bodyBytes) Sort() {
	sort.Slice(body.Percentiles, func(i, j int) bool {
		return body.Percentiles[i] < body.Percentiles[j]
	})
}
