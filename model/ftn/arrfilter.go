package ftn

import (
	"fmt"
	"lottery/model/common"
	"lottery/model/df"
	"time"
)

func (fa FTNArray) FilterByGroupIndex(group *FTNGroup, cs []int) FTNArray {
	defer common.TimeTaken(time.Now(), fmt.Sprintf("FilterByGroupIndex : %d\n", len(fa)))
	arr := FTNArray{}
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

func (fa FTNArray) FilterMatchBall(params []PickParam, staticmap map[string]BallsInfo) FTNArray {
	arr := FTNArray{}
	for _, p := range params {
		fa.IntervalBallsCountStatic(p)
		static := staticmap[p.GetKey()]
		numbers := []string{}
		for _, b := range static {
			if b.Count == 0 {
				numbers = append(numbers, b.Ball.Number)
			}
		}

		for _, b := range numbers {
			arr = append(arr, fa.findNumbers([]string{b}, df.None)...)
		}
	}
	return arr.Distinct()
}

func (fa FTNArray) FilterExcludes(tops FTNArray, sb []int) FTNArray {
	defer common.TimeTaken(time.Now(), fmt.Sprintf("FilterExcludes : %d\n", len(fa)))
	result := FTNArray{}
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

	if len(search) == 0 {
		return fa
	}

	// fmt.Println(search.Presentation())

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

func (fa FTNArray) FilterIncludes(tops FTNArray, sb []int) FTNArray {
	defer common.TimeTaken(time.Now(), fmt.Sprintf("FilterIncludes : %d\n", len(fa)))

	// if len(tops) == 0 && len(sb) == 0 {
	// 	return fa
	// }
	result := FTNArray{}
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

	if len(search) == 0 {
		return fa
	}
	fmt.Println(search.Presentation())

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

func (fa FTNArray) FilterHighFreqNumber(highFreqs FTNArray, p PickParam) FTNArray {
	defer common.TimeTaken(time.Now(), fmt.Sprintf("FilterHighFreqNumber : %d\n", len(fa)))
	result := FTNArray{}
	ballsCount := highFreqs.IntervalBallsCountStatic(p)
	fmt.Println(ballsCount.AppearBalls.Presentation(false))

	total := uint(0)
	for _, bc := range ballsCount.AppearBalls {
		total = total + bc.Count
	}

	numbers := []string{}
	for _, b := range ballsCount.AppearBalls {
		if b.Count > uint(p.Freq) {
			numbers = append(numbers, b.Ball.Number)
		}
	}

	for _, b := range numbers {
		result = append(result, fa.findNumbers([]string{b}, df.None)...)
	}
	return result.Distinct()
}

func (fa FTNArray) FilterPickBySpecConfition(cts []int) FTNArray {
	defer common.TimeTaken(time.Now(), fmt.Sprintf("FilterPickBySpecConfition : %d\n", len(fa)))
	result := FTNArray{}
	for _, ftn := range fa {
		for _, ct := range cts {
			if ct == df.ContinueRowNone && ftn.Feature.IsContinueNo() {
				result = append(result, ftn)
			} else if ct == df.ContinueRow2 && ftn.Feature.IsContinue2() {
				result = append(result, ftn)
			} else if ct == df.ContinueRow3 && ftn.Feature.IsContinue3() {
				result = append(result, ftn)
			} else if ct == df.ContinueRow4 && ftn.Feature.IsContinue4() {
				result = append(result, ftn)
			} else if ct == df.ContinueRow22 && ftn.Feature.IsContinue22() {
				result = append(result, ftn)
			}
		}
	}
	return result
}

func (fa FTNArray) FilterFeatureExcludes(tops FTNArray) FTNArray {
	defer common.TimeTaken(time.Now(), fmt.Sprintf("FilterFeatureExcludes : %d\n", len(fa)))
	result := FTNArray{}

	for _, ftn := range fa {
		add := true
		for _, top := range tops {
			if ftn.MatchFeature(&top) {
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

func (fa FTNArray) FilterFeatureIncludes(tops FTNArray) FTNArray {
	defer common.TimeTaken(time.Now(), fmt.Sprintf("FilterFeatureIncludes : %d\n", len(fa)))
	result := FTNArray{}

	for _, ftn := range fa {
		for _, top := range tops {
			if ftn.MatchFeature(&top) {
				result = append(result, ftn)
				break
			}
		}
	}
	return result
}

func (fa FTNArray) FilterExcludeNote(tops FTNArray) FTNArray {
	defer common.TimeTaken(time.Now(), fmt.Sprintf("FilterExcludeNote : %d\n", len(fa)))
	result := FTNArray{}
	sames := FTNArray{}
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

	return result
}

func (fa FTNArray) FilterNeighber(top *FTN, cs []int) FTNArray {
	defer common.TimeTaken(time.Now(), fmt.Sprintf("FilterNeighberNumber : %d\n", len(fa)))
	if len(cs) == 0 {
		return fa
	}
	result := FTNArray{}
	for _, f := range fa {
		for _, c := range cs {
			if f.haveNeighber(top, c) {
				result = append(result, f)
			}
		}
	}
	return result
}

func (fa FTNArray) FilterCol(top *FTN, cs []int) FTNArray {
	defer common.TimeTaken(time.Now(), fmt.Sprintf("FilterCol : %d\n", len(fa)))
	result := FTNArray{}
	for _, f := range fa {
		for _, c := range cs {
			if f.haveCol(top, c) {
				result = append(result, f)
			}
		}
	}
	return result
}

func (fa FTNArray) FilterByTenGroupLog(tt []int, hh []int) FTNArray {
	defer common.TimeTaken(time.Now(), fmt.Sprintf("FilterByTenGroup : %d\n", len(fa)))
	return fa.FilterByTenGroup(tt, hh)
}

func (fa FTNArray) FilterByTenGroup(tt []int, hh []int) FTNArray {
	result := FTNArray{}
	if len(tt) == 0 {
		for _, f := range fa {
			if f.Feature.IsFullTenGroup() {
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

func (fa FTNArray) FilterByTebGroupC(tt []int, hhh [][]int) FTNArray {
	defer common.TimeTaken(time.Now(), fmt.Sprintf("FilterTebGroup : %d\n", len(fa)))
	if len(tt) == 0 {
		return fa
	}

	result := FTNArray{}
	if len(tt) == 0 {
		for _, f := range fa {
			if f.Feature.IsFullTenGroup() {
				result = append(result, f)
			}
		}
	} else {
		for _, f := range fa {
			add := false
			for ti, t := range tt {
				count := 0
				for _, hh := range hhh {
					if f.Feature.TenGroupCount[t] != hh[ti] {
						break
					}
					count++
					if count == len(tt) {
						result = append(result, f)
						add = true
						break
					}
				}
				if add {
					break
				}
			}
		}
	}

	return result
}

func (fa FTNArray) FilterOddEvenList(ocs []int) FTNArray {
	defer common.TimeTaken(time.Now(), fmt.Sprintf("FilterOddEvenList : %d\n", len(fa)))
	result := FTNArray{}
	for _, f := range fa {
		for _, oc := range ocs {
			if f.Feature.OddNumberCount == oc {
				result = append(result, f)
			}
		}
	}
	return result
}

func (fa FTNArray) FilterPrime(cs []int) FTNArray {
	defer common.TimeTaken(time.Now(), fmt.Sprintf("FilterPrime : %d\n", len(fa)))
	if len(cs) == 0 {
		return fa
	}
	result := FTNArray{}
	for _, f := range fa {
		for _, c := range cs {
			if f.EqualPrime(c) {
				result = append(result, f)
			}
		}
	}
	return result
}

func (fa FTNArray) FilterColN(n int) FTNArray {
	result := FTNArray{}
	for _, f := range fa {
		if f.B1.Continue == n || f.B2.Continue == n || f.B3.Continue == n || f.B4.Continue == n || f.B5.Continue == n {
			result = append(result, f)
		}
	}
	return result
}

func (fa FTNArray) FilterPeriodN(n, p int) FTNArray {
	result := FTNArray{}
	for i, f := range fa {
		if f.B1.Disappear(n, p) || f.B2.Disappear(n, p) || f.B3.Disappear(n, p) || f.B4.Disappear(n, p) || f.B5.Disappear(n, p) {
			result = append(result, f)
			result = append(result, fa[i+1])
			result = append(result, *Empty())
		}
	}
	return result
}

func (fa FTNArray) FilterNoContinue() FTNArray {
	result := FTNArray{}
	for _, f := range fa {
		if f.Feature.IsContinueNo() {
			result = append(result, f)
		}
	}
	return result
}

func (fa FTNArray) FilterContinue2() FTNArray {
	result := FTNArray{}
	for _, f := range fa {
		if f.Feature.IsContinue2() {
			result = append(result, f)
		}
	}
	return result
}

func (fa FTNArray) FilterContinue3() FTNArray {
	result := FTNArray{}
	for _, f := range fa {
		if f.Feature.IsContinue3() {
			result = append(result, f)
		}
	}
	return result
}
