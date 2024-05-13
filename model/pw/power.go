package pw

import (
	"fmt"
	"lottery/model/df"
	"strconv"

	"github.com/sirupsen/logrus"
)

const ballsCountPower = 38
const ballsCountPowerS2 = 8

var numberToIndex = map[string]int{}

// PickParam ...
type PickParam struct {
	Ball       Power
	Key        string
	SortType   uint
	Interval   uint
	Whichfront uint
	Spliter    uint
	Hot        uint
}

func NewPower(arr []string) *Power {
	if len(arr) == arrPowerCount {
		i1, _ := strconv.Atoi(arr[arrIdxB1])
		i2, _ := strconv.Atoi(arr[arrIdxB2])
		i3, _ := strconv.Atoi(arr[arrIdxB3])
		i4, _ := strconv.Atoi(arr[arrIdxB4])
		i5, _ := strconv.Atoi(arr[arrIdxB5])
		i6, _ := strconv.Atoi(arr[arrIdxB6])
		return &Power{
			arr[arrIdxYear],
			arr[arrIdxMonthDay],
			arr[arrIdxLIdx],
			arr[arrIdxB1],
			arr[arrIdxB2],
			arr[arrIdxB3],
			arr[arrIdxB4],
			arr[arrIdxB5],
			arr[arrIdxB6],
			arr[arrIdxS1],
			arr[arrIdxTIdx],
			[]int{i1, i2, i3, i4, i5, i6},
			*df.NewFeature([]int{i1, i2, i3, i4, i5, i6}, ballsCountPower),
		}
	}
	logrus.Error("POWER 資料格式錯誤")
	return nil
}

func NewPowerWithString(arr []string) *Power {
	if len(arr) == 6 {
		i1, _ := strconv.Atoi(arr[0])
		i2, _ := strconv.Atoi(arr[1])
		i3, _ := strconv.Atoi(arr[2])
		i4, _ := strconv.Atoi(arr[3])
		i5, _ := strconv.Atoi(arr[4])
		i6, _ := strconv.Atoi(arr[5])
		return &Power{
			"",
			"",
			"",
			arr[0],
			arr[1],
			arr[2],
			arr[3],
			arr[4],
			arr[5],
			"",
			"",
			[]int{i1, i2, i3, i4, i5, i6},
			*df.NewFeature([]int{i1, i2, i3, i4, i5, i6}, ballsCountPower),
		}
	}
	logrus.Error("POWER 資料格式錯誤")
	return nil
}

func NewPowerWithInts(arr []int) *Power {
	if len(arr) == 6 {
		i1 := arr[0] + 1
		i2 := arr[1] + 1
		i3 := arr[2] + 1
		i4 := arr[3] + 1
		i5 := arr[4] + 1
		i6 := arr[5] + 1
		return &Power{
			"",
			"",
			"",
			fmt.Sprintf("%02d", i1),
			fmt.Sprintf("%02d", i2),
			fmt.Sprintf("%02d", i3),
			fmt.Sprintf("%02d", i4),
			fmt.Sprintf("%02d", i5),
			fmt.Sprintf("%02d", i6),
			"",
			"",
			[]int{i1, i2, i3, i4, i5, i6},
			*df.NewFeature([]int{i1, i2, i3, i4, i5, i6}, ballsCountPower),
		}
	}
	logrus.Error("POWER 資料格式錯誤")
	return nil
}

// Power ...
type Power struct {
	Year     string
	MonthDay string
	LIdx     string
	B1       string
	B2       string
	B3       string
	B4       string
	B5       string
	B6       string
	S1       string
	TIdx     string
	IBalls   []int
	Feature  df.Feature
}

func (fa *Power) toStringArray() []string {
	return []string{fa.B1, fa.B2, fa.B3, fa.B4, fa.B5, fa.B6}
}

func (b *Power) Key() string {
	return fmt.Sprintf("%s_%s_%s_%s_%s_%s", b.B1, b.B2, b.B3, b.B4, b.B5, b.B6)
}

func Empty() *Power {
	return &Power{"====", "====", "==", "==", "==", "==", "==", "==", "==", "==", "==", []int{0, 0, 0, 0, 0, 0}, *df.DefaultFeature()}
}

func (fa Power) formRow() string {
	rowmsg := fmt.Sprintf("%s|", fa.Year)
	rowmsg = rowmsg + fmt.Sprintf("%s|", fa.MonthDay)
	iB1, _ := strconv.Atoi(fa.B1)
	iB2, _ := strconv.Atoi(fa.B2)
	iB3, _ := strconv.Atoi(fa.B3)
	iB4, _ := strconv.Atoi(fa.B4)
	iB5, _ := strconv.Atoi(fa.B5)
	iB6, _ := strconv.Atoi(fa.B6)
	iS1, _ := strconv.Atoi(fa.S1)
	ballarr := []int{iB1, iB2, iB3, iB4, iB5, iB6}
	bi := 0
	for i := 1; i <= ballsCountPower; i++ {
		if ballarr[bi] == i {
			rowmsg = rowmsg + fmt.Sprintf("%02d|", ballarr[bi])
			if bi < 5 {
				bi++
			}
		} else {
			rowmsg = rowmsg + "  |"
		}
	}

	rowmsg = rowmsg + "      |"

	for i := 1; i <= ballsCountPowerS2; i++ {
		if iS1 == i {
			rowmsg = rowmsg + fmt.Sprintf("%02d|", iS1)
		} else {
			rowmsg = rowmsg + "  |"
		}
	}
	return rowmsg
}

func (fa *Power) matchCount(n Power) int {
	set := make(map[string]bool)

	for _, num := range n.toStringArray() {
		set[num] = true // setting the initial value to true
	}

	// Check elements in the second array against the set
	count := 0
	for _, num := range fa.toStringArray() {
		if set[num] {
			count++
		}
	}

	return count
}

func (fa *Power) smatch(n Power) bool {
	return fa.S1 == n.S1
}

/*
Power38 索引
*/
const (
	arrIdxYear = iota
	arrIdxMonthDay
	arrIdxLIdx
	arrIdxB1
	arrIdxB2
	arrIdxB3
	arrIdxB4
	arrIdxB5
	arrIdxB6
	arrIdxS1
	arrIdxTIdx
	arrPowerCount
)

func initNumberToIndex() {
	for i := 0; i < ballsCountPower; i++ {
		key := fmt.Sprintf("%02d", i+1)
		numberToIndex[key] = i
	}
}

func (fa *Power) MatchFeature(t *Power) bool {
	return fa.Feature.Compare(&t.Feature)
}

func (fa *Power) AdariPrice(fb *Power) int {
	mc := fa.matchCount(*fb)
	smatch := fa.smatch(*fb)
	if mc == 6 && smatch {
		return PriceTop
	} else if mc == 6 {
		return PriceSecond
	} else if mc == 5 && smatch {
		return PriceThird
	} else if mc == 5 {
		return PriceFourth
	} else if mc == 4 && smatch {
		return PriceFifth
	} else if mc == 4 {
		return PriceSixth
	} else if mc == 3 && smatch {
		return PriceSeventh
	} else if mc == 2 && smatch {
		return PriceEigth
	} else if (mc == 1 && smatch) || (mc == 3) {
		return PriceNinth
	} else {
		return 0
	}
}
