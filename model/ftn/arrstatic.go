package ftn

import (
	"lottery/model/df"
	"sort"
)

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

func (ar FTNArray) StaticContinue2Percent() float64 {

	count := 0
	for _, f := range ar {
		if f.Feature.IsContinue2() {
			count++
		}
	}
	return (float64(count) / float64(len(ar))) * 100
}

func (ar FTNArray) StaticContinue3Percent(show bool) float64 {
	count := 0
	for _, f := range ar {
		if f.Feature.IsContinue3() {
			if show {
				f.ShowRow()
			}
			count++
		}
	}
	return (float64(count) / float64(len(ar))) * 100
}

func (ar FTNArray) StaticGroupTen(t int, v int) float64 {
	count := 0
	for _, f := range ar {
		if f.Feature.TenGroupCount[t] == v {
			count++
		}
	}
	return (float64(count) / float64(len(ar))) * 100
}

func (ar FTNArray) StaticColPercent(n int) float64 {
	count := 0.
	for i, f := range ar {
		if i < len(ar)-1 {
			if f.haveCol(&ar[i+1], n) {
				count++
			}
		}
	}
	return (count / float64(len(ar)-1)) * 100
}

func (ar FTNArray) StaticHaveNeighberPercent(n int) float64 {
	count := 0.
	for i, f := range ar {
		if i < len(ar)-1 {
			if f.haveNeighber(&ar[i+1], n) {
				count++
			}
		}
	}
	return (count / float64(len(ar)-1)) * 100
}

func (ar FTNArray) StaticFullTenGroupPercent() float64 {
	count := 0.
	for _, f := range ar {
		if f.Feature.IsFullTenGrouop() {
			count++
		}
	}
	return (count / float64(len(ar)-1)) * 100
}

func (ar FTNArray) StaticExclude(r int, show bool) float64 {
	result := ar.Exclude(r)
	if show {
		result.ShowAll()
	}
	return (float64(len(result)) / float64(len(ar)-r)) * 100
}

func (ar FTNArray) StaticTenGroupByTKey() df.LMap {
	result := df.LMap{}
	for _, f := range ar {
		if v, ok := result[f.Feature.TGKey()]; ok {
			result[f.Feature.TGKey()] = v + 1
		} else {
			result[f.Feature.TGKey()] = 1
		}
	}
	return result

}

func (ar FTNArray) StaticHaveNeighberAndColsPercent(n int, c int) float64 {
	result := ar.NeighberAndCols(n, c)
	return (float64(len(result)) / float64(len(ar)-1)) * 100
}
