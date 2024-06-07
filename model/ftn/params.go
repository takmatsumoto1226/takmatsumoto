package ftn

// PickParams ...
type PickParams []PickParam

func (pms PickParams) Add(p PickParam) {
	pms = append(pms, p)
}
