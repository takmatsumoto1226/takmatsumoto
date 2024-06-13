package pw

import (
	"fmt"
	"lottery/algorithm"
	"lottery/model/df"
	"lottery/model/interf"
	"sort"
	"strconv"
)

type PowerList []Power

func (fa PowerList) Len() int {
	return len(fa)
}

// Less ...
func (fa PowerList) Less(i, j int) bool {
	ii, _ := strconv.Atoi(fa[i].TIdx)
	jj, _ := strconv.Atoi(fa[j].TIdx)
	return ii > jj
}

// Swap swaps the elements with indexes i and j.
func (fa PowerList) Swap(i, j int) {
	fa[i], fa[j] = fa[j], fa[i]
}

func (fa PowerList) Presentation() string {
	msg := ""
	for _, f := range fa {
		msg = msg + f.formRow() + "\n"
	}
	return msg
}

func (fa PowerList) WithRange(i, r int) PowerList {
	if r > 0 {
		return fa[i : i+r]
	}
	return fa
}

func (fa PowerList) Reverse() PowerList {
	sort.Sort(sort.Reverse(fa))
	return fa
}

func (fa PowerList) FragmentRange(indexs []int) PowerList {
	result := PowerList{}
	for _, i := range indexs {
		result = append(result, fa[i])
	}
	return result
}

func (fa PowerList) FeatureRange(th interf.Threshold) PowerList {
	lottos := fa.WithRange(th.Interval.Index, th.Interval.Length)
	lottos = append(lottos, fa.WithAI()...)
	return lottos.Distinct()
}

func (fa PowerList) WithAI() PowerList {
	features := PowerList{}
	result := algorithm.Combinations(fa[0].toStringArray(), 3)
	for _, v := range result {
		features = append(features, fa.findNumbers(v, df.NextOnly)...)
	}
	return features
}

func (ar PowerList) findNumbers(numbers []string, t int) PowerList {
	intersection := PowerList{}
	set := make(map[string]bool)

	for i, ns := range ar {
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
				intersection = append(intersection, ar[i-1])
			}

			if t != df.NextOnly && t != df.BeforeOnly {
				intersection = append(intersection, ns)
			}

			if t == df.NextOnly || t == df.Next || t == df.Both {
				if i+1 < len(ar) {
					intersection = append(intersection, ar[i+1])
				}
			}
			if t != df.None && t != df.NextOnly && t != df.BeforeOnly {
				intersection = append(intersection, *Empty())
			}

		}

	}

	return intersection
}

func (fa PowerList) ShowWithRange(r int) {
	fmt.Println(fa[:r].Presentation())
}

func (fa PowerList) PresentationWithRange(r int) string {
	msg := ""
	tmp := fa
	if r > 0 {
		tmp = fa[:r]
	}
	for _, ftn := range tmp {
		msg = msg + ftn.formRow() + "\n"

	}
	return msg
}

func (fa PowerList) ShowAll() {
	fa.ShowWithRange(0)
}

func (fa PowerList) GetNode(i int) Power {
	if i >= len(fa) {
		return fa[0]
	}
	return fa[i]
}
