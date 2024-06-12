package ftn

import (
	"lottery/model/df"
	"lottery/model/interf"
	"sort"
)

func (ar *FTNsManager) findDate(date string, t int) FTNArray {
	intersection := FTNArray{}

	for i, ns := range ar.List {

		// Check elements in the second array against the set

		if ns.MonthDay == date {

			if (t == df.Before || t == df.Both) && t != df.None {
				intersection = append(intersection, ar.List[i-1])
			}

			intersection = append(intersection, ns)

			if t == df.Next || t == df.Both && t != df.None {
				if i+1 < len(ar.List) {
					intersection = append(intersection, ar.List[i+1])
				}
			}
			intersection = append(intersection, *Empty())
		}

	}
	// Create a set from the first array

	return intersection
}

func (ar FTNArray) StaticTotalNotInclude(r int) float64 {
	count := 0.
	if r > len(ar) {
		return count
	}
	for i, _ := range ar {
		if i >= r {
			fs := ar.WithRange(i-r, 1)
			tops := ar.WithRange(i-(r-1), r)
			result := fs.FilterExcludes(tops, []int{})
			if len(result) > 0 {
				// fmt.Println(fs.Presentation())
				// fmt.Println(tops.Presentation())
				count++
			}
		}
	}
	return count / float64(len(ar))
}

func (ar FTNArray) StaticNumberShowTwiceup(r int) float64 {
	count := 0.
	if r >= len(ar)-1 {
		return count
	}

	for i, _ := range ar {
		if i >= len(ar)-1 {
			continue
		}
		intcounts := []int{}
		var first = ar.GetNode(i)
		staticr := ar.WithRange(i+1, r-1)
		sort.Sort(staticr)
		for _, rf := range staticr {
			intcounts = append(intcounts, first.IncludeNumbers(rf)...)
		}

		static := map[int]int{}
		for _, v := range intcounts {
			if c, ok := static[v]; ok {
				static[v] = c + 1
			} else {
				static[v] = 1
			}
		}

		for _, c := range static {
			if c > r {
				count++
				break
			}
		}
	}
	return count / float64(len(ar))
}

func (ar FTNArray) StaticContinue2Percent(i interf.Interval) float64 {

	sr := ar
	if i.Length > 0 && i.Index+i.Length < len(ar) {
		sr = ar[i.Index : i.Index+i.Length]
	}
	count := 0
	for _, f := range sr {
		if f.Feature.IsContinue2() {
			count++
		}
	}
	return (float64(count) / float64(len(sr))) * 100
}

func (ar FTNArray) StaticContinue3Percent(i interf.Interval) float64 {

	sr := ar
	if i.Length > 0 && i.Index+i.Length < len(ar) {
		sr = ar[i.Index : i.Index+i.Length]
	}
	count := 0
	for _, f := range sr {
		if f.Feature.IsContinue3() {
			count++
		}
	}
	return (float64(count) / float64(len(sr))) * 100
}

func (ar FTNArray) StaticGroupTen(i interf.Interval, t int, v int) float64 {
	sr := ar
	if i.Length > 0 && i.Index+i.Length < len(ar) {
		sr = ar[i.Index : i.Index+i.Length]
	}
	count := 0
	for _, f := range sr {
		if f.Feature.TenGroupCount[t] == v {
			count++
		}
	}
	return (float64(count) / float64(len(sr))) * 100
}
