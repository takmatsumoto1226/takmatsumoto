package pw

import "lottery/model/df"

func (ar *PowerManager) findNumbers(numbers []string, t int) PowerList {
	intersection := PowerList{}
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

			if (t == df.Before || t == df.Both) && i > 0 {
				intersection = append(intersection, ar.List[i-1])
			}

			intersection = append(intersection, ns)

			if t == df.Next || t == df.Both {
				if i+1 < len(ar.List) {
					intersection = append(intersection, ar.List[i+1])
				}
			}
			if t != df.None {
				intersection = append(intersection, *Empty())
			}

		}

	}

	return intersection
}
