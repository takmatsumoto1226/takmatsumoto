package bl

import "strconv"

type BigLotteryList []BigLottery

func (fa BigLotteryList) Len() int {
	return len(fa)
}

func (fa BigLotteryList) Presentation() {
	for _, bl := range fa {
		bl.formRow()
	}
}

// Less ...
func (fa BigLotteryList) Less(i, j int) bool {
	ii, _ := strconv.Atoi(fa[i].TIdx)
	jj, _ := strconv.Atoi(fa[j].TIdx)
	return ii > jj
}

// Swap swaps the elements with indexes i and j.
func (fa BigLotteryList) Swap(i, j int) {
	fa[i], fa[j] = fa[j], fa[i]
}
