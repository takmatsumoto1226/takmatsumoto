package interf

type Ball interface {
	Balls() []string
}

type Interval struct {
	Index  int
	Length int
}

type Threshold struct {
	Round      int
	SampleTime float32
	Sample     int
	Value      int
	RealSale   int32
	Interval   Interval
}
