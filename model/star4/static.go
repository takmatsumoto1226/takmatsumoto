package star4

import "fmt"

func (ar *Star4Manager) findNumbers(numbers string) Star4s {
	intersection := Star4s{}

	fmt.Println(len(ar.List))
	for i, ns := range ar.List {
		if ns.Balls == numbers {
			intersection = append(intersection, ns)
			if i+1 < len(ar.List) {
				intersection = append(intersection, ar.List[i+1])
				intersection = append(intersection, *Empty())
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
