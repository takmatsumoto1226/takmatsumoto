package ftn

import "fmt"

type TenGroup struct {
	ID            string
	TenGroupCount []int
	FTNs          FTNArray
	Percent       float64
}

type TenGroupMgr struct {
	List      FTNArray
	TenGroups []TenGroup
}

func NewTenGroup(gs []int) TenGroup {
	return TenGroup{
		ID:            fmt.Sprintf("%v", gs),
		TenGroupCount: gs,
		FTNs:          FTNArray{},
		Percent:       0.,
	}
}

func NewTenGroupMgr(fa FTNArray) TenGroupMgr {
	return TenGroupMgr{
		List:      fa,
		TenGroups: []TenGroup{},
	}
}
