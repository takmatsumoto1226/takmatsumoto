package ftn

import "lottery/model/df"

func (ar *FTNsManager) findNumbers(numbers []string, t int) FTNArray {
	intersection := FTNArray{}
	set := make(map[string]bool)

	for i, ns := range ar.List {
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

			if t == df.Before || t == df.Both {
				intersection = append(intersection, ar.List[i-1])
			}

			intersection = append(intersection, ns)

			if t == df.Next || t == df.Both {
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

// Number MA - 移動平均數字
type BallMA struct {
	Interval int
}
