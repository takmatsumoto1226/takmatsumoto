package pw

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
)

func (ar PowerList) IntervalBallsCountStatic(p PickParam) Balls {

	if p.Interval == 0 {
		logrus.Error(errors.New("不可指定0"))
		return Balls{}
	}
	var FTNIntervalCount = [ballsCountPower]uint{}
	var disappearCount = [ballsCountPower]uint{}

	for _, t := range ar {
		FTNIntervalCount[numberToIndex[t.B1.Number]]++
		FTNIntervalCount[numberToIndex[t.B2.Number]]++
		FTNIntervalCount[numberToIndex[t.B3.Number]]++
		FTNIntervalCount[numberToIndex[t.B4.Number]]++
		FTNIntervalCount[numberToIndex[t.B5.Number]]++
		FTNIntervalCount[numberToIndex[t.B6.Number]]++
		for i := 0; i < ballsCountPower; i++ {
			if i != numberToIndex[t.B1.Number] ||
				i != numberToIndex[t.B2.Number] ||
				i != numberToIndex[t.B3.Number] ||
				i != numberToIndex[t.B4.Number] ||
				i != numberToIndex[t.B5.Number] ||
				i != numberToIndex[t.B6.Number] {
				disappearCount[i]++
			} else {
				disappearCount[i] = 0
			}
		}
	}

	arr := Balls{}
	for i, count := range FTNIntervalCount {
		b := Ball{
			Number:   fmt.Sprintf("%02d", i+1),
			Position: 0,
			Digit:    i + 1,
			Period:   0,
			Continue: 0,
			Count:    count,
		}
		arr = append(arr, b)
	}

	return arr
}
