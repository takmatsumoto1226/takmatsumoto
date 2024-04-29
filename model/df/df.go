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
	Descending = iota // raw data 年份大到小
	Ascending         // raw data 年份小到大
)

const (
	Biggerfront = iota // 球數出現次數統計後, 出現次數多得在前面
	Smallfront         // 球數出現次數統計後, 出現次數少的在前面
	Normal             // ball的數字由小到大排序
)

// 特徵種類
const (
	ContinuePickupNumber1 = iota // 跟前一期號碼相比, 有1個號碼連續出現
	ContinuePickupNumber2        // 跟前一期號碼相比, 有2個號碼連續出現
	ContinuePickupNumber3        // 跟前一期號碼相比, 有3個號碼連續出現
	ContinuePickupNumber4        // 跟前一期號碼相比, 有4個號碼連續出現
	ContinueNumber2              // 同一期出現相連號碼(2個) ex: 01 05 06 23 33
	ContinueNumber3              // 同一期出現相連號碼(3個) ex: 01 05 06 07 33
	ContinueNumber4              // 同一期出現相連號碼(3個) ex: 01 05 06 07 08
	ContinueNumber5              // 同一期出現相連號碼(3個) ex: 04 05 06 07 08
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

var Primes = []byte{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47}

const (
	FilterOddCount = iota
	FilterEvenCount
	FilterTenGroupOddCount
	FilterTenGroupEvenCount
	FilterTenGroup
	FilterTailDigit
	FilterPrime
	FilterPrimeCount
	FilterContinue2
	FilterContinue3
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
	true, // continue2
}

func DistableFilters(fs []int) {
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
	TenGroupCount           []int
	OddNumberCount          int
	TenGroupOddNumberCount  []int
	TenGroupEvenNumberCount []int
	EvenNumberCount         int
	TailDigit               []int
	PrimeCount              int
	MultiplesOfs            []int // 2,3,....helf of ball count
}

func NewFeature(numbers []int, ballsCount int) *Feature {
	oc := 0
	ec := 0
	gt := []int{0, 0, 0, 0, 0}
	td := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	tgonc := []int{0, 0, 0, 0, 0}
	tgenc := []int{0, 0, 0, 0, 0}
	primec := 0
	for i := 0; i < len(numbers); i++ {
		if numbers[i]%2 == 1 {
			oc++
			tgonc[numbers[i]/10]++
		}
		if numbers[i]%2 == 0 {
			ec++
			tgenc[numbers[i]/10]++
		}

		gt[numbers[i]/10]++
		td[numbers[i]%10]++
		if bytes.IndexByte(Primes, byte(numbers[i])) >= 0 {
			primec++
		}
	}
	return &Feature{
		TenGroupCount:           gt,
		OddNumberCount:          oc,
		TenGroupOddNumberCount:  tgonc,
		TenGroupEvenNumberCount: tgenc,
		EvenNumberCount:         ec,
		TailDigit:               td,
		PrimeCount:              primec,
	}
}

func DefaultFeature() *Feature {
	return &Feature{
		TenGroupCount:   []int{UndefinedFeature},
		OddNumberCount:  UndefinedFeature,
		EvenNumberCount: UndefinedFeature,
		TailDigit:       []int{UndefinedFeature},
	}
}

func (f *Feature) CompareWithFilter(t *Feature, fs []bool) bool {
	setFilter(fs)
	return f.Compare(t)
}

func (f *Feature) Compare(t *Feature) bool {

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
			j := t.TenGroupCount[idx]
			if i != j {
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
			return false
		}
	}

	if filters[FilterPrimeCount] {
		if f.PrimeCount != t.PrimeCount {
			return false
		}
	}

	return true
}
