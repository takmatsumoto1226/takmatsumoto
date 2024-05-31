package ftn

import (
	"lottery/model/df"
)

func (fa FTNArray) FilterByGroupIndex(group *FTNGroup) FTNArray {
	arr := FTNArray{}
	for _, ftn := range fa {
		if _, gcount := group.FindGroupIndex(ftn); gcount == 0 {
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
			if b.Count > 3 {
				numbers = append(numbers, b.Ball.Number)
			}
		}

		for _, b := range numbers {
			arr = append(arr, fa.findNumbers([]string{b}, df.None)...)
		}
	}
	return arr.Distinct()
}

func (fa FTNArray) FilterExcludes(tops FTNArray) FTNArray {
	result := FTNArray{}
	search := map[int]bool{}
	for _, t := range tops {
		for _, i := range t.Feature.IBalls {
			search[i] = true
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
		if ftn.Feature.IsContinue2() {
			result = append(result, ftn)
		}
	}
	return result
}
