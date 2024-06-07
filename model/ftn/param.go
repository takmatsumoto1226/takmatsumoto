package ftn

import "fmt"

// PickParam ...
type PickParam struct {
	Key        string
	SortType   uint
	Interval   uint
	Whichfront uint
	Spliter    uint
	Hot        uint
	Freq       int
}

// GetKey get key
func (p *PickParam) GetKey() string {
	return fmt.Sprintf("%d_%d_%d", p.SortType, p.Interval, p.Whichfront)
}
