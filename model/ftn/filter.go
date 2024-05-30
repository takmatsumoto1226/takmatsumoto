package ftn

import "lottery/model/df"

func (ar FTNArray) FilterByGroupIndex(group *FTNGroup) FTNArray {
	arr := FTNArray{}
	for _, ftn := range ar {
		if _, gcount := group.FindGroupIndex(ftn); gcount == 0 {
			arr = append(arr, ftn)
		}
	}
	return arr.Distinct()
}

func (ar FTNArray) FilterMatchBall(params []PickParam, staticmap map[string]BallsInfo) FTNArray {
	arr := FTNArray{}
	for _, p := range params {
		ar.intervalBallsCountStatic(p)
		static := staticmap[p.GetKey()]
		numbers := []string{}
		for _, b := range static {
			if b.Count > 3 {
				numbers = append(numbers, b.Ball.Number)
			}
		}

		for _, b := range numbers {
			arr = append(arr, ar.findNumbers([]string{b}, df.None)...)
		}
	}
	return arr.Distinct()
}
