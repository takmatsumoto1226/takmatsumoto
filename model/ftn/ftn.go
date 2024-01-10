package ftn

import (
	"fmt"
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
	Position int // 由小到大排序後的位置
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
}

func Ball39() []string {
	arr := []string{}
	for i := 0; i < 30; i++ {
		arr = append(arr, fmt.Sprintf("%02d", i+1))
	}
	return arr
}

func (fa *FTN) toStringArray() []string {
	return []string{fa.B1, fa.B2, fa.B3, fa.B4, fa.B5}
}

func (fa *FTN) formRow() {
	rowmsg := fmt.Sprintf("%s|", fa.Year)
	rowmsg = rowmsg + fmt.Sprintf("%s|", fa.MonthDay)
	iB1, _ := strconv.Atoi(fa.B1)
	iB2, _ := strconv.Atoi(fa.B2)
	iB3, _ := strconv.Atoi(fa.B3)
	iB4, _ := strconv.Atoi(fa.B4)
	iB5, _ := strconv.Atoi(fa.B5)
	ballarr := []int{iB1, iB2, iB3, iB4, iB5}
	bi := 0
	for i := 1; i <= ballsCountFTN; i++ {
		if ballarr[bi] == i {
			rowmsg = rowmsg + fmt.Sprintf("%02d|", ballarr[bi])
			if bi < 4 {
				bi++
			}
		} else {
			rowmsg = rowmsg + "  |"
		}
	}
	fmt.Println(rowmsg)
}

// NewFTN ...
func NewFTN(arr []string) *FTN {
	if len(arr) == arrFTNCount {
		return &FTN{arr[arrIdxYear], arr[arrIdxMonthDay], arr[arrIdxLIdx], arr[arrIdxB1], arr[arrIdxB2], arr[arrIdxB3], arr[arrIdxB4], arr[arrIdxB5], arr[arrIdxTIdx]}
	}
	logrus.Error("FTN 資料格式錯誤")
	return Empty()
}

func Empty() *FTN {
	return &FTN{"====", "====", "==", "==", "==", "==", "==", "==", "=="}
}

func AdjacentNumberRecordCount() error {
	return nil
}

type Options struct {
	next bool
}
