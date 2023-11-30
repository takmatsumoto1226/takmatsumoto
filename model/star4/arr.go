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

// func NewStar4s(arr []int)Star4s {

// }
