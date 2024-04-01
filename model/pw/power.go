package pw

import (
	"fmt"
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
	Hot        uint // 熱門號碼
}

func NewPower(arr []string) *Power {
	if len(arr) == arrPowerCount {
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
}

func (fa *Power) toStringArray() []string {
	return []string{fa.B1, fa.B2, fa.B3, fa.B4, fa.B5, fa.B6}
}

func Empty() *Power {
	return &Power{"====", "====", "==", "==", "==", "==", "==", "==", "==", "==", "=="}
}

func (fa Power) formRow() {
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
	fmt.Println(rowmsg)
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

type IBalls struct {
	B1 int
	B2 int
	B3 int
	B4 int
	B5 int
	B6 int
	S1 int
}

func (b *IBalls) Key() string {
	return fmt.Sprintf("%02d_%02d_%02d_%02d_%02d_%02d", b.B1, b.B2, b.B3, b.B4, b.B5, b.B6)
}

func NewBalls(n []int) *IBalls {
	return &IBalls{B1: n[0], B2: n[1], B3: n[2], B4: n[3], B5: n[4], B6: n[5]}
}
