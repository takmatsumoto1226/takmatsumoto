package ftn

import "lottery/model/df"

func (ar *FTNsManager) findDate(date string, t int) FTNArray {
	intersection := FTNArray{}

	for i, ns := range ar.List {

		// Check elements in the second array against the set

		if ns.MonthDay == date {

			if (t == df.Before || t == df.Both) && t != df.None {
				intersection = append(intersection, ar.List[i-1])
			}

			intersection = append(intersection, ns)

			if t == df.Next || t == df.Both && t != df.None {
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
