package pw

import (
	"fmt"
	"lottery/model/common"
)

type PWGroup struct {
	GroupCount   int            `json:"group_count"`
	GroupMapping map[string]int `json:"groupmapping"`
	Statics      []int          `json:"statics"`
	ZeroCount    int            `json:"zero_count"`
	Max          int            `json:"max"`
	Avg          float64        `json:"avg"`
	// next         Filter
}

func NewPWGroup(gc int, combinations [][]int, arr PowerList) *PWGroup {
	if gc == 0 || len(combinations) == 0 {
		return nil
	}

	groupMapping := map[string]int{}
	for idx, v := range combinations {
		nftn := NewPowerWithInts(v)
		groupMapping[nftn.Key()] = idx / gc
	}

	length := len(combinations)/gc + 1
	statics := make([]int, length)

	for _, pw := range arr {
		idx := groupMapping[pw.Key()]
		statics[idx]++
	}

	return &PWGroup{GroupCount: gc, GroupMapping: groupMapping, Statics: statics}
}

func (pg *PWGroup) FindGroupIndex(ftn Power) (int, int) {
	gi := pg.GroupMapping[ftn.Key()]
	return gi, pg.Statics[gi]
}

func (pg *PWGroup) Presentation() string {
	msg := ""
	sum := 0.0
	for i, _ := range pg.Statics {
		msg = msg + fmt.Sprintf("%05d|", i)

	}
	msg = msg + "\n"
	for _, v := range pg.Statics {
		if v == 0 {
			pg.ZeroCount++
		}
		msg = msg + fmt.Sprintf("%5d|", v)
		pg.Max = common.MAX(pg.Max, v)
		sum = sum + float64(v)
	}
	pg.Avg = sum / float64(len(pg.Statics))

	total := (len(pg.GroupMapping) / pg.GroupCount) + 1
	msg = msg + "\n"
	msg = msg + fmt.Sprintf("Max : %d\n", pg.Max)
	msg = msg + fmt.Sprintf("Avg : %f\n", pg.Avg)
	msg = msg + fmt.Sprintf("Total : %d\n", total)
	msg = msg + fmt.Sprintf("Zero Count : %d\n", pg.ZeroCount)
	msg = msg + fmt.Sprintf("Zero Percent : %.2f\n", float64(pg.ZeroCount)/float64(total))
	return msg
}

func (pg *PWGroup) filter(bt *Power) {

}

// func (fg *FTNGroup) setNext(next Filter) {
// 	fg.next = next
// }
