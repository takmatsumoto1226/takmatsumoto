package ftnapi

import (
	"fmt"
	"lottery/api"
	"lottery/model/df"
	"lottery/model/ftn"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	defaultLimit  = 20
	maxLimit      = 200
	defaultWindow = 100
	queryDate     = "2006-01-02"
)

type handler struct {
	store *Store
}

// Register 註冊今彩539 路由
func Register(r gin.IRouter, s *Store) {
	h := &handler{store: s}
	r.GET("/draws", h.listDraws)
	r.GET("/draws/latest", h.latestDraw)
	r.GET("/draws/:period", h.getDraw)
	r.GET("/numbers/stats", h.numberStats)
	r.POST("/analyze", h.analyze)
}

// HealthHandler 服務狀態, 不需 api key
func HealthHandler(s *Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		draws, loadedAt := s.snapshot()
		h := Health{Status: "ok", Draws: len(draws)}
		if len(draws) == 0 {
			h.Status = "data_not_ready"
		} else {
			latest := newDraw(&draws[0])
			h.LatestPeriod = latest.Period
			h.LatestDate = latest.Date
			h.LoadedAt = loadedAt.Format(time.RFC3339)
		}
		api.OK(c, h)
	}
}

// draws 取得資料, 沒資料時直接回錯誤
func (h *handler) draws(c *gin.Context) (ftn.FTNArray, bool) {
	draws, _ := h.store.snapshot()
	if len(draws) == 0 {
		api.Fail(c, http.StatusServiceUnavailable, api.CodeDataNotReady, "draw data not loaded")
		return nil, false
	}
	return draws, true
}

func queryInt(c *gin.Context, key string, def, min, max int) (int, error) {
	v := c.Query(key)
	if v == "" {
		return def, nil
	}
	i, err := strconv.Atoi(v)
	if err != nil || i < min || i > max {
		return 0, fmt.Errorf("%s must be an integer between %d and %d", key, min, max)
	}
	return i, nil
}

// queryDateKey YYYY-MM-DD 轉成資料比對用的 YYYYMMDD
func queryDateKey(c *gin.Context, key string) (string, error) {
	v := c.Query(key)
	if v == "" {
		return "", nil
	}
	t, err := time.Parse(queryDate, v)
	if err != nil {
		return "", fmt.Errorf("%s must be YYYY-MM-DD", key)
	}
	return t.Format("20060102"), nil
}

// GET /draws?limit=20&offset=0&from=2026-01-01&to=2026-12-31
func (h *handler) listDraws(c *gin.Context) {
	limit, err := queryInt(c, "limit", defaultLimit, 1, maxLimit)
	if err != nil {
		api.Fail(c, http.StatusBadRequest, api.CodeInvalidParam, err.Error())
		return
	}
	offset, err := queryInt(c, "offset", 0, 0, 1<<30)
	if err != nil {
		api.Fail(c, http.StatusBadRequest, api.CodeInvalidParam, err.Error())
		return
	}
	from, err := queryDateKey(c, "from")
	if err != nil {
		api.Fail(c, http.StatusBadRequest, api.CodeInvalidParam, err.Error())
		return
	}
	to, err := queryDateKey(c, "to")
	if err != nil {
		api.Fail(c, http.StatusBadRequest, api.CodeInvalidParam, err.Error())
		return
	}

	draws, ok := h.draws(c)
	if !ok {
		return
	}

	matched := draws
	if from != "" || to != "" {
		matched = ftn.FTNArray{}
		for _, d := range draws {
			date := d.Date()
			if (from == "" || date >= from) && (to == "" || date <= to) {
				matched = append(matched, d)
			}
		}
	}

	page := DrawPage{Items: []Draw{}, Total: len(matched), Limit: limit, Offset: offset}
	for i := offset; i < len(matched) && i < offset+limit; i++ {
		page.Items = append(page.Items, newDraw(&matched[i]))
	}
	api.OK(c, page)
}

// GET /draws/latest
func (h *handler) latestDraw(c *gin.Context) {
	draws, ok := h.draws(c)
	if !ok {
		return
	}
	api.OK(c, newDraw(&draws[0]))
}

// GET /draws/:period
func (h *handler) getDraw(c *gin.Context) {
	if _, ok := h.draws(c); !ok {
		return
	}
	d, ok := h.store.findPeriod(c.Param("period"))
	if !ok {
		api.Fail(c, http.StatusNotFound, api.CodeNotFound, "period not found")
		return
	}
	api.OK(c, newDraw(&d))
}

// GET /numbers/stats?window=100
func (h *handler) numberStats(c *gin.Context) {
	draws, ok := h.draws(c)
	if !ok {
		return
	}
	window, err := queryInt(c, "window", defaultWindow, 1, len(draws))
	if err != nil {
		api.Fail(c, http.StatusBadRequest, api.CodeInvalidParam, err.Error())
		return
	}

	stats := make([]NumberStat, maxNumber+1)
	running := make([]int, maxNumber+1)
	for n := 1; n <= maxNumber; n++ {
		stats[n].Number = n
	}

	// draws 新到舊, 倒著走才是時間順序
	for i := len(draws) - 1; i >= 0; i-- {
		hit := make([]bool, maxNumber+1)
		for _, b := range draws[i].IBalls {
			if b >= 1 && b <= maxNumber {
				hit[b] = true
			}
		}
		for n := 1; n <= maxNumber; n++ {
			if hit[n] {
				stats[n].TotalCount++
				if i < window {
					stats[n].WindowCount++
				}
				running[n] = 0
			} else {
				running[n]++
				if running[n] > stats[n].MaxOmission {
					stats[n].MaxOmission = running[n]
				}
			}
		}
	}
	for n := 1; n <= maxNumber; n++ {
		stats[n].Omission = running[n]
	}

	api.OK(c, NumberStats{
		Window:         window,
		TotalDraws:     len(draws),
		ExpectedWindow: float64(window*pickCount) / maxNumber,
		LatestPeriod:   draws[0].TIdx,
		Numbers:        stats[1:],
	})
}

func validateNumbers(numbers []int) ([]int, error) {
	if len(numbers) != pickCount {
		return nil, fmt.Errorf("numbers must contain exactly %d numbers", pickCount)
	}
	seen := map[int]bool{}
	for _, n := range numbers {
		if n < 1 || n > maxNumber {
			return nil, fmt.Errorf("numbers must be between 1 and %d", maxNumber)
		}
		if seen[n] {
			return nil, fmt.Errorf("numbers must not repeat")
		}
		seen[n] = true
	}
	sorted := append([]int(nil), numbers...)
	sort.Ints(sorted)
	return sorted, nil
}

// POST /analyze {"numbers":[3,7,15,17,30]}
func (h *handler) analyze(c *gin.Context) {
	var req AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.Fail(c, http.StatusBadRequest, api.CodeInvalidParam, "body must be {\"numbers\":[5 integers]}")
		return
	}
	numbers, err := validateNumbers(req.Numbers)
	if err != nil {
		api.Fail(c, http.StatusBadRequest, api.CodeInvalidParam, err.Error())
		return
	}
	draws, ok := h.draws(c)
	if !ok {
		return
	}

	pick := map[int]bool{}
	for _, n := range numbers {
		pick[n] = true
	}

	result := AnalyzeResult{
		Numbers:    numbers,
		Features:   newFeatures(df.NewFeature(numbers, maxNumber)),
		TotalDraws: len(draws),
	}
	for i := range draws {
		m := 0
		for _, b := range draws[i].IBalls {
			if pick[b] {
				m++
			}
		}
		switch m {
		case 5:
			result.History.Match5++
		case 4:
			result.History.Match4++
		case 3:
			result.History.Match3++
		case 2:
			result.History.Match2++
		}
		if m >= 3 && result.LastMatch3 == nil {
			d := newDraw(&draws[i])
			result.LastMatch3 = &d
		}
	}
	api.OK(c, result)
}
