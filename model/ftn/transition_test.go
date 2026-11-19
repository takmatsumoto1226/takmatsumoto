package ftn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFTNTransitionAndFilter(t *testing.T) {
	a := assert.New(t)

	// 歷史 (新到舊)
	history := FTNArray{
		*NewFTNWithStrings([]string{"04", "14", "16", "25", "35"}),
		*NewFTNWithStrings([]string{"03", "10", "16", "20", "30"}),
	}

	c1 := *NewFTNWithStrings([]string{"05", "15", "16", "20", "34"}) // repeat=1, neighbor=3, isolated=1, sandwich=1
	c2 := *NewFTNWithStrings([]string{"01", "08", "19", "28", "39"}) // repeat=0, neighbor=0, isolated=5

	tr1 := c1.Transition(history)
	a.NotNil(tr1)
	a.Equal(1, tr1.RepeatCount)
	a.Equal(3, tr1.NeighborCount)
	a.Equal(1, tr1.SandwichCount)
	a.Equal(1, tr1.Skip1RepeatCount)

	candidates := FTNArray{c1, c2}
	filtered := candidates.FilterByTransition(history, TransitionFilter{
		RepeatCounts:   []int{1},
		SandwichCounts: []int{1},
	})
	a.Len(filtered, 1)
	a.True(filtered[0].IsSame(&c1))

	minDiff := -10
	maxDiff := 10
	filteredBySum := candidates.FilterByTransition(history, TransitionFilter{
		IsolatedCounts: []int{5},
		MinSumDiff:     &minDiff,
		MaxSumDiff:     &maxDiff,
	})
	a.Len(filteredBySum, 1)
	a.True(filteredBySum[0].IsSame(&c2))
}
