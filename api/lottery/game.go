package lotteryapi

import (
	"lottery/model/bl"
	"lottery/model/ftn"
	"lottery/model/pw"
	"strconv"
)

// SpecialSpec 特別號 / 第二區
type SpecialSpec struct {
	Name      string `json:"name"`      // 顯示名稱 : 特別號, 第二區
	MaxNumber int    `json:"maxNumber"` // 號碼範圍 1..MaxNumber
	SamePool  bool   `json:"samePool"`  // true : 與主號同一號碼池開出(大樂透), 玩家不另外選
	UserPick  bool   `json:"userPick"`  // true : 玩家需另外選一個號碼(威力彩)
}

// Prize 獎項 : 對中主號 Main 個, 且特別號是否對中
type Prize struct {
	Main    int    `json:"main"`
	Special bool   `json:"special"`
	Name    string `json:"name"`
}

// Record 一期開獎, 各彩種共用
type Record struct {
	Period    string
	Year      string
	MonthDay  string
	YearIndex int
	Numbers   []int // 主號, 由小到大
	Special   int   // 0 : 無
}

// Date YYYYMMDD
func (r *Record) Date() string { return r.Year + r.MonthDay }

// Game 彩種設定
type Game struct {
	Code      string
	Name      string
	MaxNumber int
	PickCount int
	BigFrom   int // >= BigFrom 算大號
	Special   *SpecialSpec
	Prizes    []Prize
	load      func() []Record
}

// TenGroups 十位分區數 : 1-10, 11-20, ...
func (g *Game) TenGroups() int {
	return (g.MaxNumber + 9) / 10
}

// PrizeName 依對中數取得獎項名稱, 沒中獎回傳空字串
func (g *Game) PrizeName(main int, special bool) string {
	for _, p := range g.Prizes {
		if p.Main == main && p.Special == special {
			return p.Name
		}
	}
	return ""
}

// Games 支援的彩種, 順序即 /games 回傳順序
func Games() []*Game {
	return []*Game{FTN(), BL(), PW()}
}

// FTN 今彩539 : 1-39 選 5
func FTN() *Game {
	return &Game{
		Code: "ftn", Name: "今彩539", MaxNumber: 39, PickCount: 5, BigFrom: 20,
		Prizes: []Prize{
			{5, false, "頭獎"}, {4, false, "貳獎"}, {3, false, "參獎"}, {2, false, "肆獎"},
		},
		load: func() []Record {
			list := ftn.LoadHistory()
			records := make([]Record, 0, len(list))
			for _, d := range list {
				records = append(records, newRecord(d.TIdx, d.Year, d.MonthDay, d.LIdx, d.IBalls, 0))
			}
			return records
		},
	}
}

// BL 大樂透 : 1-49 選 6, 另開 1 個特別號
func BL() *Game {
	return &Game{
		Code: "bl", Name: "大樂透", MaxNumber: 49, PickCount: 6, BigFrom: 25,
		Special: &SpecialSpec{Name: "特別號", MaxNumber: 49, SamePool: true},
		Prizes: []Prize{
			{6, false, "頭獎"}, {5, true, "貳獎"}, {5, false, "參獎"}, {4, true, "肆獎"},
			{4, false, "伍獎"}, {3, true, "陸獎"}, {2, true, "柒獎"}, {3, false, "普獎"},
		},
		load: func() []Record {
			list := bl.LoadHistory()
			records := make([]Record, 0, len(list))
			for _, d := range list {
				if len(d.IBalls) != 7 {
					continue
				}
				records = append(records, newRecord(d.TIdx, d.Year, d.MonthDay, d.LIdx, d.IBalls[:6], d.IBalls[6]))
			}
			return records
		},
	}
}

// PW 威力彩 : 第一區 1-38 選 6, 第二區 1-8 選 1
func PW() *Game {
	return &Game{
		Code: "pw", Name: "威力彩", MaxNumber: 38, PickCount: 6, BigFrom: 20,
		Special: &SpecialSpec{Name: "第二區", MaxNumber: 8, UserPick: true},
		Prizes: []Prize{
			{6, true, "頭獎"}, {6, false, "貳獎"}, {5, true, "參獎"}, {5, false, "肆獎"},
			{4, true, "伍獎"}, {4, false, "陸獎"}, {3, true, "柒獎"}, {2, true, "捌獎"},
			{3, false, "玖獎"}, {1, true, "普獎"},
		},
		load: func() []Record {
			list := pw.LoadHistory()
			records := make([]Record, 0, len(list))
			for _, d := range list {
				records = append(records, newRecord(d.TIdx, d.Year, d.MonthDay, d.LIdx, d.IBalls, d.S1.Digit))
			}
			return records
		},
	}
}

func newRecord(period, year, monthDay, yearIndex string, numbers []int, special int) Record {
	yi, _ := strconv.Atoi(yearIndex)
	return Record{
		Period:    period,
		Year:      year,
		MonthDay:  monthDay,
		YearIndex: yi,
		Numbers:   append([]int(nil), numbers...),
		Special:   special,
	}
}
