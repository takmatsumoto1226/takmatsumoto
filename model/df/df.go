package df

import (
	"bytes"
	"fmt"
)

const (
	InfoFTN = iota
	Info49
	InfoPOWER
	Info4STAR
)

const (
	HOT = iota
	COOL
)

const (
	REVERSE = iota
)

const (
	Descending = iota // raw data 年份\
	Ascending         // raw data 年份小到大
)

const (
	Biggerfront = iota // 球數出現次數統計後, 出現次數多得在前面
	Smallfront         // 球數出現次數統計後, 出現次數少的在前面
	Normal             // ball的數字由小到大排序
)

// 特徵種類
const (
	ContinueColNone = iota
	ContinueCol1    // 跟前一期號碼相比, 有1個號碼連續出現
	ContinueCol2    // 跟前一期號碼相比, 有2個號碼連續出現
	ContinueCol3    // 跟前一期號碼相比, 有3個號碼連續出現
	ContinueCol4    // 跟前一期號碼相比, 有4個號碼連續出現

)

const (
	ContinueRowNone = iota
	ContinueRow2    // 同一期出現相連號碼(2個) ex: 01 05 06 23 33
	ContinueRow3    // 同一期出現相連號碼(3個) ex: 01 05 06 07 33
	ContinueRow4    // 同一期出現相連號碼(4個) ex: 01 05 06 07 08
	ContinueRow5    // 同一期出現相連號碼(5個) ex: 04 05 06 07 08
	ContinueRow22   // 同一期出現相連號碼(2個*2)
	ContinueRow32   // 同一期出現相連號碼(2個, 3個 葫蘆)

)

const (
	Next = iota
	NextOnly
	Before
	BeforeOnly
	Both
	None
)

/*
*

	特徵值定義
*/
const (
	FeatureTenGroup1 = iota // 1~10
	FeatureTenGroup2        // 11~20
	FeatureTenGroup3        // 21~30
	FeatureTenGroup4        // 31~39
)

const (
	TailDigit1 = iota
	TailDigit2
	TailDigit3
	TailDigit4
	TailDigit5
	TailDigit6
	TailDigit7
	TailDigit8
	TailDigit9
	TailDigit0
)

var Primes = []byte{1, 2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47}

const (
	FilterOddCount = iota
	FilterEvenCount
	FilterTenGroupOddCount
	FilterTenGroupEvenCount
	FilterTenGroup
	FilterTailDigit
	FilterPrime
	FilterPrimeCount
	FilterContinueRowType
	FilterDisableAll
)

var filters = []bool{
	true, // odd count
	true, // even count
	true, // ten group odd count
	true, // ten group even count
	true, // ten group
	true, // tail digit
	true, // prime
	true, // prime count
	true, // continue type
	true, // disable all
}

func DisableFilters(fs []int) {
	for _, i := range fs {
		filters[i] = false
	}
}

func setFilter(fs []bool) {
	if len(filters) != len(fs) {
		fmt.Errorf("Filter Format Error %d:%d", len(filters), len(fs))
		return
	}
	filters = fs
}

type GROUP int

const UndefinedFeature = -1

type Feature struct {
	IBalls                  []int   `json:"iballs"`
	TenGroupCount           [5]int  `json:"tengroupcount"`
	OddNumberCount          int     `json:"oddnumbercount"`
	TenGroupOddNumberCount  [5]int  `json:"tengroupoddnumbercount"`
	TenGroupEvenNumberCount [5]int  `json:"tengroupevennumbercount"`
	EvenNumberCount         int     `json:"evennumbercount"`
	TailDigit               [10]int `json:"taildigit"`
	PrimeCount              int     `json:"primecount"`
	MultiplesOfs            [19]int `json:"multiplesofs"`
	ContinueRowType         int     `json:"continuerowtype"`
	Point                   float32 `json:"point"`
}

func NewFeature(numbers []int, ballsCount int) *Feature {
	oc := 0
	ec := 0
	gt := [5]int{0, 0, 0, 0, 0}
	td := [10]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	tgonc := [5]int{0, 0, 0, 0, 0}
	tgenc := [5]int{0, 0, 0, 0, 0}
	primec := 0
	for _, n := range numbers {
		if n%2 == 1 {
			oc++
			tgonc[n/10]++
		}
		if n%2 == 0 {
			ec++
			tgenc[n/10]++
		}

		gt[n/10]++
		td[n%10]++
		if bytes.IndexByte(Primes, byte(n)) >= 0 {
			primec++
		}
	}
	f := &Feature{
		IBalls:                  numbers,
		TenGroupCount:           gt,
		OddNumberCount:          oc,
		TenGroupOddNumberCount:  tgonc,
		TenGroupEvenNumberCount: tgenc,
		EvenNumberCount:         ec,
		TailDigit:               td,
		PrimeCount:              primec,
	}
	f.setContinueRowType()
	return f
}

func DefaultFeature() *Feature {
	return &Feature{
		TenGroupCount:   [5]int{UndefinedFeature},
		OddNumberCount:  UndefinedFeature,
		EvenNumberCount: UndefinedFeature,
		TailDigit:       [10]int{UndefinedFeature},
	}
}

func (f *Feature) setContinueRowType() {
	if f.IsContinue2() {
		f.ContinueRowType = ContinueRow2
	} else if f.IsContinue3() {
		f.ContinueRowType = ContinueRow3
	} else if f.IsContinue4() {
		f.ContinueRowType = ContinueRow4
	} else if f.IsContinue5() {
		f.ContinueRowType = ContinueRow5
	} else if f.IsContinue22() {
		f.ContinueRowType = ContinueRow22
	} else if f.IsContinue2() && f.IsContinue3() {
		f.ContinueRowType = ContinueRow32
	} else {
		f.ContinueRowType = ContinueRowNone
	}
}

func (f *Feature) CompareWithFilter(t *Feature, fs []bool) bool {
	setFilter(fs)
	return f.Compare(t)
}

func (f *Feature) Compare(t *Feature) bool {
	if !filters[FilterDisableAll] {
		return true
	}

	if filters[FilterOddCount] {
		if f.OddNumberCount != t.OddNumberCount {
			return false
		}
	}

	if filters[FilterEvenCount] {
		if f.EvenNumberCount != t.EvenNumberCount {
			return false
		}
	}

	if filters[FilterTenGroup] {
		for idx, i := range f.TenGroupCount {
			if i != t.TenGroupCount[idx] {
				return false
			}
		}
	}

	if filters[FilterTailDigit] {
		for idx, i := range f.TailDigit {
			if i != t.TailDigit[idx] {
				return false
			}
		}
	}

	if filters[FilterTenGroupOddCount] {
		for idx, i := range f.TenGroupOddNumberCount {
			if i != t.TenGroupOddNumberCount[idx] {
				return false
			}
		}
	}

	if filters[FilterTenGroupEvenCount] {
		for idx, i := range f.TenGroupEvenNumberCount {
			if i != t.TenGroupEvenNumberCount[idx] {
				return false
			}
		}
	}

	if filters[FilterPrime] {
		if f.PrimeCount > 0 && t.PrimeCount > 0 {

		} else {
			return false
		}
	}

	if filters[FilterPrimeCount] {
		if f.PrimeCount != t.PrimeCount {
			return false
		}
	}

	if filters[FilterContinueRowType] {
		if f.ContinueRowType != t.ContinueRowType {
			return false
		}
	}

	return true
}

func (f *Feature) RCompare(t *Feature) bool {

	if filters[FilterOddCount] {
		if f.OddNumberCount == t.OddNumberCount {
			return true
		}
	}

	if filters[FilterEvenCount] {
		if f.EvenNumberCount == t.EvenNumberCount {
			return true
		}
	}

	if filters[FilterTenGroup] {
		for idx, i := range f.TenGroupCount {
			if i == t.TenGroupCount[idx] {
				return true
			}
		}
	}

	if filters[FilterTailDigit] {
		for idx, i := range f.TailDigit {
			if i == t.TailDigit[idx] {
				return true
			}
		}
	}

	if filters[FilterTenGroupOddCount] {
		for idx, i := range f.TenGroupOddNumberCount {
			if i == t.TenGroupOddNumberCount[idx] {
				return true
			}
		}
	}

	if filters[FilterTenGroupEvenCount] {
		for idx, i := range f.TenGroupEvenNumberCount {
			if i == t.TenGroupEvenNumberCount[idx] {
				return true
			}
		}
	}

	if filters[FilterPrime] {
		if f.PrimeCount > 0 && t.PrimeCount > 0 {
			return true
		}
	}

	if filters[FilterPrimeCount] {
		if f.PrimeCount == t.PrimeCount {
			return true
		}
	}

	if filters[FilterContinueRowType] {
		if f.ContinueRowType != t.ContinueRowType {
			return true
		}
	}

	return false
}

func (fa *Feature) IsContinue2() bool {
	i1 := fa.IBalls[0]
	i2 := fa.IBalls[1]
	i3 := fa.IBalls[2]
	i4 := fa.IBalls[3]
	i5 := fa.IBalls[4]
	return (i2-i1 == 1 || i3-i2 == 1 || i4-i3 == 1 || i5-i4 == 1) && !fa.IsContinue3() && !fa.IsContinue4() && !fa.IsContinue5() && !fa.IsContinue22()
}
func (fa *Feature) IsContinue3() bool {
	i1 := fa.IBalls[0]
	i2 := fa.IBalls[1]
	i3 := fa.IBalls[2]
	i4 := fa.IBalls[3]
	i5 := fa.IBalls[4]
	return ((i2-i1 == 1 && i3-i2 == 1) || (i3-i2 == 1 && i4-i3 == 1) || (i4-i3 == 1 && i5-i4 == 1)) && !fa.IsContinue4() && !fa.IsContinue5()
}

func (fa *Feature) IsContinue4() bool {
	i1 := fa.IBalls[0]
	i2 := fa.IBalls[1]
	i3 := fa.IBalls[2]
	i4 := fa.IBalls[3]
	i5 := fa.IBalls[4]
	return ((i2-i1 == 1 && i3-i2 == 1 && i4-i3 == 1) || (i3-i2 == 1 && i4-i3 == 1 && i5-i4 == 1)) && !fa.IsContinue5()
}

func (fa *Feature) IsContinue5() bool {
	i1 := fa.IBalls[0]
	i2 := fa.IBalls[1]
	i3 := fa.IBalls[2]
	i4 := fa.IBalls[3]
	i5 := fa.IBalls[4]
	return (i2-i1 == 1 && i3-i2 == 1 && i4-i3 == 1 && i5-i4 == 1)
}

func (fa *Feature) IsContinue22() bool {
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

func (f *Feature) Presentation() string {
	return fmt.Sprintf("Balls:%v TenGroup : %v, Odd:Even==%d:%d, OddTen:EvenTen===%v:%v, DigitTail : %v, PrimeCount:%d",
		f.IBalls,
		f.TenGroupCount,
		f.OddNumberCount, f.EvenNumberCount,
		f.TenGroupOddNumberCount, f.TenGroupEvenNumberCount,
		f.TailDigit,
		f.PrimeCount)
}

func (f *Feature) NoContinue() bool {
	return !f.IsContinue2() && !f.IsContinue3() && !f.IsContinue4() && !f.IsContinue5() && !f.IsContinue22()
}

func (f *Feature) IsJump2() bool {
	return f.IBalls[0]-f.IBalls[1] == 1 || f.IBalls[1]-f.IBalls[2] == 1 || f.IBalls[2]-f.IBalls[3] == 1 || f.IBalls[3]-f.IBalls[4] == 1
}

func (f *Feature) IsJump3() bool {
	return (f.IBalls[0]-f.IBalls[1] == 1 && f.IBalls[1]-f.IBalls[2] == 1) || (f.IBalls[1]-f.IBalls[2] == 1 && f.IBalls[2]-f.IBalls[3] == 1) || (f.IBalls[2]-f.IBalls[3] == 1 && f.IBalls[3]-f.IBalls[4] == 1)
}

func (f *Feature) IsJump4() bool {
	return f.IBalls[0]-f.IBalls[1] == 1 && f.IBalls[1]-f.IBalls[2] == 1 && f.IBalls[2]-f.IBalls[3] == 1 && f.IBalls[3]-f.IBalls[4] == 1
}
