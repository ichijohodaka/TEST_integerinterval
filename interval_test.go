package interval

import (
	"fmt"
	"testing"
)

func TestIntervalSet_Union(t *testing.T) {
	set := IntervalSet{{Start: 0, End: 2}, {Start: 5, End: 6}}
	other := IntervalSet{{Start: 1, End: 4}, {Start: 6, End: 8}}
	union := set.Union(other)
	fmt.Println(union.String())
}
