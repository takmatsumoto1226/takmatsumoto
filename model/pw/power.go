package pw

import (
	"fmt"
	"lottery/model/df"

	"github.com/sirupsen/logrus"
)

const ballsCountPower = 38
const ballsCountPowerS2 = 8

var numberToIndex = map[string]int{}

func NewPower(arr []string) *Power {
	if len(arr) == arrPowerCount {
		b1 := *NewBallS(arr[arrIdxB1], 0)
		b2 := *NewBallS(arr[arrIdxB2], 0)
		b3 := *NewBallS(arr[arrIdxB3], 0)
		b4 := *NewBallS(arr[arrIdxB4], 0)
		b5 := *NewBallS(arr[arrIdxB5], 0)
		b6 := *NewBallS(arr[arrIdxB6], 0)
		s1 := *NewBallS(arr[arrIdxS1], 0)

		return &Power{
			Year:     arr[arrIdxYear],
			MonthDay: arr[arrIdxMonthDay],
			LIdx:     arr[arrIdxLIdx],
			B1:       b1,
			B2:       b2,
			B3:       b3,
			B4:       b4,
			B5:       b5,
			B6:       b6,
			S1:       s1,
			TIdx:     arr[arrIdxTIdx],
			IBalls:   []int{b1.Digit, b2.Digit, b3.Digit, b4.Digit, b5.Digit, b6.Digit},
			Feature:  *df.NewFeature([]int{b1.Digit, b2.Digit, b3.Digit, b4.Digit, b5.Digit, b6.Digit}, ballsCountPower),
		}
	}
	logrus.Error("POWER 資料格式錯誤")
	return nil
}

func NewPowerWithString(arr []string) *Power {
	if len(arr) == 6 {
		B1 := *NewBallS(arr[0], 0)
		B2 := *NewBallS(arr[1], 0)
		B3 := *NewBallS(arr[2], 0)
		B4 := *NewBallS(arr[3], 0)
		B5 := *NewBallS(arr[4], 0)
		B6 := *NewBallS(arr[5], 0)
		S1 := *NewBallS("==", 0)
		return &Power{
			"",
			"",
			"",
			B1,
			B2,
			B3,
			B4,
			B5,
			B6,
			S1,
			"",
			[]int{B1.Digit, B2.Digit, B3.Digit, B4.Digit, B5.Digit, B6.Digit},
			*df.NewFeature([]int{B1.Digit, B2.Digit, B3.Digit, B4.Digit, B5.Digit, B6.Digit}, ballsCountPower),
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
			*NewBallI(i1, 0),
			*NewBallI(i2, 0),
			*NewBallI(i3, 0),
			*NewBallI(i4, 0),
			*NewBallI(i5, 0),
			*NewBallI(i6, 0),
			*NewBallI(0, 0),
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
	B1       Ball
	B2       Ball
	B3       Ball
	B4       Ball
	B5       Ball
	B6       Ball
	S1       Ball
	TIdx     string
	IBalls   []int
	Feature  df.Feature
}

func (pw *Power) toStringArray() []string {
	return []string{pw.B1.Number, pw.B2.Number, pw.B3.Number, pw.B4.Number, pw.B5.Number, pw.B6.Number}
}

func (pw *Power) Key() string {
	return fmt.Sprintf("%s_%s_%s_%s_%s_%s", pw.B1.Number, pw.B2.Number, pw.B3.Number, pw.B4.Number, pw.B5.Number, pw.B6.Number)
}

func Empty() *Power {
	return &Power{
		"====",
		"====",
		"==",
		*NewBallS("==", 0),
		*NewBallS("==", 0),
		*NewBallS("==", 0),
		*NewBallS("==", 0),
		*NewBallS("==", 0),
		*NewBallS("==", 0),
		*NewBallS("==", 0),
		"==",
		[]int{0, 0, 0, 0, 0, 0},
		*df.DefaultFeature(),
	}
}

func (pw *Power) formRow() string {
	rowmsg := fmt.Sprintf("%s|", pw.Year)
	rowmsg = rowmsg + fmt.Sprintf("%s|", pw.MonthDay)

	ballarr := []int{
		pw.B1.Digit,
		pw.B2.Digit,
		pw.B3.Digit,
		pw.B4.Digit,
		pw.B5.Digit,
		pw.B6.Digit,
		pw.S1.Digit,
	}
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
		if pw.S1.Digit == i {
			rowmsg = rowmsg + pw.S1.Number
		} else {
			rowmsg = rowmsg + "  |"
		}
	}
	return rowmsg
}

func (pw *Power) matchCount(n Power) int {
	set := make(map[string]bool)

	for _, num := range n.toStringArray() {
		set[num] = true // setting the initial value to true
	}

	// Check elements in the second array against the set
	count := 0
	for _, num := range pw.toStringArray() {
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

func (pw *Power) MatchFeature(t *Power) bool {
	return pw.Feature.Compare(&t.Feature)
}

func (pw *Power) AdariPrice(fb *Power) int {
	mc := pw.matchCount(*fb)
	smatch := pw.smatch(*fb)
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

func (pw *Power) ShowRow() {
	fmt.Println(pw.formRow())
}

func (fa *Power) IsFullSame(t *Power) bool {
	return fa.Year == t.Year && fa.MonthDay == t.MonthDay && fa.B1.Digit == t.B1.Digit && fa.B2.Digit == t.B2.Digit && fa.B3.Digit == t.B3.Digit && fa.B4.Digit == t.B4.Digit && fa.B5.Digit == t.B5.Digit && fa.B6.Digit == t.B6.Digit
}

func (fa *Power) IsSame(t *Power) bool {
	return fa.B1.Digit == t.B1.Digit && fa.B2.Digit == t.B2.Digit && fa.B3.Digit == t.B3.Digit && fa.B4.Digit == t.B4.Digit && fa.B5.Digit == t.B5.Digit && fa.B6.Digit == t.B6.Digit
}
