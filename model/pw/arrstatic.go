package pw

import (
	"errors"
	"fmt"
	"lottery/model/interf"

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

func (ar PowerList) StaticContinue2Percent(i interf.Interval) float64 {

	sr := ar
	if i.Length > 0 && i.Index+i.Length < len(ar) {
		sr = ar[i.Index : i.Index+i.Length]
	}
	count := 0
	for _, f := range sr {
		if f.Feature.IsContinue2() {
			f.ShowRow()
			count++
		}
	}
	return (float64(count) / float64(len(sr))) * 100
}

func (ar PowerList) StaticContinue3Percent(i interf.Interval) float64 {

	sr := ar
	if i.Length > 0 && i.Index+i.Length < len(ar) {
		sr = ar[i.Index : i.Index+i.Length]
	}
	count := 0
	for _, f := range sr {
		if f.Feature.IsContinue3() {
			f.ShowRow()
			count++
		}
	}
	return (float64(count) / float64(len(sr))) * 100
}

func (ar PowerList) StaticContinue22Percent(i interf.Interval) float64 {

	sr := ar
	if i.Length > 0 && i.Index+i.Length < len(ar) {
		sr = ar[i.Index : i.Index+i.Length]
	}
	count := 0
	for _, f := range sr {
		if f.Feature.IsContinue22() {
			f.ShowRow()
			count++
		}
	}
	return (float64(count) / float64(len(sr))) * 100
}

func (ar PowerList) StaticContinue4Percent(i interf.Interval) float64 {

	sr := ar
	if i.Length > 0 && i.Index+i.Length < len(ar) {
		sr = ar[i.Index : i.Index+i.Length]
	}
	count := 0
	for i, f := range sr {
		if f.Feature.IsContinue4() {
			sr[i+1].ShowRow()
			f.ShowRow()
			fmt.Println("")
			count++
		}
	}
	return (float64(count) / float64(len(sr))) * 100
}
