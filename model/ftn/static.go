package ftn

func (ar *FTNsManager) findNumbers(numbers []string, nextplus bool) FTNArray {
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
			if nextplus {
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
	// Create a set from the first array

	return intersection
}
