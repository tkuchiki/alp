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
