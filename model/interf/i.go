package interf

import "fmt"

type Ball interface {
	Balls() []string
}

type Interval struct {
	Index  int
	Length int
}

func NewInterval(index int, length int) Interval {
	return Interval{index, length}
}

func NewIntervalR(length int) Interval {
	return Interval{0, length}
}

func (i Interval) String() string {
	return fmt.Sprintf("Start : %d, Len:%d", i.Index, i.Length)
}

const (
	RangeTypeLatestDefault = iota
	RangeTypeLatestRange
	RangeTypeLatestSame
	RangeTypeSpecStartAndRangeNotes
	RangeTypeFullHistoryOnly
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
	case RangeTypeSpecStartAndRangeNotes:
		return "RangeTypeSpecStartRange"
	case RangeTypeFullHistoryOnly:
		return "RangeTypeFullHistoryOnly"
	default:
		return "RangeTypeLatestDefault"
	}
}

type Threshold struct {
	Randomer   int      `json:"randomer"`
	Round      int      `json:"round"`
	SampleTime float64  `json:"sampletime"`
	Sample     int      `json:"sample"`
	Value      int      `json:"value"`
	RealSale   int32    `json:"realsale"`
	Interval   Interval `json:"interval"`
	Smart      Smart    `json:"smart"`
	Match      bool     `json:"match"`
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
