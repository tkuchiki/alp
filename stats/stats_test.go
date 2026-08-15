package stats

import (
	"fmt"
	"testing"
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
