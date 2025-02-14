package bl

import "lottery/model/df"

func (ar BLList) findNumbers(numbers []string, nextplus int) BLList {
	intersection := BLList{}
	set := make(map[string]bool)

	for i, ns := range ar {
		for _, num := range numbers {
			set[num] = true // setting the initial value to true
		}

		// Check elements in the second array against the set
		count := 0
		for _, num := range ns.toStringArray2() {
			if set[num] {
				count++
			}
		}

		if len(set) == count {
			if nextplus == df.Before || nextplus == df.Both {
				intersection = append(intersection, ar[i-1])
				intersection = append(intersection, *Empty())
			}

			intersection = append(intersection, ns)
			if (nextplus == df.Next || nextplus == df.Both) && i+1 < len(ar) {
				intersection = append(intersection, ar[i+1])
				intersection = append(intersection, *Empty())
			}
		}

	}
	// Create a set from the first array

	return intersection
}
