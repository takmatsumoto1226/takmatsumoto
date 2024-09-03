package ftn

import (
	"fmt"
	"sort"
)

type TenGroupStatics []TenGroupStatic

func (fa TenGroupStatics) Len() int {
	return len(fa)
}

// Less ...
func (tgs TenGroupStatics) Less(i, j int) bool {
	return tgs[i].Percent > tgs[j].Percent
}

// Swap swaps the elements with indexes i and j.
func (tgs TenGroupStatics) Swap(i, j int) {
	tgs[i], tgs[j] = tgs[j], tgs[i]
}

type TenGroup struct {
	ID       string
	Features []int
}

type TenGroupStatic struct {
	TenGroup
	FTNs    FTNArray
	Percent float64
}

type TenGroupMgr struct {
	Arr       TenGroupStatics
	TenGroups map[string]TenGroupStatic
}

func NewTenGroupStatic(gs []int) TenGroupStatic {
	return TenGroupStatic{
		TenGroup: TenGroup{
			ID:       fmt.Sprintf("%v", gs),
			Features: gs,
		},
		FTNs:    FTNArray{},
		Percent: 0.,
	}
}

func (tg *TenGroupStatic) Presentation() string {
	return tg.ID + fmt.Sprintf("  %.2f", tg.Percent)
}

func (tg *TenGroupStatic) Add(f *FTN) {
	tg.FTNs = append(tg.FTNs, *f)
}

func NewTenGroupMgr(fa FTNArray) TenGroupMgr {
	n := 5
	k := 4
	result := [][]int{}
	findSolutions(k, n, []int{}, &result)
	tgs := map[string]TenGroupStatic{}
	for _, s := range result {

		ID := fmt.Sprintf("%v", s)
		tgs[ID] = NewTenGroupStatic(s)
	}

	for _, f := range fa {
		ID := fmt.Sprintf("%v", f.Feature.TenGroupCount[:4])
		cf := f
		tg := tgs[ID]
		tg.FTNs = append(tg.FTNs, cf)
		tgs[ID] = tg
	}

	arr := TenGroupStatics{}
	for key, tg := range tgs {
		tg.Percent = (float64(len(tg.FTNs)) / float64(len(fa)) * 100)
		tgs[key] = tg
		atg := tg
		arr = append(arr, atg)
	}

	return TenGroupMgr{
		Arr:       arr,
		TenGroups: tgs,
	}
}

func (mgr *TenGroupMgr) Presentation() string {
	msg := ""
	sort.Sort(mgr.Arr)
	for _, tg := range mgr.Arr {
		msg = msg + tg.Presentation() + "\n"
	}
	return msg
}

// findSolutions recursively finds and prints all solutions
func findSolutions(k, n int, prefix []int, result *[][]int) {
	if k == 1 {
		prefix = append(prefix, n)
		*result = append(*result, prefix)
		return
	}

	for i := 0; i <= n; i++ {
		findSolutions(k-1, n-i, append([]int(nil), append(prefix, i)...), result)
	}
}
