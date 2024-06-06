package ftn

import (
	"fmt"
	"lottery/model/df"
)

func (fa FTNArray) FilterByGroupIndex(group *FTNGroup, c int) FTNArray {
	arr := FTNArray{}
	for _, ftn := range fa {
		if _, gcount := group.FindGroupIndex(ftn); gcount == c {
			arr = append(arr, ftn)
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
	result := FTNArray{}
	search := map[int]bool{}
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

func (fa FTNArray) FilterIncludeLatest(tops FTNArray) FTNArray {
	result := FTNArray{}
	search := map[int]bool{}
	for _, t := range tops {
		for _, i := range t.Feature.IBalls {
			search[i] = true
		}
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

func (fa FTNArray) FilterHighFreqNumber(highFreqs FTNArray, p PickParam) FTNArray {
	result := FTNArray{}
	ballsCount := highFreqs.IntervalBallsCountStatic(p)
	fmt.Println(ballsCount.AppearBalls.Presentation(false))

	numbers := []string{}
	for _, b := range ballsCount.AppearBalls {
		if b.Count > uint(p.Freq) {
			numbers = append(numbers, b.Ball.Number)
		}
	}

	for _, b := range numbers {
		result = append(result, fa.findNumbers([]string{b}, df.None)...)
	}
	return result
}

func (fa FTNArray) FilterPickBySpecConfition() FTNArray {
	result := FTNArray{}
	for _, ftn := range fa {
		if ftn.Feature.IsContinue2() || ftn.Feature.NoContinue() {
			result = append(result, ftn)
		}
	}
	return result
}

func (fa FTNArray) FilterFeatureExcludes(tops FTNArray) FTNArray {
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

func (fa FTNArray) FilterExcludeNode(tops FTNArray) FTNArray {
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
	if len(sames) > 0 {
		fmt.Println("same ....")
		for _, s := range sames {
			s.ShowRow()
		}
		fmt.Println("so much...")
	}
	return result
}
