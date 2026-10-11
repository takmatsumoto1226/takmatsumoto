package lotteryapi

import (
	"lottery/model/df"
)

// GameInfo 彩種規則
type GameInfo struct {
	Code      string       `json:"code"`
	Name      string       `json:"name"`
	MaxNumber int          `json:"maxNumber"`
	PickCount int          `json:"pickCount"`
	BigFrom   int          `json:"bigFrom"`
	Special   *SpecialSpec `json:"special,omitempty"`
	Prizes    []Prize      `json:"prizes"`
}

func newGameInfo(g *Game) GameInfo {
	return GameInfo{
		Code:      g.Code,
		Name:      g.Name,
		MaxNumber: g.MaxNumber,
		PickCount: g.PickCount,
		BigFrom:   g.BigFrom,
		Special:   g.Special,
		Prizes:    g.Prizes,
	}
}

// Draw 一期開獎
type Draw struct {
	Period    string   `json:"period"`            // 總期數
	Date      string   `json:"date"`              // YYYY-MM-DD
	YearIndex int      `json:"yearIndex"`         // 當年第幾期
	Numbers   []int    `json:"numbers"`           // 主號, 由小到大
	Special   *int     `json:"special,omitempty"` // 大樂透特別號 / 威力彩第二區
	Features  Features `json:"features"`
}

// Features 主號組合特徵
type Features struct {
	Sum            int    `json:"sum"`
	Odd            int    `json:"odd"`
	Even           int    `json:"even"`
	Big            int    `json:"big"`   // >= bigFrom
	Small          int    `json:"small"` // < bigFrom
	Prime          int    `json:"prime"`
	TenGroup       []int  `json:"tenGroup"` // 1-10, 11-20, ...
	FullTenGroup   bool   `json:"fullTenGroup"`
	Consecutive    string `json:"consecutive"` // none, 2, 3, 22, 32, 222 ...
	AC             int    `json:"ac"`
	Span           int    `json:"span"`
	Mod3           [3]int `json:"mod3"` // 除以3 餘0, 餘1, 餘2 個數
	TailSum        int    `json:"tailSum"`
	SameTailGroups int    `json:"sameTailGroups"`
	MaxGap         int    `json:"maxGap"`
	MinGap         int    `json:"minGap"`
}

func (g *Game) newFeatures(numbers []int) Features {
	f := df.NewFeature(numbers, g.MaxNumber)
	big := f.CountAtLeast(g.BigFrom)
	tenGroup := make([]int, g.TenGroups())
	full := true
	for i := range tenGroup {
		tenGroup[i] = f.TenGroupCount[i]
		if tenGroup[i] == 0 {
			full = false
		}
	}
	return Features{
		Sum:            f.SUM,
		Odd:            f.OddNumberCount,
		Even:           f.EvenNumberCount,
		Big:            big,
		Small:          len(numbers) - big,
		Prime:          f.PrimeCount,
		TenGroup:       tenGroup,
		FullTenGroup:   full,
		Consecutive:    f.ConsecutivePattern(),
		AC:             f.AC(),
		Span:           f.Span(),
		Mod3:           f.Mod3Counts(),
		TailSum:        f.TailSum(),
		SameTailGroups: f.SameTailGroups(),
		MaxGap:         f.MaxGap(),
		MinGap:         f.MinGap(),
	}
}

func (g *Game) newDraw(r *Record) Draw {
	d := Draw{
		Period:    r.Period,
		Date:      r.Year + "-" + r.MonthDay[:2] + "-" + r.MonthDay[2:],
		YearIndex: r.YearIndex,
		Numbers:   append([]int(nil), r.Numbers...),
		Features:  g.newFeatures(r.Numbers),
	}
	if g.Special != nil {
		sp := r.Special
		d.Special = &sp
	}
	return d
}

// DrawPage 分頁結果
type DrawPage struct {
	Items  []Draw `json:"items"`
	Total  int    `json:"total"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

// NumberStat 單一號碼統計
type NumberStat struct {
	Number      int `json:"number"`
	WindowCount int `json:"windowCount"` // 最近 window 期出現次數
	TotalCount  int `json:"totalCount"`  // 全部歷史出現次數
	Omission    int `json:"omission"`    // 目前連續未開期數
	MaxOmission int `json:"maxOmission"` // 歷史最長未開期數
}

// NumberStats 號碼統計結果
type NumberStats struct {
	Window                int          `json:"window"`
	TotalDraws            int          `json:"totalDraws"`
	LatestPeriod          string       `json:"latestPeriod"`
	ExpectedWindow        float64      `json:"expectedWindow"` // 隨機下 window 期內每個主號的期望次數
	Numbers               []NumberStat `json:"numbers"`
	SpecialExpectedWindow float64      `json:"specialExpectedWindow,omitempty"`
	Special               []NumberStat `json:"special,omitempty"` // 特別號 / 第二區統計
}

// AnalyzeRequest 分析一組號碼
type AnalyzeRequest struct {
	Numbers []int `json:"numbers" binding:"required"`
	Special *int  `json:"special"` // 威力彩第二區, 其他彩種不可帶
}

// MatchCount 歷史上對中 Main 個主號(且特別號是否對中)的期數
type MatchCount struct {
	Main    int    `json:"main"`
	Special bool   `json:"special"`
	Count   int    `json:"count"`
	Prize   string `json:"prize,omitempty"` // 對應獎項, 沒有獎項時省略
}

// AnalyzeResult 分析結果
type AnalyzeResult struct {
	Numbers    []int        `json:"numbers"`
	Special    *int         `json:"special,omitempty"`
	Features   Features     `json:"features"`
	History    []MatchCount `json:"history"` // 主號對中數由多到少
	TotalDraws int          `json:"totalDraws"`
	LastMatch3 *Draw        `json:"lastMatch3,omitempty"` // 最近一次主號對中 3 個以上的期別
}

// GameHealth 單一彩種資料狀態
type GameHealth struct {
	Draws        int    `json:"draws"`
	LatestPeriod string `json:"latestPeriod,omitempty"`
	LatestDate   string `json:"latestDate,omitempty"`
	LoadedAt     string `json:"loadedAt,omitempty"`
}

// Health 服務狀態
type Health struct {
	Status string                `json:"status"` // ok, partial, data_not_ready
	Games  map[string]GameHealth `json:"games"`
}
