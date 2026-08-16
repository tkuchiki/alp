package stats

import (
	"bytes"
	"fmt"
	"testing"

	corev1 "github.com/tkuchiki/logschema/core/v1"
	httpv1 "github.com/tkuchiki/logschema/http/v1"
)

func Test_percentRank(t *testing.T) {
	n := 100
	percentiles := make([]int, 0, n)

	for i := 1; i <= n; i++ {
		percentiles = append(percentiles, i)
	}

	for i := 1; i <= n; i++ {
		i := i
		t.Run(fmt.Sprintf("%d percentile", i), func(t *testing.T) {
			p := percentRank(n, i)
			val := percentiles[p]
			if val != i {
				t.Errorf("want: %d, got: %d", i, val)
			}
		})
	}
}

func TestBodyBytesIgnoresMissingValues(t *testing.T) {
	body := newFloatStats(true)
	body.Set(10)

	if got, want := body.Avg(2), float64(10); got != want {
		t.Fatalf("average = %v, want %v", got, want)
	}
	if got, want := body.PN(2, 100), float64(10); got != want {
		t.Fatalf("p100 = %v, want %v", got, want)
	}
	if got := body.Stddev(2); got != 0 {
		t.Fatalf("standard deviation = %v, want 0", got)
	}
}

func TestBodyBytesWithNoSamples(t *testing.T) {
	body := newFloatStats(true)

	if got := body.Avg(1); got != 0 {
		t.Fatalf("average = %v, want 0", got)
	}
	if got := body.PN(1, 99); got != 0 {
		t.Fatalf("p99 = %v, want 0", got)
	}
	if got := body.Stddev(1); got != 0 {
		t.Fatalf("standard deviation = %v, want 0", got)
	}
}

func TestFloatStatsPreservesMeasuredZeroAsMinimum(t *testing.T) {
	stats := newFloatStats(false)
	stats.Set(0)
	stats.Set(10)

	if got := stats.Min; got != 0 {
		t.Fatalf("minimum = %v, want 0", got)
	}
}

func TestNewHTTPStatsKeepsRequestBodyPercentileOption(t *testing.T) {
	stats := NewHTTPStats(false, true, false)

	if !stats.useRequestBodyBytesPercentile {
		t.Fatal("request body percentile option was not enabled")
	}
}

func TestHTTPStatsDumpPreservesMissingMetricSampleCount(t *testing.T) {
	status := 200
	bodySize := corev1.DecimalUint64(10)
	withBody, err := httpv1.NewRequest(100_000_000, corev1.Source{Kind: corev1.SourceOther}, httpv1.RequestData{
		Method:                "GET",
		URLPath:               "/mixed",
		StatusCode:            &status,
		ResponseBodySizeBytes: &bodySize,
	})
	if err != nil {
		t.Fatal(err)
	}
	withoutBody, err := httpv1.NewRequest(200_000_000, corev1.Source{Kind: corev1.SourceOther}, httpv1.RequestData{
		Method:     "GET",
		URLPath:    "/mixed",
		StatusCode: &status,
	})
	if err != nil {
		t.Fatal(err)
	}

	before := NewHTTPStats(false, false, false)
	before.Observe(&withBody)
	before.Observe(&withoutBody)

	var dump bytes.Buffer
	if err := before.DumpStats(&dump); err != nil {
		t.Fatal(err)
	}
	after := NewHTTPStats(false, false, false)
	if err := after.LoadStats(&dump); err != nil {
		t.Fatal(err)
	}

	if got, want := before.Stats()[0].AvgResponseBodyBytes(), float64(10); got != want {
		t.Fatalf("average before dump = %v, want %v", got, want)
	}
	if got, want := after.Stats()[0].AvgResponseBodyBytes(), float64(10); got != want {
		t.Fatalf("average after load = %v, want %v", got, want)
	}
}
