package counter

import "testing"

func TestGroupKeyPreservesEntryBoundaries(t *testing.T) {
	left := groupKey([]string{"a", "bc"})
	right := groupKey([]string{"ab", "c"})

	if left == right {
		t.Fatalf("different groups produced the same key %q", left)
	}
}
