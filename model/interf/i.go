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
	return fmt.Sprintf("Start : %d, Len:%d", i.Index, i.Length)
}

const (
	RangeTypeLatestDefault = iota
	RangeTypeLatestRange
	RangeTypeLatestSame
	RangeTypeSpecStartRange
)

type Smart struct {
	Enable bool
	Type   int
}

func (s *Smart) String() string {
	return fmt.Sprintf("%t:%s", s.Enable, s.typeName())
}

func (s *Smart) typeName() string {
	switch s.Type {
	// case RangeTypeRandomRange:
	// 	return "RangeTypeRandomRange"
	case RangeTypeLatestRange:
		return "RangeTypeLatest"
	case RangeTypeLatestSame:
		return "RangeTypeLatestSame"
	case RangeTypeSpecStartRange:
		return "RangeTypeSpecStartRange"
	default:
		return "RangeTypeLatestDefault"
	}
}

type Threshold struct {
	Randomer   int      `json:"randomer"`
	Round      int      `json:"round"`
	SampleTime float32  `json:"sampletime"`
	Sample     int      `json:"sample"`
	Value      int      `json:"value"`
	RealSale   int32    `json:"realsale"`
	Interval   Interval `json:"interval"`
	Smart      Smart    `json:"smart"`
}

func (th *Threshold) GetRandomer() string {
	switch th.Randomer {
	case 1:
		return "mt19937"
	default:
		return "defaultRand"
	}
}

func (th *Threshold) Presentation() string {
	return fmt.Sprintf("Randormer : %v\nRound:%d\nSampleTime:%.2f\nSample:%d\nValue:%d\nInterval === %s\nSmart === %s\n",
		th.GetRandomer(),
		th.Round,
		th.SampleTime,
		th.Sample,
		th.Value,
		th.Interval.String(),
		th.Smart.String())

}

func (th *Threshold) ShowAll() {
	fmt.Println(th.Presentation())
}

func PureIntervalTH(i, l int) *Threshold {
	return &Threshold{Interval: Interval{Index: i, Length: l}}
}

func SmartPureIntervalTH(i, l int) *Threshold {
	return &Threshold{Interval: Interval{Index: i, Length: l}, Smart: Smart{Enable: true, Type: RangeTypeLatestRange}}
}
