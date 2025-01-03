package bl

import (
	"fmt"
	"lottery/model/df"
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
	IBalls   []int
	Feature  df.Feature
}

func Empty() *BigLottery {
	return &BigLottery{"====", "====", "==", "==", "==", "==", "==", "==", "==", "==", "==", []int{0, 0, 0, 0, 0, 0, 0}, *df.DefaultFeature()}
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
	if len(arr) == arrBigLCount && arr[arrIdxB1] != "" && arr[arrIdxB1] != "00" {
		i1, _ := strconv.Atoi(arr[arrIdxB1])
		i2, _ := strconv.Atoi(arr[arrIdxB2])
		i3, _ := strconv.Atoi(arr[arrIdxB3])
		i4, _ := strconv.Atoi(arr[arrIdxB4])
		i5, _ := strconv.Atoi(arr[arrIdxB5])
		i6, _ := strconv.Atoi(arr[arrIdxB6])
		i7, _ := strconv.Atoi(arr[arrIdxB7])
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
			[]int{i1, i2, i3, i4, i5, i6, i7},
			*df.NewFeature([]int{i1, i2, i3, i4, i5, i6}, ballsCountBigLottery),
		}
	}
	logrus.Error("NewBigLottery 資料格式錯誤")
	return nil
}

func NewPowerWithString(arr []string) *BigLottery {
	if len(arr) == 6 {
		i1, _ := strconv.Atoi(arr[0])
		i2, _ := strconv.Atoi(arr[1])
		i3, _ := strconv.Atoi(arr[2])
		i4, _ := strconv.Atoi(arr[3])
		i5, _ := strconv.Atoi(arr[4])
		i6, _ := strconv.Atoi(arr[5])
		return &BigLottery{
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
			*df.NewFeature([]int{i1, i2, i3, i4, i5, i6}, ballsCountBigLottery),
		}
	}
	logrus.Error("POWER 資料格式錯誤")
	return nil
}

func NewPowerWithInts(arr []int) *BigLottery {
	if len(arr) == 6 {
		i1 := arr[0] + 1
		i2 := arr[1] + 1
		i3 := arr[2] + 1
		i4 := arr[3] + 1
		i5 := arr[4] + 1
		i6 := arr[5] + 1
		return &BigLottery{
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
			*df.NewFeature([]int{i1, i2, i3, i4, i5, i6}, ballsCountBigLottery),
		}
	}
	logrus.Error("POWER 資料格式錯誤")
	return nil
}

func Ball49() []int {
	arr := []int{}
	for i := 0; i < 49; i++ {
		arr = append(arr, i+1)
	}
	return arr
}

func (fa *BigLottery) CompareFeature(t *BigLottery) bool {
	return fa.Feature.Compare(&t.Feature)
}

func (b *BigLottery) Key() string {
	return fmt.Sprintf("%s_%s_%s_%s_%s_%s", b.B1, b.B2, b.B3, b.B4, b.B5, b.B6)
}

func (fa *BigLottery) ToStringArr() []string {
	return []string{fa.B1, fa.B2, fa.B3, fa.B4, fa.B5, fa.B6, fa.B7}
}
