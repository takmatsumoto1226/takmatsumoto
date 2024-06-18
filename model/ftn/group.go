package ftn

import (
	"fmt"
	"lottery/model/common"
	"strings"
)

type FTNGroup struct {
	GroupCount   int            `json:"group_count"`
	GroupMapping map[string]int `json:"groupmapping"`
	CountMapping map[int]int    `json:"groupmapping"`
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

	length := (len(combinations) / gc) + 1
	statics := make([]int, length)

	for _, ftn := range arr {
		idx := groupMapping[ftn.Key()]
		statics[idx]++
	}

	return &FTNGroup{GroupCount: gc, GroupMapping: groupMapping, Statics: statics, CountMapping: map[int]int{}}
}

func (fg *FTNGroup) FindGroupIndex(ftn FTN) (int, int) {
	gi := fg.GroupMapping[ftn.Key()]
	return gi, fg.Statics[gi]
}

func (fg *FTNGroup) StaticCounts() {
	for _, v := range fg.Statics {
		if i, ok := fg.CountMapping[v]; ok {
			fg.CountMapping[v] = i + 1
		} else {
			fg.CountMapping[v] = 1
		}
	}
}

func (fg *FTNGroup) FindGroupNotes(c int) FTNArray {
	result := FTNArray{}
	gi := -1
	for i, v := range fg.Statics {
		if v == c {
			gi = i
			break
		}
	}

	for k, v := range fg.GroupMapping {
		if v == gi {
			thNumb := NewFTNWithStrings(strings.Split(k, "_"))
			result = append(result, *thNumb)
		}
	}
	return result
}

func (fg *FTNGroup) Presentation() string {
	msg := ""
	sum := 0.0
	for i, _ := range fg.Statics {
		msg = msg + fmt.Sprintf("%06d|", i)

	}
	msg = msg + "\n"
	for _, v := range fg.Statics {
		if v == 0 {
			fg.ZeroCount++
		}
		msg = msg + fmt.Sprintf("%6d|", v)
		fg.Max = common.MAX(fg.Max, v)
		sum = sum + float64(v)
	}
	fg.Avg = sum / float64(len(fg.Statics))

	msg = msg + "\n"
	msg = msg + fmt.Sprintf("Max : %d\n", fg.Max)
	msg = msg + fmt.Sprintf("Avg : %f\n", fg.Avg)
	msg = msg + fmt.Sprintf("Zero Count : %d\n", fg.ZeroCount)
	msg = msg + fmt.Sprintf("Zero Percent : %.2f\n", float64(fg.ZeroCount)/float64(len(fg.Statics)))
	for idx, i := range fg.CountMapping {
		msg = msg + fmt.Sprintf("Count Percent %02d: %.2f%%\n", idx, (float64(i)/float64(len(fg.Statics)))*100)
	}
	return msg
}

// func (fg *FTNGroup) setNext(next Filter) {
// 	fg.next = next
// }
