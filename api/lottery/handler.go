package lotteryapi

import (
	"fmt"
	"lottery/api"
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
	store         *Store
	game          *Game
	freeDrawLimit int
}

// Options 路由設定
type Options struct {
	FreeDrawLimit int // 免費用戶只能使用最新幾期, 0 表示不限
}

// Register 註冊單一彩種路由
func Register(r gin.IRouter, s *Store, opts Options) {
	h := &handler{store: s, game: s.Game, freeDrawLimit: opts.FreeDrawLimit}
	r.GET("", h.gameInfo)
	r.GET("/draws", h.listDraws)
	r.GET("/draws/latest", h.latestDraw)
	r.GET("/draws/:period", h.getDraw)
	r.GET("/numbers/stats", h.numberStats)
	r.POST("/analyze", h.analyze)
}

// GamesHandler 列出所有彩種規則
func GamesHandler(stores []*Store) gin.HandlerFunc {
	infos := make([]GameInfo, 0, len(stores))
	for _, s := range stores {
		infos = append(infos, newGameInfo(s.Game))
	}
	return func(c *gin.Context) {
		api.OK(c, infos)
	}
}

// HealthHandler 服務狀態, 不需 api key
func HealthHandler(stores []*Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := Health{Games: map[string]GameHealth{}}
		ready := 0
		for _, s := range stores {
			records, loadedAt := s.snapshot()
			gh := GameHealth{Draws: len(records)}
			if len(records) > 0 {
				ready++
				latest := s.Game.newDraw(&records[0])
				gh.LatestPeriod = latest.Period
				gh.LatestDate = latest.Date
				gh.LoadedAt = loadedAt.Format(time.RFC3339)
			}
			h.Games[s.Game.Code] = gh
		}
		switch ready {
		case len(stores):
			h.Status = "ok"
		case 0:
			h.Status = "data_not_ready"
		default:
			h.Status = "partial"
		}
		api.OK(c, h)
	}
}

// records 取得資料(新到舊), 免費用戶只取最新 freeDrawLimit 期; limited 表示有被截斷.
// 沒資料時直接回錯誤
func (h *handler) records(c *gin.Context) (records []Record, limited bool, ok bool) {
	records, _ = h.store.snapshot()
	if len(records) == 0 {
		api.Fail(c, http.StatusServiceUnavailable, api.CodeDataNotReady, h.game.Code+" draw data not loaded")
		return nil, false, false
	}
	if h.freeDrawLimit > 0 && !api.IsPremium(c) && len(records) > h.freeDrawLimit {
		return records[:h.freeDrawLimit], true, true
	}
	return records, false, true
}

func premiumRequired(c *gin.Context, msg string) {
	api.Fail(c, http.StatusForbidden, api.CodePremiumRequired, msg)
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

// GET /{game}
func (h *handler) gameInfo(c *gin.Context) {
	api.OK(c, newGameInfo(h.game))
}

// GET /{game}/draws?limit=20&offset=0&from=2026-01-01&to=2026-12-31
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

	records, limited, ok := h.records(c)
	if !ok {
		return
	}

	matched := records
	if from != "" || to != "" {
		matched = []Record{}
		for _, r := range records {
			date := r.Date()
			if (from == "" || date >= from) && (to == "" || date <= to) {
				matched = append(matched, r)
			}
		}
	}

	page := DrawPage{Items: []Draw{}, Total: len(matched), Limit: limit, Offset: offset, Limited: limited}
	for i := offset; i < len(matched) && i < offset+limit; i++ {
		page.Items = append(page.Items, h.game.newDraw(&matched[i]))
	}
	api.OK(c, page)
}

// GET /{game}/draws/latest
func (h *handler) latestDraw(c *gin.Context) {
	records, _, ok := h.records(c)
	if !ok {
		return
	}
	api.OK(c, h.game.newDraw(&records[0]))
}

// GET /{game}/draws/:period
func (h *handler) getDraw(c *gin.Context) {
	records, limited, ok := h.records(c)
	if !ok {
		return
	}
	r, idx, ok := h.store.findPeriod(c.Param("period"))
	if !ok {
		api.Fail(c, http.StatusNotFound, api.CodeNotFound, "period not found")
		return
	}
	if limited && idx >= len(records) {
		premiumRequired(c, fmt.Sprintf("free users can only view the latest %d draws", len(records)))
		return
	}
	api.OK(c, h.game.newDraw(&r))
}

// numberStats 統計 1..max 每個號碼, pick 取出每期要統計的號碼
func numberStats(records []Record, max, window int, pick func(r *Record) []int) []NumberStat {
	stats := make([]NumberStat, max+1)
	running := make([]int, max+1)
	for n := 1; n <= max; n++ {
		stats[n].Number = n
	}

	// records 新到舊, 倒著走才是時間順序
	for i := len(records) - 1; i >= 0; i-- {
		hit := make([]bool, max+1)
		for _, b := range pick(&records[i]) {
			if b >= 1 && b <= max {
				hit[b] = true
			}
		}
		for n := 1; n <= max; n++ {
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
	for n := 1; n <= max; n++ {
		stats[n].Omission = running[n]
	}
	return stats[1:]
}

// GET /{game}/numbers/stats?window=100
func (h *handler) numberStats(c *gin.Context) {
	records, limited, ok := h.records(c)
	if !ok {
		return
	}
	def := defaultWindow
	if def > len(records) {
		def = len(records)
	}
	full, _ := h.store.snapshot()
	window, err := queryInt(c, "window", def, 1, len(full))
	if err != nil {
		api.Fail(c, http.StatusBadRequest, api.CodeInvalidParam, err.Error())
		return
	}
	if window > len(records) {
		premiumRequired(c, fmt.Sprintf("free users can only use window up to %d", len(records)))
		return
	}

	g := h.game
	result := NumberStats{
		Window:         window,
		TotalDraws:     len(records),
		Limited:        limited,
		LatestPeriod:   records[0].Period,
		ExpectedWindow: float64(window*g.PickCount) / float64(g.MaxNumber),
		Numbers:        numberStats(records, g.MaxNumber, window, func(r *Record) []int { return r.Numbers }),
	}
	if g.Special != nil {
		result.SpecialExpectedWindow = float64(window) / float64(g.Special.MaxNumber)
		result.Special = numberStats(records, g.Special.MaxNumber, window, func(r *Record) []int { return []int{r.Special} })
	}
	api.OK(c, result)
}

func (g *Game) validNumbers(numbers []int) error {
	if len(numbers) != g.PickCount {
		return fmt.Errorf("numbers must contain exactly %d numbers", g.PickCount)
	}
	seen := map[int]bool{}
	for _, n := range numbers {
		if n < 1 || n > g.MaxNumber {
			return fmt.Errorf("numbers must be between 1 and %d", g.MaxNumber)
		}
		if seen[n] {
			return fmt.Errorf("numbers must not repeat")
		}
		seen[n] = true
	}
	return nil
}

func (g *Game) validSpecialPick(special *int) error {
	if g.Special == nil || !g.Special.UserPick {
		if special != nil {
			return fmt.Errorf("special is not allowed for %s", g.Code)
		}
		return nil
	}
	if special == nil {
		return fmt.Errorf("special is required for %s", g.Code)
	}
	if *special < 1 || *special > g.Special.MaxNumber {
		return fmt.Errorf("special must be between 1 and %d", g.Special.MaxNumber)
	}
	return nil
}

// POST /{game}/analyze {"numbers":[...], "special": n}
func (h *handler) analyze(c *gin.Context) {
	g := h.game
	var req AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.Fail(c, http.StatusBadRequest, api.CodeInvalidParam, fmt.Sprintf("body must be {\"numbers\":[%d integers]}", g.PickCount))
		return
	}
	if err := g.validNumbers(req.Numbers); err != nil {
		api.Fail(c, http.StatusBadRequest, api.CodeInvalidParam, err.Error())
		return
	}
	if err := g.validSpecialPick(req.Special); err != nil {
		api.Fail(c, http.StatusBadRequest, api.CodeInvalidParam, err.Error())
		return
	}
	records, limited, ok := h.records(c)
	if !ok {
		return
	}

	numbers := append([]int(nil), req.Numbers...)
	sort.Ints(numbers)
	pick := map[int]bool{}
	for _, n := range numbers {
		pick[n] = true
	}

	type matchKey struct {
		main    int
		special bool
	}
	counts := map[matchKey]int{}
	result := AnalyzeResult{
		Numbers:    numbers,
		Special:    req.Special,
		Features:   g.newFeatures(numbers),
		TotalDraws: len(records),
		Limited:    limited,
	}
	for i := range records {
		r := &records[i]
		m := 0
		for _, b := range r.Numbers {
			if pick[b] {
				m++
			}
		}
		sp := false
		if g.Special != nil {
			if g.Special.UserPick {
				sp = r.Special == *req.Special
			} else {
				sp = pick[r.Special]
			}
		}
		counts[matchKey{m, sp}]++
		if m >= 3 && result.LastMatch3 == nil {
			d := g.newDraw(r)
			result.LastMatch3 = &d
		}
	}

	specials := []bool{false}
	if g.Special != nil {
		specials = []bool{true, false}
	}
	result.History = []MatchCount{}
	for m := g.PickCount; m >= 1; m-- {
		for _, sp := range specials {
			// 同一號碼池的特別號不可能在主號全中時又對中
			if sp && g.Special.SamePool && m == g.PickCount {
				continue
			}
			k := matchKey{m, sp}
			result.History = append(result.History, MatchCount{Main: m, Special: sp, Count: counts[k], Prize: g.PrizeName(m, sp)})
		}
	}
	api.OK(c, result)
}
