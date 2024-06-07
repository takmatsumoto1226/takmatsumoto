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
		s := al - r - i
		e := al - i
		ra := fa[s:e]
		return ra
	}
	return fa
}

func (fa FTNArray) FragmentRange(indexs []int) FTNArray {
	result := FTNArray{}
	for _, i := range indexs {
		result = append(result, fa[i])
	}
	return result
}

func (fa FTNArray) GetFTN(i int) FTN {
	if i <= len(fa) {
		return fa[i]
	}
	return fa[len(fa)-1]
}

func (fa FTNArray) GetFTNWithDate(date string) FTN {
	for _, ftn := range fa {
		if ftn.SameDate(date) {
			return ftn
		}
	}
	return *Empty()
}

func (fa FTNArray) FeatureRange(th interf.Threshold) FTNArray {
	features := fa.WithRange(th.Interval.Index, th.Interval.Length)
	if th.Smart.Enable {
		if th.Smart.Type == interf.RangeTypeLatestRange {
			features = append(features, fa.SmartWithTh(*interf.PureIntervalTH(0, 1))...)
		} else if th.Smart.Type == interf.RangeTypeLatestSame {
			features = append(features, fa.SmartWithFeature(*interf.PureIntervalTH(0, 1))...)
		} else if th.Smart.Type == interf.RangeTypeSpecStartAndRangeNotes {
			features = append(features, fa.SmartWithTh(*interf.PureIntervalTH(th.Interval.Index, 1))...)
		} else if th.Smart.Type == interf.RangeTypeFullHistoryOnly {

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
			results = append(results, f)
		}
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
			if latest.MatchFeature(&bs) {
				features = append(features, bs)
			}
		}
	}
	return features
}

func (fa FTNArray) MatchElements(fb FTNArray) FTNArray {
	ta := FTNArray{}
	for _, a := range fa {
		for _, b := range fb {
			if a.MatchFeature(&b) {
				ta = append(ta, b)
			}
		}
	}
	return ta
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

func (fa FTNArray) AdariPrice(adari *FTN) int {
	total := 0
	for _, pn := range fa {
		currentPrice := pn.AdariPrice(adari)
		total = total + currentPrice
	}
	return 0
}

func (ar FTNArray) IntervalBallsCountStatic(p PickParam) NormalizeInfo {

	if p.Interval == 0 {
		logrus.Error(errors.New("不可指定0"))
		return NormalizeInfo{}
	}
	var FTNIntervalCount = [ballsCountFTN]uint{}
	var disappearCount = [ballsCountFTN]uint{}

	for _, t := range ar {
		FTNIntervalCount[numberToIndex[t.B1.Number]]++
		FTNIntervalCount[numberToIndex[t.B2.Number]]++
		FTNIntervalCount[numberToIndex[t.B3.Number]]++
		FTNIntervalCount[numberToIndex[t.B4.Number]]++
		FTNIntervalCount[numberToIndex[t.B5.Number]]++
		for i := 0; i < ballsCountFTN; i++ {
			if i != numberToIndex[t.B1.Number] ||
				i != numberToIndex[t.B2.Number] ||
				i != numberToIndex[t.B3.Number] ||
				i != numberToIndex[t.B4.Number] ||
				i != numberToIndex[t.B5.Number] {
				disappearCount[i]++
			} else {
				disappearCount[i] = 0
			}
		}
	}

	arr := BallsCount{}
	disarr := BallsCount{}
	for i, count := range FTNIntervalCount {
		b := BallInfo{Count: count, Ball: Ball{fmt.Sprintf("%02d", i+1), i, i + 1, 0, 0}}
		arr = append(arr, b)
		c := BallInfo{Count: disappearCount[i], Ball: Ball{fmt.Sprintf("%02d", i+1), i, i + 1, 0, 0}}
		disarr = append(disarr, c)
	}

	return NormalizeInfo{AppearBalls: arr, Param: p}
}

func (ar FTNArray) BallsCountStatic() BallsCount {
	var FTNIntervalCount = [ballsCountFTN]uint{}
	var disappearCount = [ballsCountFTN]uint{}

	for _, t := range ar {
		FTNIntervalCount[numberToIndex[t.B1.Number]]++
		FTNIntervalCount[numberToIndex[t.B2.Number]]++
		FTNIntervalCount[numberToIndex[t.B3.Number]]++
		FTNIntervalCount[numberToIndex[t.B4.Number]]++
		FTNIntervalCount[numberToIndex[t.B5.Number]]++
		for i := 0; i < ballsCountFTN; i++ {
			if i != numberToIndex[t.B1.Number] ||
				i != numberToIndex[t.B2.Number] ||
				i != numberToIndex[t.B3.Number] ||
				i != numberToIndex[t.B4.Number] ||
				i != numberToIndex[t.B5.Number] {
				disappearCount[i]++
			} else {
				disappearCount[i] = 0
			}
		}
	}

	arr := BallsCount{}
	for i, count := range FTNIntervalCount {
		b := BallInfo{Count: count, Ball: Ball{fmt.Sprintf("%02d", i+1), i, i + 1, 0, 0}}
		arr = append(arr, b)
	}

	return arr
}

func (ar FTNArray) featureBackTesting() map[string]FTNArray {
	result := map[string]FTNArray{}
	for _, ftn := range ar {
		tmpArr := FTNArray{}
		for _, t := range ar {
			if ftn.Key() != t.Key() {
				if ftn.MatchFeature(&t) {
					tmpArr = append(tmpArr, t)
				}
			}
		}
		result[ftn.DateKey()] = tmpArr

	}
	return result
}
