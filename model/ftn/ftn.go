package ftn

import (
	"fmt"
	"math"
	"strconv"

	"github.com/sirupsen/logrus"
)

var ballPools = map[string]BallsInfo{}

// PickParams ...
type PickParams []PickParam

// PickParam ...
type PickParam struct {
	Key        string
	SortType   uint
	Interval   uint
	Whichfront uint
	Spliter    uint
	Hot        uint // 熱門號碼
}

// GetKey 取得key
func (p *PickParam) GetKey() string {
	return fmt.Sprintf("%d_%d_%d", p.SortType, p.Interval, p.Whichfront)
}

const ballsCountFTN = 39

// Balls ...
type Balls []Ball

// Ball 球
type Ball struct {
	Number   string
	Position int // 出球的順序
}

// BallsCount ...
type BallsCount []BallInfo
type NormalizeInfo struct {
	NorBalls BallsCount // 統計號碼多久沒出現
	Param    PickParam
}

// Len ...
func (fa BallsCount) Len() int {
	return len(fa)
}

// Less ...
func (fa BallsCount) Less(i, j int) bool {
	return fa[i].Count < fa[j].Count
}

// Swap swaps the elements with indexes i and j.
func (fa BallsCount) Swap(i, j int) {
	fa[i], fa[j] = fa[j], fa[i]
}

// BallCount ...
type BallInfo struct {
	Count    uint
	Intervel uint
	Ball     Ball
}

type BallsInfo []BallInfo

func (bsi BallsInfo) Presentation() {
	rowmsg := "          "
	for _, bi := range bsi {
		rowmsg = rowmsg + fmt.Sprintf("%2d ", bi.Count)
	}
	fmt.Println(rowmsg)
}

/*
FTN
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
	arrIdxTIdx
	arrFTNCount
)

// Len ...

var numberToIndex = map[string]int{}

func initNumberToIndex() {
	for i := 0; i < ballsCountFTN; i++ {
		key := fmt.Sprintf("%02d", i+1)
		numberToIndex[key] = i
	}
}

// FTN ...
type FTN struct {
	Year     string
	MonthDay string
	LIdx     string
	B1       string
	B2       string
	B3       string
	B4       string
	B5       string
	TIdx     string
	IBalls   []int
}

func Ball39() []string {
	arr := []string{}
	for i := 0; i < 39; i++ {
		arr = append(arr, fmt.Sprintf("%02d", i+1))
	}
	return arr
}

func (fa *FTN) toStringArray() []string {
	return []string{fa.B1, fa.B2, fa.B3, fa.B4, fa.B5}
}

func (fa *FTN) ArrInt() []int {
	iB1, _ := strconv.Atoi(fa.B1)
	iB2, _ := strconv.Atoi(fa.B2)
	iB3, _ := strconv.Atoi(fa.B3)
	iB4, _ := strconv.Atoi(fa.B4)
	iB5, _ := strconv.Atoi(fa.B5)
	return []int{iB1, iB2, iB3, iB4, iB5}
}

func (fa *FTN) formRow() {
	if fa.Year == "====" {
		rowmsg := fmt.Sprintf("%s|", fa.Year)
		rowmsg = rowmsg + fmt.Sprintf("%s|", fa.MonthDay)
		for i := 1; i <= ballsCountFTN; i++ {
			rowmsg = rowmsg + "==="
		}
		fmt.Println(rowmsg)
	} else {
		rowmsg := fmt.Sprintf("%s|", fa.Year)
		rowmsg = rowmsg + fmt.Sprintf("%s|", fa.MonthDay)
		bi := 0
		for i := 1; i <= ballsCountFTN; i++ {
			if fa.IBalls[bi] == i {
				rowmsg = rowmsg + fmt.Sprintf("%02d|", fa.IBalls[bi])
				if bi < 4 {
					bi++
				}
			} else {
				rowmsg = rowmsg + "  |"
			}
		}
		fmt.Println(rowmsg)
	}

}

// NewFTN ...
func NewFTN(arr []string) *FTN {
	if len(arr) == arrFTNCount {
		i1, _ := strconv.Atoi(arr[arrIdxB1])
		i2, _ := strconv.Atoi(arr[arrIdxB2])
		i3, _ := strconv.Atoi(arr[arrIdxB3])
		i4, _ := strconv.Atoi(arr[arrIdxB4])
		i5, _ := strconv.Atoi(arr[arrIdxB5])
		return &FTN{arr[arrIdxYear], arr[arrIdxMonthDay], arr[arrIdxLIdx], arr[arrIdxB1], arr[arrIdxB2], arr[arrIdxB3], arr[arrIdxB4], arr[arrIdxB5], arr[arrIdxTIdx], []int{i1, i2, i3, i4, i5}}
	}
	logrus.Error("FTN 資料格式錯誤")
	return Empty()
}

func Empty() *FTN {
	return &FTN{"====", "====", "==", "==", "==", "==", "==", "==", "==", []int{}}
}

func AdjacentNumberRecordCount() error {
	return nil
}

func (fa *FTN) IsContinue2() bool {
	i1 := fa.IBalls[0]
	i2 := fa.IBalls[1]
	i3 := fa.IBalls[2]
	i4 := fa.IBalls[3]
	i5 := fa.IBalls[4]
	return i2-i1 == 1 || i3-i2 == 1 || i4-i3 == 1 || i5-i4 == 1
}
func (fa *FTN) IsContinue3() bool {
	i1 := fa.IBalls[0]
	i2 := fa.IBalls[1]
	i3 := fa.IBalls[2]
	i4 := fa.IBalls[3]
	i5 := fa.IBalls[4]
	return (i2-i1 == 1 && i3-i2 == 1) || (i3-i2 == 1 && i4-i3 == 1) || (i4-i3 == 1 && i5-i4 == 1)
}

func (fa *FTN) IsContinue4() bool {
	i1 := fa.IBalls[0]
	i2 := fa.IBalls[1]
	i3 := fa.IBalls[2]
	i4 := fa.IBalls[3]
	i5 := fa.IBalls[4]
	return (i2-i1 == 1 && i3-i2 == 1 && i4-i3 == 1) || (i3-i2 == 1 && i4-i3 == 1 && i5-i4 == 1)
}

func (fa *FTN) IsContinue5() bool {
	i1 := fa.IBalls[0]
	i2 := fa.IBalls[1]
	i3 := fa.IBalls[2]
	i4 := fa.IBalls[3]
	i5 := fa.IBalls[4]
	return i2-i1 == 1 && i3-i2 == 1 && i4-i3 == 1 && i5-i4 == 1
}

func (fa *FTN) IsContinue22() bool {
	i1 := fa.IBalls[0]
	i2 := fa.IBalls[1]
	i3 := fa.IBalls[2]
	i4 := fa.IBalls[3]
	i5 := fa.IBalls[4]

	count := 0
	if i2-i1 == 1 {
		count++
	}

	if i3-i2 == 1 {
		count++
	}

	if i4-i3 == 1 {
		count++
	}

	if i5-i4 == 1 {
		count++
	}

	return count == 2 && !fa.IsContinue3()
}

type Options struct {
	next bool
}

/* 本期號碼有出現34, 上一期或下一期出現35 36*/
func (fa *FTN) IsDTree(next *FTN) bool {
	for i := 0; i < 5; i++ {
		for j := 0; j < 4; j++ {
			if math.Abs(float64(fa.IBalls[i]-next.IBalls[j])) == 1 && math.Abs(float64(fa.IBalls[i]-next.IBalls[j+1])) == 1 {
				return true
			}
		}
	}
	return false
}

func (fa *FTN) IsUTree(before *FTN) bool {
	for i := 0; i < 4; i++ {
		for j := 0; j < 5; j++ {
			if math.Abs(float64(fa.IBalls[j]-before.IBalls[i])) == 1 && math.Abs(float64(fa.IBalls[j]-before.IBalls[i+1])) == 1 {
				return true
			}
		}
	}
	return false
}
