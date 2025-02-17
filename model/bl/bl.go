package bl

import (
	"fmt"
	"lottery/model/df"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

const ballsCountBigLottery = 49

type BLBall df.IBall

// BL ...
type BL struct {
	Year     string
	MonthDay string
	LIdx     string
	B1       BLBall
	B2       BLBall
	B3       BLBall
	B4       BLBall
	B5       BLBall
	B6       BLBall
	B7       BLBall
	TIdx     string
	IBalls   []int
	Feature  df.Feature
}

func NewBallS(n string, pos int) *BLBall {
	digit, _ := strconv.Atoi(n)
	if strings.Contains(n, "==") {
		return &BLBall{
			Number:   n,
			Position: 0,
			Digit:    0,
		}
	} else {
		return &BLBall{
			Number:   n,
			Position: pos,
			Digit:    digit,
		}
	}
}

func Empty() *BL {
	return &BL{
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
		[]int{0, 0, 0, 0, 0, 0, 0},
		*df.DefaultFeature(),
	}
}

func (fa *BL) toStringArray() []string {
	return []string{
		fa.B1.Number,
		fa.B2.Number,
		fa.B3.Number,
		fa.B4.Number,
		fa.B5.Number,
		fa.B6.Number,
		fa.B7.Number,
	}
}

func (fa *BL) toStringArray2() []string {
	return []string{
		fa.B1.Number,
		fa.B2.Number,
		fa.B3.Number,
		fa.B4.Number,
		fa.B5.Number,
		fa.B6.Number,
	}
}

func (fa BL) formRow(b_optional ...bool) string {
	rowmsg := fmt.Sprintf("%s|", fa.Year)
	rowmsg = rowmsg + fmt.Sprintf("%s|", fa.MonthDay)

	ballarr := []int{
		fa.B1.Digit,
		fa.B2.Digit,
		fa.B3.Digit,
		fa.B4.Digit,
		fa.B5.Digit,
		fa.B6.Digit,
		fa.B7.Digit,
	}
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

	useB7 := true
	if len(b_optional) > 0 {
		useB7 = b_optional[0]
	}
	if useB7 {
		rowmsg = rowmsg + fmt.Sprintf("      |%02d|", fa.B7.Digit)
	}
	return rowmsg
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

func NewBL(arr []string) *BL {
	if len(arr) == arrBigLCount && arr[arrIdxB1] != "" && arr[arrIdxB1] != "00" {
		i1, _ := strconv.Atoi(arr[arrIdxB1])
		i2, _ := strconv.Atoi(arr[arrIdxB2])
		i3, _ := strconv.Atoi(arr[arrIdxB3])
		i4, _ := strconv.Atoi(arr[arrIdxB4])
		i5, _ := strconv.Atoi(arr[arrIdxB5])
		i6, _ := strconv.Atoi(arr[arrIdxB6])
		i7, _ := strconv.Atoi(arr[arrIdxB7])
		return &BL{
			arr[arrIdxYear],
			arr[arrIdxMonthDay],
			arr[arrIdxLIdx],
			*NewBallS(arr[arrIdxB1], 0),
			*NewBallS(arr[arrIdxB2], 0),
			*NewBallS(arr[arrIdxB3], 0),
			*NewBallS(arr[arrIdxB4], 0),
			*NewBallS(arr[arrIdxB5], 0),
			*NewBallS(arr[arrIdxB6], 0),
			*NewBallS(arr[arrIdxB7], 0),
			arr[arrIdxTIdx],
			[]int{i1, i2, i3, i4, i5, i6, i7},
			*df.NewFeature([]int{i1, i2, i3, i4, i5, i6}, ballsCountBigLottery),
		}
	}
	logrus.Error("NewBigLottery 資料格式錯誤")
	return nil
}

func NewPowerWithString(arr []string) *BL {
	if len(arr) == 6 {
		i1, _ := strconv.Atoi(arr[0])
		i2, _ := strconv.Atoi(arr[1])
		i3, _ := strconv.Atoi(arr[2])
		i4, _ := strconv.Atoi(arr[3])
		i5, _ := strconv.Atoi(arr[4])
		i6, _ := strconv.Atoi(arr[5])
		return &BL{
			"",
			"",
			"",
			*NewBallS(arr[0], 0),
			*NewBallS(arr[1], 0),
			*NewBallS(arr[2], 0),
			*NewBallS(arr[3], 0),
			*NewBallS(arr[4], 0),
			*NewBallS(arr[5], 0),
			*NewBallS("", 0),
			"",
			[]int{i1, i2, i3, i4, i5, i6},
			*df.NewFeature([]int{i1, i2, i3, i4, i5, i6}, ballsCountBigLottery),
		}
	}
	logrus.Error("POWER 資料格式錯誤")
	return nil
}

func NewPowerWithInts(arr []int) *BL {
	if len(arr) == 6 {
		i1 := arr[0] + 1
		i2 := arr[1] + 1
		i3 := arr[2] + 1
		i4 := arr[3] + 1
		i5 := arr[4] + 1
		i6 := arr[5] + 1
		return &BL{
			"",
			"",
			"",
			*NewBallS(fmt.Sprintf("%02d", i1), 0),
			*NewBallS(fmt.Sprintf("%02d", i2), 0),
			*NewBallS(fmt.Sprintf("%02d", i3), 0),
			*NewBallS(fmt.Sprintf("%02d", i4), 0),
			*NewBallS(fmt.Sprintf("%02d", i5), 0),
			*NewBallS(fmt.Sprintf("%02d", i6), 0),
			*NewBallS("", 0),
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

func (fa *BL) CompareFeature(t *BL) bool {
	return fa.Feature.Compare(&t.Feature)
}

func (b *BL) Key() string {
	return fmt.Sprintf("%s_%s_%s_%s_%s_%s",
		b.B1.Number,
		b.B2.Number,
		b.B3.Number,
		b.B4.Number,
		b.B5.Number,
		b.B6.Number)
}

func (fa *BL) ToStringArr() []string {
	return []string{
		fa.B1.Number,
		fa.B2.Number,
		fa.B3.Number,
		fa.B4.Number,
		fa.B5.Number,
		fa.B6.Number,
		fa.B7.Number,
	}
}

func (fa *BL) ToStringArr6() []string {
	return []string{
		fa.B1.Number,
		fa.B2.Number,
		fa.B3.Number,
		fa.B4.Number,
		fa.B5.Number,
		fa.B6.Number,
	}
}

func (fa *BL) INT6Balls() []int {
	return []int{
		fa.B1.Digit,
		fa.B2.Digit,
		fa.B3.Digit,
		fa.B4.Digit,
		fa.B5.Digit,
		fa.B6.Digit,
	}
}

func (fa *BL) INT6SBall() []int {
	return []int{
		fa.B7.Digit,
	}
}

func (fa *BL) simpleFormRow() string {
	if fa.Year == "====" {
		rowmsg := ""
		for i := 1; i <= ballsCountBigLottery; i++ {
			rowmsg = rowmsg + "==="
		}
		fmt.Println(rowmsg)
		return rowmsg
	} else {
		rowmsg := "  "
		bi := 0
		for i := 1; i <= ballsCountBigLottery; i++ {
			if fa.IBalls[bi] == i {
				rowmsg = rowmsg + fmt.Sprintf("%02d|", fa.IBalls[bi])
				if bi < 5 {
					bi++
				}
			} else {
				rowmsg = rowmsg + "  |"
			}
		}
		return rowmsg
	}
}

// 設定中獎金額 (僅供參考，實際獎金依台灣彩券公告)
var prizeTable = map[string]int{
	"6+0": 100000000, // 頭獎 (數億元)
	"5+1": 2000000,   // 貳獎 (百萬級)
	"5+0": 150000,    // 參獎 (十萬級)
	"4+1": 20000,     // 肆獎 (幾萬元)
	"4+0": 4000,      // 伍獎 (幾千元)
	"3+1": 1000,      // 陸獎 (一千元)
	"3+0": 400,       // 柒獎 (400元)
	"2+1": 400,       // 捌獎 (400元)
}

// 計算中獎結果
func (fa *BL) AdariPrice(winningNumbers, specialNumbers []int) int {
	useNumbers := fa.IBalls[:6]
	// 計算對中的普通號碼數量
	matchedRegular := fa.countMatches(winningNumbers, useNumbers)

	// 檢查是否對中特別號
	matchedSpecial := fa.contains(specialNumbers, useNumbers)

	// 建立中獎分類字串
	key := fmt.Sprintf("%d+%d", matchedRegular, matchedSpecial)

	// 查詢獎金
	if prize, exists := prizeTable[key]; exists {
		return prize
	}

	// 沒有中獎
	return 0
}

// 計算對中的號碼數量
func (fa *BL) countMatches(set1, set2 []int) int {
	matches := 0
	numberMap := make(map[int]bool)

	// 轉成 map 來加速查找
	for _, num := range set1 {
		numberMap[num] = true
	}

	// 計算對中的數量
	for _, num := range set2 {
		if numberMap[num] {
			matches++
		}
	}

	return matches
}

// 檢查是否有中特別號
func (fa *BL) contains(specialNumbers, userNumbers []int) int {
	for _, num := range userNumbers {
		for _, special := range specialNumbers {
			if num == special {
				return 1
			}
		}
	}
	return 0
}
