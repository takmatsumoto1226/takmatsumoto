package ftn

import (
	"fmt"
	"lottery/algorithm"
	"lottery/model/df"
	"strconv"
)

// FTNArray ...
type FTNArray []FTN

func NewFTNArray(numberss [][]string) *FTNArray {
	arr := FTNArray{}
	for _, numbers := range numberss {
		ftn := NewFTNWithStrings(numbers)
		arr = append(arr, *ftn)
	}
	return &arr
}

func NewFTNArrayWithInts(numberss [][]int) *FTNArray {
	arr := FTNArray{}
	for _, numbers := range numberss {
		ftn := NewFTNWithInts(numbers)
		arr = append(arr, *ftn)
	}
	return &arr
}

func (fa *FTNArray) Head() {
	rowmsg := "====|====|"
	for i := 1; i <= ballsCountFTN; i++ {
		rowmsg = rowmsg + fmt.Sprintf("%02d|", i)
	}
	fmt.Println(rowmsg)
	fmt.Println("")
}

func (fa FTNArray) Len() int {
	return len(fa)
}

// Less ...
func (fa FTNArray) Less(i, j int) bool {
	ii, _ := strconv.Atoi(fa[i].TIdx)
	jj, _ := strconv.Atoi(fa[j].TIdx)
	return ii > jj
}

// Swap swaps the elements with indexes i and j.
func (fa FTNArray) Swap(i, j int) {
	fa[i], fa[j] = fa[j], fa[i]
}
func (fa FTNArray) Presentation() {
	fa.PresentationWithRange(0)
}

func (fa FTNArray) PresentationWithRange(r int) {
	tmp := fa
	al := len(fa)
	if r > 0 {
		tmp = fa[al-r : al]
	}
	for _, ftn := range tmp {
		ftn.formRow()
	}
}

func (list FTNArray) findNumbers(numbers []string, t int) FTNArray {
	intersection := FTNArray{}
	set := make(map[string]bool)

	for i, ns := range list {
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
				intersection = append(intersection, list[i-1])
			}

			intersection = append(intersection, ns)

			if t == df.Next || t == df.Both {
				if i+1 < len(list) {
					intersection = append(intersection, list[i+1])
				}
			}
			if t != df.None {
				intersection = append(intersection, *Empty())
			}

		}

	}
	// Create a set from the first array

	return intersection
}

func (fa FTNArray) adariPrice(adari *FTN) {
	for i := 5; i > 1; i-- {
		combinations := algorithm.Combinations(adari.toStringArray(), i)
		// fmt.Println(combinations)
		ftnarr := NewFTNArray(combinations)
		for _, ftn := range *ftnarr {
			for _, fav := range fa {
				mc := fav.matchCount(ftn)
				if mc == 5 {
					fmt.Println("8000000元")
				} else if mc == 4 {
					fmt.Println("20000元")
				} else if mc == 3 {
					fmt.Println("300元")
				} else {
					fmt.Println("50元")
				}
			}
		}
	}
}
