package ftnapi

import (
	"lottery/model/df"
	"lottery/model/ftn"
	"strconv"
)

const (
	maxNumber = 39
	pickCount = 5
	bigNumber = 20 // 大號 >= 20
)

// Draw 一期開獎
type Draw struct {
	Period    string   `json:"period"`    // 總期數, 例 "5990"
	Date      string   `json:"date"`      // YYYY-MM-DD
	YearIndex int      `json:"yearIndex"` // 當年第幾期
	Numbers   []int    `json:"numbers"`   // 由小到大
	Features  Features `json:"features"`
}

// Features 號碼組合特徵
type Features struct {
	Sum            int    `json:"sum"`
	Odd            int    `json:"odd"`
	Even           int    `json:"even"`
	Big            int    `json:"big"`   // >= 20
	Small          int    `json:"small"` // < 20
	Prime          int    `json:"prime"`
	TenGroup       [4]int `json:"tenGroup"` // 1-10, 11-20, 21-30, 31-39
	FullTenGroup   bool   `json:"fullTenGroup"`
	Consecutive    string `json:"consecutive"` // none, 2, 3, 4, 5, 22, 32
	AC             int    `json:"ac"`
	Span           int    `json:"span"`
	Mod3           [3]int `json:"mod3"` // 除以3 餘0, 餘1, 餘2 個數
	TailSum        int    `json:"tailSum"`
	SameTailGroups int    `json:"sameTailGroups"`
	MaxGap         int    `json:"maxGap"`
	MinGap         int    `json:"minGap"`
}

func newFeatures(f *df.Feature) Features {
	big := f.CountAtLeast(bigNumber)
	return Features{
		Sum:            f.SUM,
		Odd:            f.OddNumberCount,
		Even:           f.EvenNumberCount,
		Big:            big,
		Small:          len(f.IBalls) - big,
		Prime:          f.PrimeCount,
		TenGroup:       [4]int{f.TenGroupCount[0], f.TenGroupCount[1], f.TenGroupCount[2], f.TenGroupCount[3]},
		FullTenGroup:   f.IsFullTenGroup(),
		Consecutive:    df.ContinueRowTypeName(f.ContinueRowType),
		AC:             f.AC(),
		Span:           f.Span(),
		Mod3:           f.Mod3Counts(),
		TailSum:        f.TailSum(),
		SameTailGroups: f.SameTailGroups(),
		MaxGap:         f.MaxGap(),
		MinGap:         f.MinGap(),
	}
}

func newDraw(d *ftn.FTN) Draw {
	yi, _ := strconv.Atoi(d.LIdx)
	date := d.Year + "-"
	if len(d.MonthDay) == 4 {
		date += d.MonthDay[:2] + "-" + d.MonthDay[2:]
	}
	return Draw{
		Period:    d.TIdx,
		Date:      date,
		YearIndex: yi,
		Numbers:   append([]int(nil), d.IBalls...),
		Features:  newFeatures(&d.Feature),
	}
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
	Window         int          `json:"window"`
	TotalDraws     int          `json:"totalDraws"`
	ExpectedWindow float64      `json:"expectedWindow"` // 隨機下 window 期內每個號碼的期望次數
	LatestPeriod   string       `json:"latestPeriod"`
	Numbers        []NumberStat `json:"numbers"`
}

// AnalyzeRequest 分析一組號碼
type AnalyzeRequest struct {
	Numbers []int `json:"numbers" binding:"required"`
}

// HistoryMatch 歷史上與該組號碼對中 n 個的期數
type HistoryMatch struct {
	Match5 int `json:"match5"`
	Match4 int `json:"match4"`
	Match3 int `json:"match3"`
	Match2 int `json:"match2"`
}

// AnalyzeResult 分析結果
type AnalyzeResult struct {
	Numbers    []int        `json:"numbers"`
	Features   Features     `json:"features"`
	History    HistoryMatch `json:"history"`
	TotalDraws int          `json:"totalDraws"`
	LastMatch3 *Draw        `json:"lastMatch3,omitempty"` // 最近一次對中 3 個以上的期別
}

// Health 服務狀態
type Health struct {
	Status       string `json:"status"`
	Draws        int    `json:"draws"`
	LatestPeriod string `json:"latestPeriod,omitempty"`
	LatestDate   string `json:"latestDate,omitempty"`
	LoadedAt     string `json:"loadedAt,omitempty"`
}
