package ftn

import (
	"fmt"
	"lottery/model/common"
)

type FTNGroup struct {
	GroupCount   int            `json:"group_count"`
	GroupMapping map[string]int `json:"groupmapping"`
	Statics      []int          `json:"statics"`
	ZeroCount    int            `json:"zero_count"`
	Max          int            `json:"max"`
	Avg          float64        `json:"avg"`
	TargetIndex  int            `json:"target_index"`
	// next         Filter
}

func NewGroup(gc int, combinations [][]int, arr FTNArray) *FTNGroup {
	if gc == 0 || len(combinations) == 0 {
		return nil
	}

	groupMapping := map[string]int{}
	for idx, v := range combinations {
		nftn := NewFTNWithInts(v)
		groupMapping[nftn.Key()] = idx / gc
	}

	length := len(combinations)/gc + 1
	statics := make([]int, length)

	for _, ftn := range arr {
		idx := groupMapping[ftn.Key()]
		statics[idx]++
	}

	return &FTNGroup{GroupCount: gc, GroupMapping: groupMapping, Statics: statics}
}

func (fg *FTNGroup) FindGroupIndex(ftn FTN) (int, int) {
	gi := fg.GroupMapping[ftn.Key()]
	return gi, fg.Statics[gi]
}

func (fg *FTNGroup) Presentation() string {
	msg := ""
	sum := 0.0
	for i, _ := range fg.Statics {
		msg = msg + fmt.Sprintf("%05d|", i)

	}
	msg = msg + "\n"
	for _, v := range fg.Statics {
		if v == 0 {
			fg.ZeroCount++
		}
		msg = msg + fmt.Sprintf("%5d|", v)
		fg.Max = common.MAX(fg.Max, v)
		sum = sum + float64(v)
	}
	fg.Avg = sum / float64(len(fg.Statics))

	msg = msg + "\n"
	msg = msg + fmt.Sprintf("Max : %d\n", fg.Max)
	msg = msg + fmt.Sprintf("Avg : %f\n", fg.Avg)
	msg = msg + fmt.Sprintf("Zero Count : %d\n", fg.ZeroCount)
	msg = msg + fmt.Sprintf("Zero Percent : %.2f\n", float64(fg.ZeroCount)/float64(len(fg.Statics)))
	return msg
}

// func (fg *FTNGroup) setNext(next Filter) {
// 	fg.next = next
// }
