package df

import "math"

const DefaultColdOmissionThreshold = 15

// Transition 前後期號碼關係特徵 (假設 current 與 history 各期號碼皆已由小到大排序，
// history 依時間由新到舊排列，history[0] 為上一期 t-1，history[1] 為前兩期 t-2 ...)
type Transition struct {
	// 重號、鄰號與孤號
	RepeatCount      int   `json:"repeatCount"`      // 連莊數 (與上期相同號碼個數)
	RepeatNumbers    []int `json:"repeatNumbers"`    // 連莊號碼列表
	Skip1RepeatCount int   `json:"skip1RepeatCount"` // 隔期重號數 (t-2 有開、t-1 沒開、本期開出)
	NeighborCount    int   `json:"neighborCount"`    // 鄰號數 (非重號，且與上期某號差 ±1；滿足 repeat+neighbor+isolated = 球數)
	NeighborPlus     int   `json:"neighborPlus"`     // 右斜連數 (本期有號碼 = 上期某號 + 1)
	NeighborMinus    int   `json:"neighborMinus"`    // 左斜連數 (本期有號碼 = 上期某號 - 1)
	IsolatedCount    int   `json:"isolatedCount"`    // 孤號數 (與上期所有號碼距離皆 >= 2)
	MaxRepeatStreak  int   `json:"maxRepeatStreak"`  // 本期號碼最長連開期數 (含本期，>=2 表示連莊)
	MaxDiagonalRun   int   `json:"maxDiagonalRun"`   // 跨期最長斜連長度 (含本期，例如 10->11->12 為 3)

	// 幾何型態與尾數/位置連動
	SandwichCount     int   `json:"sandwichCount"`     // 夾號數 (上期開 x-1 與 x+1，本期補中洞 x)
	SplitCount        int   `json:"splitCount"`        // 分裂號數 (上期開 x，本期同時開 x-1 與 x+1)
	SameTailPrevCount int   `json:"sameTailPrevCount"` // 同尾延續數 (本期號碼個位數在上期出現過的個數)
	PositionDiffs     []int `json:"positionDiffs"`     // 逐位位移 (本期與上期排序後同位置的號碼差值)

	// 結構變化量
	SumDiff       int `json:"sumDiff"`       // 和值差 (本期 sum - 上期 sum)
	SpanDiff      int `json:"spanDiff"`      // 跨度差 (本期 span - 上期 span)
	GapZoneRefill int `json:"gapZoneRefill"` // 斷區回補數 (上期為 0 顆的十位分區，本期開出幾顆)

	// 跨多期熱冷與遺漏連動
	Recent3Count   int     `json:"recent3Count"`   // 近 3 期熱號數 (過去 3 期內曾開出的號碼數)
	Recent5Count   int     `json:"recent5Count"`   // 近 5 期熱號數 (過去 5 期內曾開出的號碼數)
	OmissionSum    int     `json:"omissionSum"`    // 開出前遺漏期數總和
	MaxOmissionHit int     `json:"maxOmissionHit"` // 本期開出號碼的最大遺漏期數
	ColdHitCount   int     `json:"coldHitCount"`   // 冷號破冰數 (開出前遺漏 >= coldThreshold 期的號碼個數)
	FollowUpScore  float64 `json:"followUpScore"`  // 歷史拖牌平均次數 (歷史上「上期號 -> 本期號」配對平均出現次數)
}

// NewTransition 計算 current 相對於 history (由新到舊) 的前後期關係特徵。
// 當 current 或 history 為空時回傳 nil。
// coldThreshold <= 0 時使用 DefaultColdOmissionThreshold (15)。
func NewTransition(current []int, history [][]int, maxNumber int, coldThreshold int) *Transition {
	if len(current) == 0 || len(history) == 0 {
		return nil
	}
	if coldThreshold <= 0 {
		coldThreshold = DefaultColdOmissionThreshold
	}

	prev := history[0]
	if len(prev) == 0 {
		return nil
	}

	currSet := make(map[int]bool, len(current))
	for _, n := range current {
		currSet[n] = true
	}

	prevSet := make(map[int]bool, len(prev))
	prevTails := [10]bool{}
	for _, n := range prev {
		prevSet[n] = true
		if n >= 0 {
			prevTails[n%10] = true
		}
	}

	var prev2Set map[int]bool
	if len(history) >= 2 {
		prev2Set = make(map[int]bool, len(history[1]))
		for _, n := range history[1] {
			prev2Set[n] = true
		}
	}

	t := &Transition{
		RepeatNumbers: []int{},
		PositionDiffs: []int{},
	}

	// 1. 重號、隔期重號、鄰號、孤號、夾號、同尾延續
	for _, n := range current {
		inPrev := prevSet[n]
		hasMinus := prevSet[n-1] // 上期有 n-1 => 本期向右斜連 (+1)
		hasPlus := prevSet[n+1]  // 上期有 n+1 => 本期向左斜連 (-1)

		if inPrev {
			t.RepeatCount++
			t.RepeatNumbers = append(t.RepeatNumbers, n)
		} else {
			if prev2Set != nil && prev2Set[n] {
				t.Skip1RepeatCount++
			}
			if hasMinus || hasPlus {
				t.NeighborCount++
			} else {
				t.IsolatedCount++
			}
		}

		if hasMinus {
			t.NeighborPlus++
		}
		if hasPlus {
			t.NeighborMinus++
		}
		if hasMinus && hasPlus {
			t.SandwichCount++
		}
		if n >= 0 && prevTails[n%10] {
			t.SameTailPrevCount++
		}
	}

	// 分裂號 (正三角)：上期開 x，本期開 x-1 與 x+1
	for _, x := range prev {
		if currSet[x-1] && currSet[x+1] {
			t.SplitCount++
		}
	}

	// 逐位位移
	if len(current) == len(prev) {
		t.PositionDiffs = make([]int, len(current))
		for i := range current {
			t.PositionDiffs[i] = current[i] - prev[i]
		}
	}

	// 2. 結構變化量：和值差、跨度差、斷區回補數
	currSum, prevSum := 0, 0
	for _, n := range current {
		currSum += n
	}
	for _, n := range prev {
		prevSum += n
	}
	t.SumDiff = currSum - prevSum
	t.SpanDiff = (current[len(current)-1] - current[0]) - (prev[len(prev)-1] - prev[0])

	tenGroups := (maxNumber + 9) / 10
	if tenGroups > 0 {
		prevTG := make([]int, tenGroups)
		currTG := make([]int, tenGroups)
		for _, n := range prev {
			if idx := (n - 1) / 10; idx >= 0 && idx < tenGroups {
				prevTG[idx]++
			}
		}
		for _, n := range current {
			if idx := (n - 1) / 10; idx >= 0 && idx < tenGroups {
				currTG[idx]++
			}
		}
		for i := 0; i < tenGroups; i++ {
			if prevTG[i] == 0 {
				t.GapZoneRefill += currTG[i]
			}
		}
	}

	// 3. 跨期連莊長度、跨期斜連長度、近 3/5 期熱號、開出前遺漏
	for _, n := range current {
		// 連莊長度 (含本期)
		streak := 1
		for _, h := range history {
			if containsInt(h, n) {
				streak++
			} else {
				break
			}
		}
		if streak > t.MaxRepeatStreak {
			t.MaxRepeatStreak = streak
		}

		// 右斜連長度 (n, n-1, n-2, ...)
		runUp := 1
		for step, h := range history {
			target := n - (step + 1)
			if target < 1 || !containsInt(h, target) {
				break
			}
			runUp++
		}
		if runUp > t.MaxDiagonalRun {
			t.MaxDiagonalRun = runUp
		}

		// 左斜連長度 (n, n+1, n+2, ...)
		runDown := 1
		for step, h := range history {
			target := n + (step + 1)
			if (maxNumber > 0 && target > maxNumber) || !containsInt(h, target) {
				break
			}
			runDown++
		}
		if runDown > t.MaxDiagonalRun {
			t.MaxDiagonalRun = runDown
		}

		// 開出前遺漏期數 & 近 3/5 期是否出現
		omission := len(history)
		found := false
		for idx, h := range history {
			if containsInt(h, n) {
				omission = idx
				found = true
				break
			}
		}
		if found && omission < 3 {
			t.Recent3Count++
		}
		if found && omission < 5 {
			t.Recent5Count++
		}
		t.OmissionSum += omission
		if omission > t.MaxOmissionHit {
			t.MaxOmissionHit = omission
		}
		if omission >= coldThreshold {
			t.ColdHitCount++
		}
	}

	// 4. 歷史拖牌平均分數：統計歷史上 prev 每個號碼出現在 t-1 時，current 每個號碼出現在 t 的總次數 / 配對數
	if len(history) >= 2 {
		pairHits := 0
		for k := 0; k < len(history)-1; k++ {
			matchCurr := 0
			for _, b := range history[k] {
				if currSet[b] {
					matchCurr++
				}
			}
			if matchCurr == 0 {
				continue
			}
			matchPrev := 0
			for _, a := range history[k+1] {
				if prevSet[a] {
					matchPrev++
				}
			}
			pairHits += matchPrev * matchCurr
		}
		rawScore := float64(pairHits) / float64(len(prev)*len(current))
		t.FollowUpScore = math.Round(rawScore*100) / 100
	}

	return t
}

func containsInt(sorted []int, target int) bool {
	lo, hi := 0, len(sorted)-1
	for lo <= hi {
		mid := (lo + hi) / 2
		if sorted[mid] == target {
			return true
		}
		if sorted[mid] < target {
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return false
}
