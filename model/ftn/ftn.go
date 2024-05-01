package ftn

import (
	"fmt"
	"lottery/model/df"
	"math"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

var ballPools = map[string]BallsInfo{}

var ballIntervalStatic = map[int]int{}

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
const BallsOfFTN = 5

// Balls ...
type Balls []Ball

// Ball 球
type Ball struct {
	Number   string
	Position int // 出球的順序
	Digit    int // int的球號
	Period   int
}

func (b *Ball) Illegal() bool {
	return b.Number == "" || b.Number == "00"
}

func NewBallS(n string, pos int) *Ball {
	iB1, _ := strconv.Atoi(n)
	if strings.Contains(n, "==") {
		return &Ball{Number: n, Position: 0, Digit: 0}
	} else {
		return &Ball{Number: n, Position: pos, Digit: iB1}
	}

}

func NewBallI(n int, pos int) *Ball {
	return &Ball{Number: fmt.Sprintf("%02d", n), Position: pos, Digit: n}
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

func (bsi BallsInfo) ShowAll() {
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
	B1       Ball
	B2       Ball
	B3       Ball
	B4       Ball
	B5       Ball
	TIdx     string
	IBalls   []int
	Feature  df.Feature
}

func (fa *FTN) matchCount(n FTN) int {
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

	if len(set) == count {
		return count
	}
	return 0
}

func (fa *FTN) matchNumber(n int) bool {
	return true
}

func Ball39() []string {
	arr := []string{}
	for i := 0; i < 39; i++ {
		arr = append(arr, fmt.Sprintf("%02d", i+1))
	}
	return arr
}

func (fa *FTN) toStringArray() []string {
	return []string{fa.B1.Number, fa.B2.Number, fa.B3.Number, fa.B4.Number, fa.B5.Number}
}

func (fa *FTN) ArrInt() []int {
	return []int{fa.B1.Digit, fa.B2.Digit, fa.B3.Digit, fa.B4.Digit, fa.B5.Digit}
}

func (fa *FTN) formRow() string {
	if fa.Year == "====" {
		rowmsg := fmt.Sprintf("%s|", fa.Year)
		rowmsg = rowmsg + fmt.Sprintf("%s|", fa.MonthDay)
		for i := 1; i <= ballsCountFTN; i++ {
			rowmsg = rowmsg + "==="
		}
		fmt.Println(rowmsg)
		return rowmsg
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
		return rowmsg
	}
}

func (fa *FTN) ShowRow() {
	fmt.Println(fa.formRow())
}

// NewFTN ...
func NewFTN(arr []string) *FTN {
	if len(arr) == arrFTNCount {
		B1 := *NewBallS(arr[arrIdxB1], 0)
		B2 := *NewBallS(arr[arrIdxB2], 0)
		B3 := *NewBallS(arr[arrIdxB3], 0)
		B4 := *NewBallS(arr[arrIdxB4], 0)
		B5 := *NewBallS(arr[arrIdxB5], 0)

		arrInt := []int{B1.Digit, B2.Digit, B3.Digit, B4.Digit, B5.Digit}
		for i := 0; i < ballsCountFTN; i++ {
			if i == B1.Digit {
				B1.Period = ballIntervalStatic[i]
				ballIntervalStatic[i] = 0
			} else if i == B2.Digit {
				B2.Period = ballIntervalStatic[i]
				ballIntervalStatic[i] = 0
			} else if i == B3.Digit {
				B3.Period = ballIntervalStatic[i]
				ballIntervalStatic[i] = 0
			} else if i == B4.Digit {
				B4.Period = ballIntervalStatic[i]
				ballIntervalStatic[i] = 0
			} else if i == B5.Digit {
				B5.Period = ballIntervalStatic[i]
				ballIntervalStatic[i] = 0
			} else {
				ballIntervalStatic[i]++
			}
		}

		return &FTN{
			Year:     arr[arrIdxYear],
			MonthDay: arr[arrIdxMonthDay],
			LIdx:     arr[arrIdxLIdx],
			B1:       B1,
			B2:       B2,
			B3:       B3,
			B4:       B4,
			B5:       B5,
			TIdx:     arr[arrIdxTIdx],
			IBalls:   arrInt,
			Feature:  *df.NewFeature(arrInt, arrFTNCount),
		}
	}
	logrus.Error("FTN 資料格式錯誤")
	return Empty()
}

func NewFTNWithStrings(arr []string) *FTN {
	i1, _ := strconv.Atoi(arr[0])
	i2, _ := strconv.Atoi(arr[1])
	i3, _ := strconv.Atoi(arr[2])
	i4, _ := strconv.Atoi(arr[3])
	i5, _ := strconv.Atoi(arr[4])
	return &FTN{
		Year:     "",
		MonthDay: "",
		LIdx:     "",
		B1:       *NewBallI(i1, 0),
		B2:       *NewBallI(i2, 0),
		B3:       *NewBallI(i3, 0),
		B4:       *NewBallI(i4, 0),
		B5:       *NewBallI(i5, 0),
		TIdx:     "",
		IBalls:   []int{i1, i2, i3, i4, i5},
		Feature:  *df.NewFeature([]int{i1, i2, i3, i4, i5}, ballsCountFTN),
	}
}

func NewFTNWithInts(arr []int) *FTN {
	i1 := arr[0] + 1
	i2 := arr[1] + 1
	i3 := arr[2] + 1
	i4 := arr[3] + 1
	i5 := arr[4] + 1
	return &FTN{
		Year:     "",
		MonthDay: "",
		LIdx:     "",
		B1:       *NewBallI(i1, 0),
		B2:       *NewBallI(i2, 0),
		B3:       *NewBallI(i3, 0),
		B4:       *NewBallI(i4, 0),
		B5:       *NewBallI(i5, 0),
		TIdx:     "",
		IBalls:   []int{i1, i2, i3, i4, i5},
		Feature:  *df.NewFeature([]int{i1, i2, i3, i4, i5}, ballsCountFTN),
	}
}

func Empty() *FTN {
	return &FTN{
		"====",
		"====",
		"==",
		*NewBallS("==", 0),
		*NewBallS("==", 0),
		*NewBallS("==", 0),
		*NewBallS("==", 0),
		*NewBallS("==", 0),
		"==",
		[]int{}, *df.DefaultFeature(),
	}
}

func AdjacentNumberRecordCount() error {
	return nil
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

func (fa *FTN) AdariPrice(fb *FTN) int {
	return 0
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

func (fa *FTN) MatchFeature(t *FTN) bool {
	return fa.Feature.Compare(&t.Feature)
}

func (fa *FTN) IsSame(t *FTN) bool {
	return fa.Year == t.Year && fa.MonthDay == t.MonthDay && fa.B1.Digit == t.B1.Digit && fa.B2.Digit == t.B2.Digit && fa.B3.Digit == t.B3.Digit && fa.B4.Digit == t.B4.Digit && fa.B5.Digit == t.B5.Digit
}

func (fa *FTN) Key() string {
	return fmt.Sprintf("%s_%s_%s_%s_%s", fa.B1.Number, fa.B2.Number, fa.B3.Number, fa.B4.Number, fa.B5.Number)
}
