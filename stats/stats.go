package stats

import (
	stderrors "errors"
	"fmt"
	"math"
	"net/url"
	"regexp"
	"sort"
	"sync"
	"time"

	"github.com/tkuchiki/alp/errors"
	"github.com/tkuchiki/alp/helpers"
	"github.com/tkuchiki/alp/options"
	httpv1 "github.com/tkuchiki/logschema/http/v1"
)

type httpStatKey struct {
	method string
	uri    string
}

type httpStatIndex struct {
	values map[httpStatKey]int
	len    int
	mu     sync.Mutex
}

func newHTTPStatIndex() *httpStatIndex {
	return &httpStatIndex{
		values: make(map[httpStatKey]int),
	}
}

func (index *httpStatIndex) loadOrStore(key httpStatKey) int {
	index.mu.Lock()
	defer index.mu.Unlock()
	_, ok := index.values[key]
	if !ok {
		index.values[key] = index.len
		index.len++
	}

	return index.values[key]
}

type HTTPStats struct {
	index                          *httpStatIndex
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
		index:                          newHTTPStatIndex(),
		stats:                          make([]*HTTPStat, 0),
		useResponseTimePercentile:      useResTimePercentile,
		useRequestBodyBytesPercentile:  useRequestBodyBytesPercentile,
		useResponseBodyBytesPercentile: useResponseBodyBytesPercentile,
	}
}

func (hs *HTTPStats) Observe(request *httpv1.Request) {
	uri := request.Data.URLReference()
	if len(hs.uriMatchingGroups) > 0 {
		for _, re := range hs.uriMatchingGroups {
			if ok := re.Match([]byte(uri)); ok {
				pattern := re.String()
				uri = pattern
				break
			}
		}
	}

	key := httpStatKey{method: request.Data.Method, uri: uri}
	idx := hs.index.loadOrStore(key)

	if idx >= len(hs.stats) {
		hs.stats = append(hs.stats, newHTTPStat(uri, request.Data.Method, hs.useResponseTimePercentile, hs.useRequestBodyBytesPercentile, hs.useResponseBodyBytesPercentile))
	}

	hs.stats[idx].Observe(request)
}

func (hs *HTTPStats) Stats() []*HTTPStat {
	return hs.stats
}

func (hs *HTTPStats) CountURIs() int {
	return hs.index.len
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

func (hs *HTTPStats) DoFilter(record *httpv1.Request) (bool, error) {
	err := hs.filter.Do(record)
	if err != nil {
		if stderrors.Is(err, errors.SkipReadLineErr) {
			return false, nil
		}

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
	Uri               string      `yaml:"uri"`
	Cnt               int         `yaml:"count"`
	Status1xx         int         `yaml:"status1xx"`
	Status2xx         int         `yaml:"status2xx"`
	Status3xx         int         `yaml:"status3xx"`
	Status4xx         int         `yaml:"status4xx"`
	Status5xx         int         `yaml:"status5xx"`
	Method            string      `yaml:"method"`
	ResponseTime      *floatStats `yaml:"response_time"`
	RequestBodyBytes  *floatStats `yaml:"request_body_bytes"`
	ResponseBodyBytes *floatStats `yaml:"response_body_bytes"`
	Time              string
}

type httpStats []*HTTPStat

func newHTTPStat(uri, method string, useResTimePercentile, useRequestBodyBytesPercentile, useResponseBodyBytesPercentile bool) *HTTPStat {
	return &HTTPStat{
		Uri:               uri,
		Method:            method,
		ResponseTime:      newFloatStats(useResTimePercentile),
		RequestBodyBytes:  newFloatStats(useRequestBodyBytesPercentile),
		ResponseBodyBytes: newFloatStats(useResponseBodyBytesPercentile),
	}
}

func (hs *HTTPStat) Observe(request *httpv1.Request) {
	hs.Cnt++
	if request.Data.StatusCode != nil {
		hs.setStatus(*request.Data.StatusCode)
	}

	hs.ResponseTime.Set(float64(request.DurationNano) / float64(time.Second))
	if request.Data.RequestBodySizeBytes != nil {
		hs.RequestBodyBytes.Set(float64(*request.Data.RequestBodySizeBytes))
	}
	if request.Data.ResponseBodySizeBytes != nil {
		hs.ResponseBodyBytes.Set(float64(*request.Data.ResponseBodySizeBytes))
	}
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

func (hs *HTTPStat) UriWithOptions(decode bool) string {
	if !decode {
		return hs.Uri
	}

	u, err := url.Parse(hs.Uri)
	if err != nil {
		return hs.Uri
	}

	if u.RawQuery == "" {
		unescaped, _ := url.PathUnescape(u.EscapedPath())
		return unescaped
	}

	unescaped, _ := url.PathUnescape(u.EscapedPath())
	decoded, _ := url.QueryUnescape(u.Query().Encode())

	return fmt.Sprintf("%s?%s", unescaped, decoded)
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

func (hs *HTTPStat) MaxResponseBodyBytes() float64 {
	return hs.ResponseBodyBytes.Max
}

func (hs *HTTPStat) MinResponseBodyBytes() float64 {
	return hs.ResponseBodyBytes.Min
}

func (hs *HTTPStat) SumResponseBodyBytes() float64 {
	return hs.ResponseBodyBytes.Sum
}

func (hs *HTTPStat) AvgResponseBodyBytes() float64 {
	return hs.ResponseBodyBytes.Avg(hs.Cnt)
}

func (hs *HTTPStat) PNResponseBodyBytes(n int) float64 {
	return hs.ResponseBodyBytes.PN(hs.Cnt, n)
}

func (hs *HTTPStat) StddevResponseBodyBytes() float64 {
	return hs.ResponseBodyBytes.Stddev(hs.Cnt)
}

func percentRank(n int, pi int) int {
	switch pi {
	case 0:
		return 0
	case 100:
		return n - 1
	}

	p := float64(pi) / 100.0
	pos := int(float64(n+1) * p)
	if pos <= 0 {
		return 0
	}

	return pos - 1
}

type floatStats struct {
	Max           float64 `yaml:"max"`
	Min           float64 `yaml:"min"`
	Sum           float64 `yaml:"sum"`
	UsePercentile bool
	Percentiles   []float64 `yaml:"percentiles"`

	// SampleCount is a pointer to distinguish legacy dumps from new dumps with
	// zero samples.
	SampleCount *int `yaml:"sample_count,omitempty"`
}

func newFloatStats(usePercentile bool) *floatStats {
	sampleCount := 0

	return &floatStats{
		UsePercentile: usePercentile,
		Percentiles:   make([]float64, 0),
		SampleCount:   &sampleCount,
	}
}

func (stats *floatStats) Set(val float64) {
	if stats.SampleCount == nil {
		sampleCount := 0
		stats.SampleCount = &sampleCount
	}
	(*stats.SampleCount)++

	if stats.Max < val {
		stats.Max = val
	}

	if *stats.SampleCount == 1 || stats.Min > val {
		stats.Min = val
	}

	stats.Sum += val

	if stats.UsePercentile {
		stats.Percentiles = append(stats.Percentiles, val)
	}
}

func (stats *floatStats) Avg(fallbackCount int) float64 {
	cnt := stats.sampleCount(fallbackCount)
	if cnt == 0 {
		return 0
	}

	return stats.Sum / float64(cnt)
}

func (stats *floatStats) PN(_ int, n int) float64 {
	if !stats.UsePercentile || len(stats.Percentiles) == 0 {
		return 0.0
	}

	plen := percentRank(len(stats.Percentiles), n)
	stats.Sort()
	return stats.Percentiles[plen]
}

func (stats *floatStats) Stddev(fallbackCount int) float64 {
	if !stats.UsePercentile {
		return 0.0
	}
	cnt := stats.sampleCount(fallbackCount)
	if cnt == 0 {
		return 0
	}

	var stdd float64
	avg := stats.Avg(fallbackCount)
	n := float64(cnt)

	for _, v := range stats.Percentiles {
		stdd += (v - avg) * (v - avg)
	}

	return math.Sqrt(stdd / n)
}

func (stats *floatStats) Sort() {
	sort.Slice(stats.Percentiles, func(i, j int) bool {
		return stats.Percentiles[i] < stats.Percentiles[j]
	})
}

func (stats *floatStats) sampleCount(fallback int) int {
	if stats.SampleCount != nil {
		return *stats.SampleCount
	}

	return fallback
}
