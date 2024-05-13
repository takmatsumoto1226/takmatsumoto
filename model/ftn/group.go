package ftn

import "fmt"

type FTNGroup struct {
	GroupCount   int            `json:"group_count"`
	GroupMapping map[string]int `json:"groupmapping"`
	Statics      []int          `json:"statics"`
}

func NewFTNGroup(gc int, combinations [][]int, arr FTNArray) *FTNGroup {
	if gc == 0 || len(combinations) == 0 {
		return nil
	}

	groupMapping := map[string]int{}
	for idx, v := range combinations {
		nftn := NewFTNWithInts(v)
		groupMapping[nftn.Key()] = idx / gc
	}

	length := len(combinations)/gc + 1
	statics := make([]int, length, length)

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
	for i, _ := range fg.Statics {
		msg = msg + fmt.Sprintf("%03d|", i)

	}
	msg = msg + "\n"
	for _, v := range fg.Statics {
		msg = msg + fmt.Sprintf("%3d|", v)
	}

	return msg
}
