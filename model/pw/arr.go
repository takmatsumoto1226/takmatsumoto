package pw

import (
	"lottery/algorithm"
	"lottery/model/df"
	"lottery/model/interf"
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
	al := len(fa)
	if r > 0 {
		return fa[al-r-i : al-i]
	}
	return fa
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
		features = append(features, fa.findNumbers(v, df.Next)...)
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

			if (t == df.Before || t == df.Both) && i > 0 {
				intersection = append(intersection, ar[i-1])
			}

			intersection = append(intersection, ns)

			if t == df.Next || t == df.Both {
				if i+1 < len(ar) {
					intersection = append(intersection, ar[i+1])
				}
			}
			if t != df.None {
				intersection = append(intersection, *Empty())
			}

		}

	}

	return intersection
}

func (fa PowerList) Distinct() PowerList {
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
