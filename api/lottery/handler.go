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
	drawsPageSize = 1 // /draws 每次只回傳一筆, limit 超過時自動縮成此值
	defaultWindow = 100
	queryDate     = "2006-01-02"
)

type handler struct {
	store         *Store
	game          *Game
	freeDrawLimit int
	maxDrawLimit  int
}

// Options 路由設定
type Options struct {
	FreeDrawLimit int // 免費用戶只能使用最新幾期, 0 表示不限
	MaxDrawLimit  int // 所有用戶(含付費)只能使用最新幾期, 0 表示不限; 用來暫停全歷史查詢
}

// Register 註冊單一彩種路由
func Register(r gin.IRouter, s *Store, opts Options) {
	h := &handler{store: s, game: s.Game, freeDrawLimit: opts.FreeDrawLimit, maxDrawLimit: opts.MaxDrawLimit}
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

// freeCapped 目前請求的上限是否來自免費額度(付費可解除), 否則來自 maxDrawLimit
func (h *handler) freeCapped(c *gin.Context) bool {
	return h.freeDrawLimit > 0 && !api.IsPremium(c) && (h.maxDrawLimit == 0 || h.freeDrawLimit < h.maxDrawLimit)
}

// drawLimit 目前請求可使用的最新期數, 0 表示不限
func (h *handler) drawLimit(c *gin.Context) int {
	if h.freeCapped(c) {
		return h.freeDrawLimit
	}
	return h.maxDrawLimit
}

// records 取得資料(新到舊), 只取最新 drawLimit 期 (免費用戶 freeDrawLimit, 所有用戶 maxDrawLimit);
// limited 表示有被截斷. 沒資料時直接回錯誤
func (h *handler) records(c *gin.Context) (records []Record, limited bool, ok bool) {
	return h.recordsUpTo(c, h.drawLimit(c))
}

// analyzeRecords 對獎用資料: 只回統計結果不回整份歷史, 所以不受 maxDrawLimit 限制;
// 付費用戶比對全歷史, 免費用戶仍只比對最新 freeDrawLimit 期
func (h *handler) analyzeRecords(c *gin.Context) (records []Record, limited bool, ok bool) {
	n := 0
	if !api.IsPremium(c) {
		n = h.freeDrawLimit
	}
	return h.recordsUpTo(c, n)
}

// recordsUpTo 取最新 n 期(新到舊), n <= 0 表示不限; limited 表示有被截斷. 沒資料時直接回錯誤
func (h *handler) recordsUpTo(c *gin.Context, n int) (records []Record, limited bool, ok bool) {
	records, _ = h.store.snapshot()
	if len(records) == 0 {
		api.Fail(c, http.StatusServiceUnavailable, api.CodeDataNotReady, h.game.Code+" draw data not loaded")
		return nil, false, false
	}
	if n > 0 && len(records) > n {
		return records[:n], true, true
	}
	return records, false, true
}

// limitExceeded 請求超出可用期數: 免費額度回 premium required, 全站上限回 history unavailable
func (h *handler) limitExceeded(c *gin.Context, detail string) {
	if h.freeCapped(c) {
		api.Fail(c, http.StatusForbidden, api.CodePremiumRequired, "free users "+detail)
		return
	}
	api.Fail(c, http.StatusForbidden, api.CodeHistoryUnavailable, "full history is temporarily unavailable, "+detail)
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

// GET /{game}/draws?offset=0&from=2026-01-01&to=2026-12-31 (每次只回傳 drawsPageSize 筆, 用 offset 翻頁)
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

	if limit > drawsPageSize {
		limit = drawsPageSize
	}

	records, limited, ok := h.records(c)
	if !ok {
		return
	}

	matchedIdx := make([]int, 0, len(records))
	if from != "" || to != "" {
		for i, r := range records {
			date := r.Date()
			if (from == "" || date >= from) && (to == "" || date <= to) {
				matchedIdx = append(matchedIdx, i)
			}
		}
	} else {
		for i := range records {
			matchedIdx = append(matchedIdx, i)
		}
	}

	page := DrawPage{Items: []Draw{}, Total: len(matchedIdx), Limit: limit, Offset: offset, Limited: limited}
	for i := offset; i < len(matchedIdx) && i < offset+limit; i++ {
		origIdx := matchedIdx[i]
		page.Items = append(page.Items, h.game.newDrawWithHistory(&records[origIdx], records[origIdx+1:]))
	}
	api.OK(c, page)
}

// GET /{game}/draws/latest
func (h *handler) latestDraw(c *gin.Context) {
	records, _, ok := h.records(c)
	if !ok {
		return
	}
	api.OK(c, h.game.newDrawWithHistory(&records[0], records[1:]))
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
		h.limitExceeded(c, fmt.Sprintf("can only view the latest %d draws", len(records)))
		return
	}
	var prev []Record
	if idx+1 < len(records) {
		prev = records[idx+1:]
	}
	api.OK(c, h.game.newDrawWithHistory(&r, prev))
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
		h.limitExceeded(c, fmt.Sprintf("can only use window up to %d", len(records)))
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
	records, limited, ok := h.analyzeRecords(c)
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
		Transition: g.newTransition(numbers, records),
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
			d := g.newDrawWithHistory(r, records[i+1:])
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
