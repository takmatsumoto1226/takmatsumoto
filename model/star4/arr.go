package star4

import "strconv"

// Star4 ...
type Star4s []Star4

func (fa Star4s) Len() int {
	return len(fa)
}

// Less ...
func (fa Star4s) Less(i, j int) bool {
	ii, _ := strconv.Atoi(fa[i].Year + fa[i].Date)
	jj, _ := strconv.Atoi(fa[j].Year + fa[j].Date)
	return ii < jj
}

// Swap swaps the elements with indexes i and j.
func (fa Star4s) Swap(i, j int) {
	fa[i], fa[j] = fa[j], fa[i]
}

func (ar Star4s) Presentation() {
	for _, f := range ar {
		f.formRow()
	}
}

func (ar Star4s) Least() Star4 {
	return ar[len(ar)-1]
}

func (ar Star4s) First() Star4 {
	return ar[0]
}

func (ar Star4s) quickSort() Star4s {
	if len(ar) <= 1 {
		return ar
	}

	pivot := ar[len(ar)-1]
	var left, right Star4s

	for i := 0; i < len(ar)-1; i++ {
		if ar[i].Less(pivot) {
			left = append(left, ar[i])
		} else {
			right = append(right, ar[i])
		}
	}

	left = left.quickSort()
	right = right.quickSort()

	return append(append(left, pivot), right...)
}

func (ar Star4s) Statics() StaticInfos {
	infos := StaticInfos{}
	return infos
}

type StaticInfos []StaticInfo

// func NewStar4s(arr []int)Star4s {

// }

// func (ar Star4s) SortByNumber() Star4s {
// 	sortedList := Star4s{}
// }
