package pw

import (
	"fmt"
	"lottery/model/common"
	"lottery/model/df"
)

func (fa PowerList) FilterByGroupIndex(group *PWGroup, cs []int) PowerList {
	fmt.Printf("FilterByGroupIndex : %d\n", len(fa))
	arr := PowerList{}
	for _, ftn := range fa {
		for _, c := range cs {
			if _, gcount := group.FindGroupIndex(ftn); gcount == c {
				arr = append(arr, ftn)
				break
			}
		}
	}
	return arr.Distinct()
}

func (fa PowerList) FilterHighFreqNumber(highFreqs PowerList, p PickParam) PowerList {
	fmt.Printf("FilterHighFreqNumber : %d\n", len(fa))
	result := PowerList{}
	ballsCount := highFreqs.IntervalBallsCountStatic(p)
	fmt.Println(ballsCount.Presentation(false))

	numbers := []string{}
	for _, b := range ballsCount {
		if b.Count > uint(p.Freq) {
			numbers = append(numbers, b.Number)
		}
	}

	for _, b := range numbers {
		result = append(result, fa.findNumbers([]string{b}, df.None)...)
	}
	return result
}

func (fa PowerList) Distinct() PowerList {
	fmt.Printf("Distinct : %d\n", len(fa))
	results := PowerList{}
	tmp := map[string]Power{}
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

func (fa PowerList) FilterPickBySpecConfition() PowerList {
	fmt.Printf("FilterPickBySpecConfition : %d\n", len(fa))
	result := PowerList{}
	for _, pw := range fa {
		if pw.Feature.IsContinue3() || pw.Feature.IsContinue22() {
			result = append(result, pw)
		}
	}
	return result
}

func (fa PowerList) FilterIncludes(tops PowerList, sb []int) PowerList {
	fmt.Printf("FilterIncludes : %d\n", len(fa))
	result := PowerList{}
	search := common.LIMap{}
	for _, t := range tops {
		for _, i := range t.Feature.IBalls {
			search[i] = true
		}
	}

	if len(sb) > 0 {
		for _, i := range sb {
			search[i] = true
		}
	}
	fmt.Println(search.Presentation())
	if len(search) == 0 {
		return fa
	}

	for _, ftn := range fa {
		for _, n := range ftn.Feature.IBalls {
			if _, ok := search[n]; ok {
				result = append(result, ftn)
				break
			}
		}
	}
	return result
}

func (fa PowerList) FilterExcludes(tops PowerList, sb []int) PowerList {
	fmt.Printf("FilterExcludes : %d\n", len(fa))
	result := PowerList{}
	search := common.LIMap{}
	for _, t := range tops {
		for _, i := range t.Feature.IBalls {
			search[i] = true
		}
	}
	if len(sb) > 0 {
		for _, b := range sb {
			search[b] = true
		}
	}

	fmt.Println(search.Presentation())
	if len(search) == 0 {
		return fa
	}

	for _, ftn := range fa {
		add := true
		for _, n := range ftn.Feature.IBalls {
			if _, ok := search[n]; ok {
				add = false
				break
			}
		}

		if add {
			result = append(result, ftn)
		}
	}
	return result
}

func (fa PowerList) FilterExcludeNode(tops PowerList) PowerList {
	fmt.Printf("FilterExcludeNode : %d\n", len(fa))
	result := PowerList{}
	sames := PowerList{}
	for _, f := range fa {
		add := true
		for _, t := range tops {
			if f.IsSame(&t) {
				sames = append(sames, f)
				add = false
				break
			}
		}
		if add {
			result = append(result, f)
		}
	}
	if len(sames) > 0 {
		fmt.Println("same ....")
		for _, s := range sames {
			s.ShowRow()
		}
		fmt.Println("so much...")
	}
	return result
}

func (fa PowerList) FilterFeatureExcludes(tops PowerList) PowerList {
	fmt.Printf("FilterFeatureExcludes : %d\n", len(fa))
	result := PowerList{}

	for _, pw := range fa {
		add := true
		for _, top := range tops {
			if pw.MatchFeature(&top) {
				add = false
				break
			}
		}

		if add {
			result = append(result, pw)
		}
	}
	return result
}

func (fa PowerList) FilterNeighber(top *Power, c int) PowerList {
	fmt.Printf("FilterNeighberNumber : %d\n", len(fa))
	result := PowerList{}
	for _, f := range fa {
		if f.haveNeighber(top, c) {
			result = append(result, f)
		}
	}
	return result
}

func (fa PowerList) FilterCol(top *Power, c int) PowerList {
	fmt.Printf("FilterCol : %d\n", len(fa))
	result := PowerList{}
	for _, f := range fa {
		if f.haveCol(top, c) {
			result = append(result, f)
		}
	}
	return result
}

func (fa PowerList) FilterByTenGroup(tt []int, hh []int) PowerList {
	fmt.Printf("FilterByTebGroup : %d\n", len(fa))

	result := PowerList{}
	if len(tt) == 0 {
		for _, f := range fa {
			if f.Feature.IsFullTenGrouop() {
				result = append(result, f)
			}
		}
	} else {
		for _, f := range fa {
			count := 0
			for ti, t := range tt {
				if f.Feature.TenGroupCount[t] != hh[ti] {
					break
				}
				count++
				if count == len(tt) {
					result = append(result, f)
				}
			}
		}
	}

	return result
}

func (fa PowerList) FilterOddEvenList(oc int) PowerList {
	fmt.Printf("FilterOddEvenList : %d\n", len(fa))
	result := PowerList{}
	for _, f := range fa {
		if f.Feature.OddNumberCount == oc {
			result = append(result, f)
		}
	}
	return result
}

func (fa PowerList) FilterPrime(c int) PowerList {
	result := PowerList{}
	for _, f := range fa {
		if f.EqualPrime(c) {
			result = append(result, f)
		}
	}
	return result
}
