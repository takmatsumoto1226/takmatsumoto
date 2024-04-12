package interf

type Ball interface {
	Balls() []string
}

type Threshold struct {
	Round      int
	SampleTime float32
	Sample     int
	Value      int
	RealSale   int32
}
