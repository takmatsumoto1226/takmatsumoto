package bl

func (ar BigLotteryList) findNumbers(numbers []string, nextplus bool) BigLotteryList {
	intersection := BigLotteryList{}
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
			intersection = append(intersection, ns)
			if nextplus && i+1 < len(ar) {
				intersection = append(intersection, ar[i+1])
				intersection = append(intersection, *Empty())
			}
		}

	}
	// Create a set from the first array

	return intersection
}
