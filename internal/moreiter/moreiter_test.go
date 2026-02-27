package moreiter

import (
	"slices"
	"testing"

	"github.com/cszczepaniak/cribbly/internal/assert"
)

func TestOfCollectsValues(t *testing.T) {
	got := slices.Collect(Of(1, 2, 3, 4))
	assert.Equal(t, []int{1, 2, 3, 4}, got)
}

func TestMapTransformsSequence(t *testing.T) {
	got := slices.Collect(Map(Of(1, 2, 3), func(x int) int { return x * 2 }))
	assert.Equal(t, []int{2, 4, 6}, got)
}
