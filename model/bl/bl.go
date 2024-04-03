package bl

import (
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
)

const ballsCountBigLottery = 49

// BigLottery ...
type BigLottery struct {
	Year     string
	MonthDay string
	LIdx     string
	B1       string
	B2       string
	B3       string
	B4       string
	B5       string
	B6       string
	B7       string
	TIdx     string
}

func Empty() *BigLottery {
	return &BigLottery{"====", "====", "==", "==", "==", "==", "==", "==", "==", "==", "=="}
}

func (fa *BigLottery) toStringArray() []string {
	return []string{fa.B1, fa.B2, fa.B3, fa.B4, fa.B5, fa.B6, fa.B7}
}

func (fa *BigLottery) toStringArray2() []string {
	return []string{fa.B1, fa.B2, fa.B3, fa.B4, fa.B5, fa.B6}
}

func (fa BigLottery) formRow() {
	rowmsg := fmt.Sprintf("%s|", fa.Year)
	rowmsg = rowmsg + fmt.Sprintf("%s|", fa.MonthDay)
	iB1, _ := strconv.Atoi(fa.B1)
	iB2, _ := strconv.Atoi(fa.B2)
	iB3, _ := strconv.Atoi(fa.B3)
	iB4, _ := strconv.Atoi(fa.B4)
	iB5, _ := strconv.Atoi(fa.B5)
	iB6, _ := strconv.Atoi(fa.B6)
	iB7, _ := strconv.Atoi(fa.B7)
	ballarr := []int{iB1, iB2, iB3, iB4, iB5, iB6}
	bi := 0
	for i := 1; i <= ballsCountBigLottery; i++ {
		if ballarr[bi] == i {
			rowmsg = rowmsg + fmt.Sprintf("%02d|", ballarr[bi])
			if bi < 5 {
				bi++
			}
		} else {
			rowmsg = rowmsg + "  |"
		}
	}

	rowmsg = rowmsg + fmt.Sprintf("      |%02d|", iB7)
	fmt.Println(rowmsg)
}

// Ball 球
type Ball struct {
	Number string
}

// PickParam ...
type PickParam struct {
	Key        string
	SortType   uint
	Interval   uint
	Whichfront uint
	Spliter    uint
	Hot        uint
}

type BallsCount []BallInfo

// BallCount ...
type BallInfo struct {
	Count uint
	Ball  Ball
}

type NormalizeInfo struct {
	NorBalls BallsCount
	Param    PickParam
}

/*
BigLottery
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
	arrIdxB7
	arrIdxTIdx
	arrBigLCount
)

func NewBigLottery(arr []string) *BigLottery {
	if len(arr) == arrBigLCount {
		return &BigLottery{
			arr[arrIdxYear],
			arr[arrIdxMonthDay],
			arr[arrIdxLIdx],
			arr[arrIdxB1],
			arr[arrIdxB2],
			arr[arrIdxB3],
			arr[arrIdxB4],
			arr[arrIdxB5],
			arr[arrIdxB6],
			arr[arrIdxB7],
			arr[arrIdxTIdx],
		}
	}
	logrus.Error("NewBigLottery 資料格式錯誤")
	return nil
}

func Ball49() []int {
	arr := []int{}
	for i := 0; i < 49; i++ {
		arr = append(arr, i+1)
	}
	return arr
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
