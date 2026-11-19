package ftn

import (
	"fmt"
	"lottery/model/common"
	"lottery/model/df"
	"sort"
	"time"
)

// Transition 計算此組 539 號碼相對於 history (由新到舊) 的前後期關係特徵
func (fa *FTN) Transition(history FTNArray) *df.Transition {
	if len(fa.IBalls) == 0 || len(history) == 0 {
		return nil
	}
	curr := append([]int(nil), fa.IBalls...)
	sort.Ints(curr)

	hist := make([][]int, 0, len(history))
	for _, h := range history {
		if len(h.IBalls) == 0 {
			continue
		}
		nums := append([]int(nil), h.IBalls...)
		sort.Ints(nums)
		hist = append(hist, nums)
	}
	return df.NewTransition(curr, hist, ballsCountFTN, df.DefaultColdOmissionThreshold)
}

// TransitionFilter 539 前後期關係篩選條件；nil 或空 slice 表示該項不限制
type TransitionFilter struct {
	RepeatCounts       []int // 允許的重號數 (例 []int{0, 1})
	Skip1RepeatCounts  []int // 允許的隔期重號數
	NeighborCounts     []int // 允許的鄰號數
	IsolatedCounts     []int // 允許的孤號數
	SandwichCounts     []int // 允許的夾號數
	SplitCounts        []int // 允許的分裂號數
	SameTailPrevCounts []int // 允許的同尾延續數
	GapZoneRefills     []int // 允許的斷區回補數
	MinSumDiff         *int  // 最小和值差
	MaxSumDiff         *int  // 最大和值差
	MinRecent3Count    *int  // 最少近 3 期熱號數
	MaxRecent3Count    *int  // 最多近 3 期熱號數
	MinRecent5Count    *int  // 最少近 5 期熱號數
	MaxRecent5Count    *int  // 最多近 5 期熱號數
	ColdThreshold      int   // 冷號遺漏門檻 (<=0 使用預設 15)
}

func (f TransitionFilter) Match(tr *df.Transition) bool {
	if tr == nil {
		return false
	}
	if len(f.RepeatCounts) > 0 && !intInSlice(tr.RepeatCount, f.RepeatCounts) {
		return false
	}
	if len(f.Skip1RepeatCounts) > 0 && !intInSlice(tr.Skip1RepeatCount, f.Skip1RepeatCounts) {
		return false
	}
	if len(f.NeighborCounts) > 0 && !intInSlice(tr.NeighborCount, f.NeighborCounts) {
		return false
	}
	if len(f.IsolatedCounts) > 0 && !intInSlice(tr.IsolatedCount, f.IsolatedCounts) {
		return false
	}
	if len(f.SandwichCounts) > 0 && !intInSlice(tr.SandwichCount, f.SandwichCounts) {
		return false
	}
	if len(f.SplitCounts) > 0 && !intInSlice(tr.SplitCount, f.SplitCounts) {
		return false
	}
	if len(f.SameTailPrevCounts) > 0 && !intInSlice(tr.SameTailPrevCount, f.SameTailPrevCounts) {
		return false
	}
	if len(f.GapZoneRefills) > 0 && !intInSlice(tr.GapZoneRefill, f.GapZoneRefills) {
		return false
	}
	if f.MinSumDiff != nil && tr.SumDiff < *f.MinSumDiff {
		return false
	}
	if f.MaxSumDiff != nil && tr.SumDiff > *f.MaxSumDiff {
		return false
	}
	if f.MinRecent3Count != nil && tr.Recent3Count < *f.MinRecent3Count {
		return false
	}
	if f.MaxRecent3Count != nil && tr.Recent3Count > *f.MaxRecent3Count {
		return false
	}
	if f.MinRecent5Count != nil && tr.Recent5Count < *f.MinRecent5Count {
		return false
	}
	if f.MaxRecent5Count != nil && tr.Recent5Count > *f.MaxRecent5Count {
		return false
	}
	return true
}

// FilterByTransition 依前後期號碼關係篩選候選組合 (history 由新到舊排列)
func (fa FTNArray) FilterByTransition(history FTNArray, cond TransitionFilter) FTNArray {
	defer common.TimeTaken(time.Now(), fmt.Sprintf("FilterByTransition : %d\n", len(fa)))
	if len(history) == 0 {
		return fa
	}
	hist := make([][]int, 0, len(history))
	for _, h := range history {
		if len(h.IBalls) == 0 {
			continue
		}
		nums := append([]int(nil), h.IBalls...)
		sort.Ints(nums)
		hist = append(hist, nums)
	}
	if len(hist) == 0 {
		return fa
	}

	result := FTNArray{}
	for _, f := range fa {
		curr := append([]int(nil), f.IBalls...)
		sort.Ints(curr)
		tr := df.NewTransition(curr, hist, ballsCountFTN, cond.ColdThreshold)
		if cond.Match(tr) {
			result = append(result, f)
		}
	}
	return result
}

func intInSlice(v int, list []int) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}
