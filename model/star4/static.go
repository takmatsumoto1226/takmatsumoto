package star4

func (ar *Star4Manager) findNumbers(numbers []string, combo bool) Star4s {
	intersection := Star4s{}
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
			if combo {
				intersection = append(intersection, ns)
				if i+1 < len(ar.List) {
					intersection = append(intersection, ar.List[i+1])
					intersection = append(intersection, *Empty())
				}
			}
		}
	}
	// Create a set from the first array

	return intersection
}

// func removeDuplicates(s Star4s) Star4s {
// 	bucket := make(map[string]bool)
// 	var result []string
// 	for _, str := range s {
// 		if _, ok := bucket[str]; !ok {
// 			bucket[str] = true
// 			result = append(result, str)
// 		}
// 	}
// 	return result
// }
