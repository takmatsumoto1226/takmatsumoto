package pw

import "lottery/model/df"

func (ar *PowerManager) findNumbers(numbers []string, nextplus int) PowerList {
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
			if nextplus == df.Next {
				intersection = append(intersection, ns)
				if i+1 < len(ar.List) {
					intersection = append(intersection, ar.List[i+1])
					intersection = append(intersection, *Empty())
				}
			} else {
				intersection = append(intersection, ar.List[i-1])
				intersection = append(intersection, ns)
				intersection = append(intersection, *Empty())
			}

		}

	}

	return intersection
}
