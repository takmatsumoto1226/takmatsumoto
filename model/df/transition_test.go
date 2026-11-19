package df

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTransitionEmpty(t *testing.T) {
	a := assert.New(t)
	a.Nil(NewTransition(nil, [][]int{{1, 2, 3, 4, 5}}, 39, 15))
	a.Nil(NewTransition([]int{1, 2, 3, 4, 5}, nil, 39, 15))
}

func TestNewTransitionCoreMetrics(t *testing.T) {
	a := assert.New(t)

	// history 由新到舊：
	// t-1: [4, 14, 16, 25, 35]  (十位分區: [1, 2, 1, 1], sum=94, span=31)
	// t-2: [3, 10, 16, 20, 30]
	// t-3: [2, 10, 16, 28, 39]
	history := [][]int{
		{4, 14, 16, 25, 35},
		{3, 10, 16, 20, 30},
		{2, 10, 16, 28, 39},
	}

	// current: [5, 15, 16, 20, 34]
	// - 5 : 右斜連 (t-3=2 -> t-2=3 -> t-1=4 -> t=5, 斜四連!), 鄰號
	// - 15: 夾號 (t-1 有 14 與 16), 同時左/右斜連, 鄰號, 同尾 5 (t-1 有 25, 35)
	// - 16: 重號 (t-1, t-2, t-3, t 連續 4 期連莊!), 同尾 6 (t-1 有 16)
	// - 20: 隔期重號 (t-2 有 20, t-1 沒 20), 孤號 (與 t-1 距離 >= 4)
	// - 34: 左斜連 (t-1 有 35), 鄰號, 同尾 4 (t-1 有 4, 14)
	// 驗證 repeat + neighbor + isolated = 1 + 3 + 1 = 5
	current := []int{5, 15, 16, 20, 34}

	tr := NewTransition(current, history, 39, 3)
	a.NotNil(tr)

	a.Equal(1, tr.RepeatCount)
	a.Equal([]int{16}, tr.RepeatNumbers)
	a.Equal(1, tr.Skip1RepeatCount)
	a.Equal(3, tr.NeighborCount)
	a.Equal(2, tr.NeighborPlus)  // 5 (from 4), 15 (from 14)
	a.Equal(2, tr.NeighborMinus) // 15 (from 16), 34 (from 35)
	a.Equal(1, tr.IsolatedCount) // 20
	a.Equal(len(current), tr.RepeatCount+tr.NeighborCount+tr.IsolatedCount)
	a.Equal(4, tr.MaxRepeatStreak) // 16 連開 4 期
	a.Equal(4, tr.MaxDiagonalRun)  // 2->3->4->5 斜四連

	a.Equal(1, tr.SandwichCount)     // 15 被 14, 16 夾
	a.Equal(0, tr.SplitCount)        // 上期無 x 滿足本期有 x-1 與 x+1 (t-1=14,16 -> 本期無 13,17)
	a.Equal(4, tr.SameTailPrevCount) // 5(5), 15(5), 16(6), 34(4)
	a.Equal([]int{1, 1, 0, -5, -1}, tr.PositionDiffs)

	// sum: curr=90, prev=94 => -4; span: curr=29, prev=31 => -2
	a.Equal(-4, tr.SumDiff)
	a.Equal(-2, tr.SpanDiff)
	a.Equal(0, tr.GapZoneRefill) // 上期四個區間都有號

	// 遺漏: 5(3), 15(3), 16(0), 20(1), 34(3) => sum=10, max=3, cold(>=3)=3
	a.Equal(2, tr.Recent3Count) // 16(0), 20(1)
	a.Equal(2, tr.Recent5Count)
	a.Equal(10, tr.OmissionSum)
	a.Equal(3, tr.MaxOmissionHit)
	a.Equal(3, tr.ColdHitCount)
}

func TestNewTransitionSplitAndGapZoneRefill(t *testing.T) {
	a := assert.New(t)

	// 上期: [1, 2, 15, 18, 19] => 十位分區 [2, 3, 0, 0] (21-30 與 31-39 斷區)
	// 本期: [14, 16, 22, 33, 38] => 15 分裂成 14, 16 (SplitCount=1); 21-30 回補 1 顆, 31-39 回補 2 顆 => GapZoneRefill=3
	prev := []int{1, 2, 15, 18, 19}
	curr := []int{14, 16, 22, 33, 38}

	tr := NewTransition(curr, [][]int{prev}, 39, 15)
	a.NotNil(tr)
	a.Equal(1, tr.SplitCount)
	a.Equal(3, tr.GapZoneRefill)
	a.Equal(0.0, tr.FollowUpScore) // 只有 1 期歷史無法算拖牌
}

func TestNewTransitionFollowUpScore(t *testing.T) {
	a := assert.New(t)

	// 歷史兩期 (t-2 -> t-1):
	// t-2: [1, 2, 3, 4, 5]
	// t-1: [10, 11, 12, 13, 14]
	// 若 prev = t-1 ([10..14]), curr = [1..5]:
	// history[1] (t-2) 與 prev 無交集 => pairHits = 0
	// 若 history 有 t-2=[10,11,20,21,22] -> t-1=[1,2,30,31,32]:
	// prev = [10, 11, 30, 31, 32], curr = [1, 2, 3, 4, 5]
	// k=0: history[1] 有 2 個在 prev (10,11), history[0] 有 2 個在 curr (1,2) => pairHits = 4
	// FollowUpScore = 4 / 25 = 0.16
	history := [][]int{
		{1, 2, 30, 31, 32},
		{10, 11, 20, 21, 22},
	}
	// 注意 NewTransition 用 history[0] 當 prev，所以構造 history[0] = [10, 11, 30, 31, 32], history[1] = [1, 2, 10, 11, 22], history[2] = [10, 11, 20, 21, 22]
	history = [][]int{
		{10, 11, 30, 31, 32}, // prev (t-1)
		{1, 2, 20, 21, 22},   // t-2
		{10, 11, 25, 26, 27}, // t-3 -> t-2 有 (10,11) -> (1,2) = 2*2 = 4 次
	}
	curr := []int{1, 2, 3, 4, 5}
	tr := NewTransition(curr, history, 39, 15)
	a.NotNil(tr)
	a.Equal(0.16, tr.FollowUpScore)
}
