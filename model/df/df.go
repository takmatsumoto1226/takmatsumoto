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
	Before
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
	FilterSingleCount = iota
	FilterDoubleCount
	FilterTenGroup
	FilterTailDigit
	FilterPrime
	FilterContinue2
)

var filters = []bool{
	true, // single count
	true, // double count
	true, // ten group
	true, // tail digit
	true, // prime
	true, // continue2
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
	TenGroupCount   []int
	OddNumberCount  int
	EvenNumberCount int
	TailDigit       []int
	HasPrime        bool
}

func NewFeature(numbers []int) *Feature {
	oc := 0
	ec := 0
	gt := []int{0, 0, 0, 0, 0}
	td := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	prime := false
	for i := 0; i < 5; i++ {
		if numbers[i]%2 == 1 {
			oc++
		}
		if numbers[i]%2 == 0 {
			ec++
		}
		gt[numbers[i]/10]++
		td[numbers[i]%10]++
		if bytes.IndexByte(Primes, byte(numbers[i])) >= 0 {
			prime = true
		}
	}
	return &Feature{
		TenGroupCount:   gt,
		OddNumberCount:  oc,
		EvenNumberCount: ec,
		TailDigit:       td,
		HasPrime:        prime,
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
	if filters[FilterSingleCount] {
		if f.OddNumberCount != t.OddNumberCount {
			return false
		}
	}

	if filters[FilterDoubleCount] {
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

	if filters[FilterPrime] {
		return f.HasPrime == t.HasPrime
	}

	return true
}
