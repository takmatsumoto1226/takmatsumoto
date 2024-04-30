package interf

import "fmt"

type Ball interface {
	Balls() []string
}

type Interval struct {
	Index  int
	Length int
}

func (i Interval) String() string {
	return fmt.Sprintf("Start : %d, Len:%d\n", i.Index, i.Length)
}

const (
	RangeTypeRandomRange = iota
	RangeTypeLatestRange
	RangeTypeLatest
	RangeTypeLatestSame
)

type Smart struct {
	Enable bool
	Type   int
}

type Threshold struct {
	Randomer     int
	Round        int
	SampleTime   float32
	Sample       int
	Value        int
	RealSale     int32
	Interval     Interval
	Smart        Smart
	Combinations [][]int
}

func PureIntervalTH(i, l int) *Threshold {
	return &Threshold{Interval: Interval{Index: i, Length: l}}
}

func SmartPureIntervalTH(i, l int) *Threshold {
	return &Threshold{Interval: Interval{Index: i, Length: l}, Smart: Smart{Enable: true, Type: RangeTypeLatest}}
}
