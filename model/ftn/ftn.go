package ftn

import (
	"fmt"
	"lottery/algorithm"
	"lottery/model/common"
	"lottery/model/df"
	"math"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

var ballPools = map[string]BallsInfo{}

var ballPeriotStatic = map[int]int{}
var ballContinueStatic = map[int]int{}
var groupMapping = map[string]int{}

const ballsCountFTN = 39
const BallsOfFTN = 5

type Ball df.IBall

func (b *Ball) Illegal() bool {
	return b.Number == "" || b.Number == "00"
}

func (b *Ball) Same(cb Ball) bool {
	return b.Digit == cb.Digit
}

func (b *Ball) Disappear(n, p int) bool {
	return b.Digit == n && b.Period == p
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
	AppearBalls    BallsCount `json:"appearballs"`    // appear interval
	DisappearBalls BallsCount `json:"disappearballs"` // disappear interval
	Param          PickParam  `json:"param"`
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

func (bsi BallsCount) Presentation(dshift bool) string {
	rowmsg := "\n  "
	if !dshift {
		rowmsg = "          "
	}

	for _, bi := range bsi {
		if dshift {
			rowmsg = rowmsg + fmt.Sprintf("%3d|", bi.Ball.Digit)
		} else {
			rowmsg = rowmsg + fmt.Sprintf("%2d|", bi.Ball.Digit)
		}

	}
	rowmsg = rowmsg + "\n  "
	if !dshift {
		rowmsg = rowmsg + "        "
	}

	for _, bi := range bsi {
		if dshift {
			rowmsg = rowmsg + fmt.Sprintf("%3d|", bi.Count)
		} else {
			rowmsg = rowmsg + fmt.Sprintf("%2d|", bi.Count)
		}

	}
	return rowmsg
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
	Year        string     `json:"year"`
	MonthDay    string     `json:"monthday"`
	LIdx        string     `json:"lidx"`
	B1          Ball       `json:"b1"`
	B2          Ball       `json:"b2"`
	B3          Ball       `json:"b3"`
	B4          Ball       `json:"b4"`
	B5          Ball       `json:"b5"`
	TIdx        string     `json:"tidx"`
	IBalls      []int      `json:"iballs"`
	Feature     df.Feature `json:"feature"`
	RandomCount int        `json:"randomcount"`
}

const PriceTop = 80000000
const PriceSecond = 20000
const PriceThird = 300
const PriceFourth = 50

func (fa *FTN) matchCount(n FTN) int {
	set := make(map[string]bool)

	for _, num := range n.ToStringArr() {
		set[num] = true // setting the initial value to true
	}

	// Check elements in the second array against the set
	count := 0
	for _, num := range fa.ToStringArr() {
		if set[num] {
			count++
		}
	}

	return count
}

func (fa *FTN) numberInclude(b string) bool {

	for _, num := range fa.ToStringArr() {
		if num == b {
			return true
		}
	}

	return false
}

func (f *FTN) IsExclude(n FTN) bool {
	for _, i := range f.Feature.IBalls {
		for _, j := range n.Feature.IBalls {
			if i == j {
				return false
			}
		}
	}
	return true
}

func (f *FTN) IncludeNumbers(n FTN) []int {

	arr := []int{}
	for _, i := range f.Feature.IBalls {
		for _, j := range n.Feature.IBalls {
			if i == j {
				arr = append(arr, i)
			}
		}
	}
	return arr
}

func (f *FTN) haveNeighber(tf *FTN, c int) bool {
	count := 0
	for _, b1 := range f.IBalls {
		for _, b2 := range tf.IBalls {
			// fmt.Printf("%02d:%02d\n", b1, b2)
			if common.ABSDiffInt(b1, b2) == 1 {
				count++
				// fmt.Println(count)
			}
		}
	}

	return count == c
}

func (f *FTN) haveCol(tf *FTN, c int) bool {
	count := 0
	for _, b1 := range f.IBalls {
		for _, b2 := range tf.IBalls {
			if b1 == b2 {
				count++
			}
		}
	}

	return count == c
}

func (fa *FTN) Similar(n FTN, b byte) bool {
	return true
}

func Ball39() []string {
	arr := []string{}
	for i := 0; i < 39; i++ {
		arr = append(arr, fmt.Sprintf("%02d", i+1))
	}
	return arr
}

func (fa *FTN) ToStringArr() []string {
	return []string{fa.B1.Number, fa.B2.Number, fa.B3.Number, fa.B4.Number, fa.B5.Number}
}

func (fa *FTN) ToIntArr() []int {
	return []int{fa.B1.Digit, fa.B2.Digit, fa.B3.Digit, fa.B4.Digit, fa.B5.Digit}
}

func (fa *FTN) formRow() string {
	if fa.Year == "====" {
		rowmsg := fmt.Sprintf("%s|", fa.Year)
		rowmsg = rowmsg + fmt.Sprintf("%s|", fa.MonthDay)
		for i := 1; i <= ballsCountFTN; i++ {
			rowmsg = rowmsg + "==="
		}
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
		rowmsg = rowmsg + fmt.Sprintf("%3d|", fa.Feature.SUM)
		return rowmsg
	}
}

func (fa *FTN) simpleFormRow() string {
	if fa.Year == "====" {
		rowmsg := ""
		for i := 1; i <= ballsCountFTN; i++ {
			rowmsg = rowmsg + "==="
		}
		fmt.Println(rowmsg)
		return rowmsg
	} else {
		rowmsg := "  "
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

func (fa *FTN) SameDate(date string) bool {
	return fa.Year == date[0:4] && fa.MonthDay == date[4:8]
}

// NewFTN ...
func NewFTN(arr []string) *FTN {
	if len(arr) == arrFTNCount && arr[arrIdxB1] != "" && arr[arrIdxB1] != "00" {
		B1 := *NewBallS(arr[arrIdxB1], 0)
		B2 := *NewBallS(arr[arrIdxB2], 0)
		B3 := *NewBallS(arr[arrIdxB3], 0)
		B4 := *NewBallS(arr[arrIdxB4], 0)
		B5 := *NewBallS(arr[arrIdxB5], 0)

		arrInt := []int{B1.Digit, B2.Digit, B3.Digit, B4.Digit, B5.Digit}
		for i := 0; i < ballsCountFTN; i++ {
			if i == B1.Digit {
				B1.Period = ballPeriotStatic[i]
				ballPeriotStatic[i] = 0
				ballContinueStatic[i]++
				B1.Continue = ballContinueStatic[i]
			} else if i == B2.Digit {
				B2.Period = ballPeriotStatic[i]
				ballPeriotStatic[i] = 0
				ballContinueStatic[i]++
				B2.Continue = ballContinueStatic[i]
			} else if i == B3.Digit {
				B3.Period = ballPeriotStatic[i]
				ballPeriotStatic[i] = 0
				ballContinueStatic[i]++
				B3.Continue = ballContinueStatic[i]
			} else if i == B4.Digit {
				B4.Period = ballPeriotStatic[i]
				ballPeriotStatic[i] = 0
				ballContinueStatic[i]++
				B4.Continue = ballContinueStatic[i]
			} else if i == B5.Digit {
				B5.Period = ballPeriotStatic[i]
				ballPeriotStatic[i] = 0
				ballContinueStatic[i]++
				B5.Continue = ballContinueStatic[i]
			} else {
				ballPeriotStatic[i]++
				ballContinueStatic[i] = 0
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
	logrus.Errorf("FTN 資料格式錯誤:%v\n", arr)
	return Empty()
}

func NewFTNWithStrings(arr []string) *FTN {
	return NewFTNWithStringsAndIndex(arr, "")
}

func NewFTNWithStringsAndIndex(arr []string, idx string) *FTN {
	if len(arr) == 0 {
		return Empty()
	}
	i1, _ := strconv.Atoi(arr[0])
	i2, _ := strconv.Atoi(arr[1])
	i3, _ := strconv.Atoi(arr[2])
	i4, _ := strconv.Atoi(arr[3])
	i5, _ := strconv.Atoi(arr[4])
	return &FTN{
		Year:     "",
		MonthDay: "",
		LIdx:     idx,
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

func NewFTNWithIntsPrediction(arr []int) *FTN {
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
		[]int{},
		*df.DefaultFeature(),
		0,
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
	mc := fa.matchCount(*fb)
	if mc == 5 {
		return PriceTop
	} else if mc == 4 {
		return PriceSecond
	} else if mc == 3 {
		return PriceThird
	} else if mc == 2 {
		return PriceFourth
	} else {
		return 0
	}
}

func (fa *FTN) MatchCombinations() FTNArray {
	ftnarr := FTNArray{}
	for i := 5; i > 1; i-- {
		combinations := algorithm.Combinations(fa.ToStringArr(), i)
		ftnarr = append(ftnarr, (*NewFTNArray(combinations))...)
	}
	return ftnarr
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

func (fa *FTN) IsFullSame(t *FTN) bool {
	return fa.Year == t.Year && fa.MonthDay == t.MonthDay && fa.B1.Digit == t.B1.Digit && fa.B2.Digit == t.B2.Digit && fa.B3.Digit == t.B3.Digit && fa.B4.Digit == t.B4.Digit && fa.B5.Digit == t.B5.Digit
}

func (fa *FTN) IsSame(t *FTN) bool {
	return fa.B1.Digit == t.B1.Digit && fa.B2.Digit == t.B2.Digit && fa.B3.Digit == t.B3.Digit && fa.B4.Digit == t.B4.Digit && fa.B5.Digit == t.B5.Digit
}

func (fa *FTN) Key() string {
	return fmt.Sprintf("%s_%s_%s_%s_%s", fa.B1.Number, fa.B2.Number, fa.B3.Number, fa.B4.Number, fa.B5.Number)
}

func (fa *FTN) DateKey() string {
	return fmt.Sprintf("%s%s==%s_%s_%s_%s_%s", fa.Year, fa.MonthDay, fa.B1.Number, fa.B2.Number, fa.B3.Number, fa.B4.Number, fa.B5.Number)
}

func (fa *FTN) Date() string { return fa.Year + fa.MonthDay }

func (fa *FTN) EqualPrime(n int) bool {
	return fa.Feature.PrimeCount == n
}

func (fa *FTN) BallsPresentation() string {
	return fmt.Sprintf("%s,%s,%s,%s,%s\n", fa.B1.Number, fa.B2.Number, fa.B3.Number, fa.B4.Number, fa.B5.Number)
}
