package ftn

import (
	"errors"
	"fmt"
	"lottery/algorithm"
	"lottery/model/df"
	"lottery/model/interf"
	"strconv"

	"github.com/sirupsen/logrus"
)

// FTNArray ...
type FTNArray []FTN

func NewFTNArray(numberss [][]string) *FTNArray {
	arr := FTNArray{}
	for _, numbers := range numberss {
		ftn := NewFTNWithStrings(numbers)
		arr = append(arr, *ftn)
	}
	return &arr
}

func NewFTNArrayWithInts(numberss [][]int) *FTNArray {
	arr := FTNArray{}
	for _, numbers := range numberss {
		ftn := NewFTNWithInts(numbers)
		arr = append(arr, *ftn)
	}
	return &arr
}

func (fa *FTNArray) Head() {
	rowmsg := "====|====|"
	for i := 1; i <= ballsCountFTN; i++ {
		rowmsg = rowmsg + fmt.Sprintf("%02d|", i)
	}
	fmt.Println(rowmsg)
	fmt.Println("")
}

func (fa FTNArray) Len() int {
	return len(fa)
}

// Less ...
func (fa FTNArray) Less(i, j int) bool {
	ii, _ := strconv.Atoi(fa[i].TIdx)
	jj, _ := strconv.Atoi(fa[j].TIdx)
	return ii > jj
}

// Swap swaps the elements with indexes i and j.
func (fa FTNArray) Swap(i, j int) {
	fa[i], fa[j] = fa[j], fa[i]
}
func (fa FTNArray) Presentation() string {
	return fa.PresentationWithRange(0)
}

func (fa FTNArray) PresentationWithRange(r int) string {
	msg := ""
	tmp := fa
	al := len(fa)
	if r > 0 {
		tmp = fa[al-r : al]
	}
	for _, ftn := range tmp {
		msg = msg + ftn.formRow() + "\n"

	}
	return msg
}

func (fa FTNArray) ShowWithRange(r int) {
	tmp := fa
	al := len(fa)
	if r > 0 {
		tmp = fa[al-r : al]
	}
	for _, ftn := range tmp {
		ftn.ShowRow()
	}
}

func (fa FTNArray) ShowAll() {
	fa.ShowWithRange(0)
}

func (fa FTNArray) WithRange(i, r int) FTNArray {
	al := len(fa)
	if r > 0 {
		return fa[al-r-i : al-i]
	}
	return fa
}

func (fa FTNArray) FeatureRange(th interf.Threshold) FTNArray {
	features := fa.WithRange(th.Interval.Index, th.Interval.Length)
	if th.Smart.Enable {
		if th.Smart.Type == interf.RangeTypeLatest {
			features = append(features, fa.SmartWithTh(*interf.PureIntervalTH(0, 1))...)
		} else if th.Smart.Type == interf.RangeTypeLatestSame {
			features = append(features, fa.SmartWithFeature(*interf.PureIntervalTH(0, 1))...)
		} else {
			features = append(features, fa.SmartWithTh(th)...)
		}
	}
	return features.Distinct()
}

func (fa FTNArray) Distinct() FTNArray {
	results := FTNArray{}
	tmp := map[string]FTN{}
	for _, f := range fa {
		if _, ok := tmp[f.Key()]; !ok {
			tmp[f.Key()] = f
		}
	}

	for _, v := range tmp {
		results = append(results, v)
	}
	return results
}

func (fa FTNArray) SmartWithTh(th interf.Threshold) FTNArray {
	features := fa.WithRange(th.Interval.Index, th.Interval.Length)
	for _, bs := range features {
		result := algorithm.Combinations(bs.toStringArray(), 3)
		for _, v := range result {
			features = append(features, fa.findNumbers(v, df.NextOnly)...)
		}
	}

	return features
}

func (fa FTNArray) SmartWithFeature(th interf.Threshold) FTNArray {
	features := fa.WithRange(th.Interval.Index, th.Interval.Length)
	if len(features) > 0 {
		latest := features[0]
		for _, bs := range fa {
			if latest.CompareFeature(&bs) {
				features = append(features, bs)
			}
		}
	}
	return features
}

func (list FTNArray) findNumbers(numbers []string, t int) FTNArray {
	intersection := FTNArray{}
	set := make(map[string]bool)

	for i, ns := range list {
		for _, num := range numbers {
			set[num] = true // setting the initial value to true
		}

		// Check elements in the second array against the set
		count := 0
		for _, num := range ns.toStringArray() {
			if set[num] {
				count++
			}
		}

		if len(set) == count {

			if (t == df.BeforeOnly || t == df.Before || t == df.Both) && i > 0 {
				intersection = append(intersection, list[i-1])
			}

			if t != df.NextOnly && t != df.BeforeOnly {
				intersection = append(intersection, ns)
			}

			if t == df.NextOnly || t == df.Next || t == df.Both {
				if i+1 < len(list) {
					intersection = append(intersection, list[i+1])
				}
			}
			if t != df.None && t != df.NextOnly && t != df.BeforeOnly {
				intersection = append(intersection, *Empty())
			}

		}

	}
	// Create a set from the first array

	return intersection
}

func (fa FTNArray) adariPrice(adari *FTN) {
	for i := 5; i > 1; i-- {
		combinations := algorithm.Combinations(adari.toStringArray(), i)
		// fmt.Println(combinations)
		ftnarr := NewFTNArray(combinations)
		for _, ftn := range *ftnarr {
			for _, fav := range fa {
				mc := fav.matchCount(ftn)
				if mc == 5 {
					fmt.Println("8000000元")
				} else if mc == 4 {
					fmt.Println("20000元")
				} else if mc == 3 {
					fmt.Println("300元")
				} else {
					fmt.Println("50元")
				}
			}
		}
	}
}

func (ar FTNArray) intervalBallsCountStatic(p PickParam) map[uint]NormalizeInfo {

	ballsCount := map[uint]NormalizeInfo{}

	if p.Interval == 0 {
		logrus.Error(errors.New("不可指定0"))
		return map[uint]NormalizeInfo{}
	}
	var FTNIntervalCount = [ballsCountFTN]uint{}

	for _, t := range ar {
		FTNIntervalCount[numberToIndex[t.B1.Number]]++
		FTNIntervalCount[numberToIndex[t.B2.Number]]++
		FTNIntervalCount[numberToIndex[t.B3.Number]]++
		FTNIntervalCount[numberToIndex[t.B4.Number]]++
		FTNIntervalCount[numberToIndex[t.B5.Number]]++
	}
	arr := BallsCount{}
	for i, count := range FTNIntervalCount {
		b := BallInfo{Count: count, Ball: Ball{fmt.Sprintf("%02d", i+1), i, i + 1, 0}}
		arr = append(arr, b)
	}
	ballsCount[p.Interval] = NormalizeInfo{NorBalls: arr, Param: p}

	return ballsCount
}
