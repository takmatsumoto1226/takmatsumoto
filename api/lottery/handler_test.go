package lotteryapi

import (
	"bytes"
	"encoding/json"
	"lottery/api"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func rec(period, date string, special int, numbers ...int) Record {
	return Record{Period: period, Year: date[:4], MonthDay: date[4:], YearIndex: 1, Numbers: numbers, Special: special}
}

func newTestRouter(key string) (*gin.Engine, []*Store) {
	gin.SetMode(gin.TestMode)
	stores := []*Store{
		NewStoreWithRecords(FTN(), []Record{
			rec("5988", "20261008", 0, 3, 7, 15, 17, 30),
			rec("5989", "20261009", 0, 20, 23, 25, 27, 31),
			rec("5990", "20261010", 0, 14, 21, 23, 32, 34),
		}),
		NewStoreWithRecords(BL(), []Record{
			rec("2488", "20261006", 23, 4, 6, 20, 32, 36, 39),
			rec("2489", "20261009", 27, 3, 9, 11, 25, 34, 41),
			rec("47", "20190201", 47, 3, 5, 7, 17, 28, 30), // 期數重複的錯誤資料
			rec("47", "20190204", 47, 2, 3, 13, 15, 19, 26),
			rec("1291", "20160315", 18, 6, 14, 30, 42, 41, 48), // 未排序
		}),
		NewStoreWithRecords(PW(), []Record{
			rec("1952", "20261005", 6, 2, 6, 17, 21, 26, 37),
			rec("1953", "20261008", 5, 4, 23, 27, 28, 29, 38),
			rec("1849", "20251009", 28, 1, 6, 11, 20, 34, 35), // 第二區超出範圍
		}),
	}

	e := gin.New()
	v1 := e.Group("/api/v1")
	v1.GET("/health", HealthHandler(stores))
	authed := v1.Group("", api.APIKeyAuth(key))
	authed.GET("/games", GamesHandler(stores))
	for _, s := range stores {
		Register(authed.Group("/"+s.Game.Code), s)
	}
	return e, stores
}

type envelope struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func do(e *gin.Engine, method, path, body string, header map[string]string) (int, envelope) {
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range header {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	var env envelope
	json.Unmarshal(w.Body.Bytes(), &env)
	return w.Code, env
}

func TestClean(t *testing.T) {
	a := assert.New(t)
	_, stores := newTestRouter("")

	bl, _ := stores[1].snapshot()
	a.Len(bl, 3)
	a.Equal("2489", bl[0].Period)
	a.Equal("1291", bl[2].Period)
	a.Equal([]int{6, 14, 30, 41, 42, 48}, bl[2].Numbers)

	pw, _ := stores[2].snapshot()
	a.Len(pw, 2)
}

func TestDedupPeriods(t *testing.T) {
	a := assert.New(t)
	s := NewStoreWithRecords(FTN(), []Record{
		rec("3790", "20191008", 0, 1, 2, 3, 4, 5),
		rec("3791", "20160729", 0, 3, 27, 32, 35, 39), // 期數抄錯
		rec("3791", "20191009", 0, 1, 8, 19, 20, 36),
		rec("3792", "20191010", 0, 6, 7, 8, 9, 10),
		rec("3793", "20191011", 0, 6, 7, 8, 9, 11),
		rec("3793", "20191011", 0, 6, 7, 8, 9, 11), // 完全相同
		rec("36", "20040506", 0, 1, 4, 10, 23, 25),
		rec("037", "20040510", 0, 3, 4, 10, 23, 25), // 前導零
		rec("38", "20040513", 0, 2, 4, 10, 23, 25),
		rec("37", "20190207", 0, 1, 5, 11, 21, 22), // 春節加碼混入, 日期不合
	})
	records, _ := s.snapshot()
	a.Len(records, 7)
	r37, ok := s.findPeriod("37")
	a.True(ok)
	a.Equal("20040510", r37.Date())
	r, ok := s.findPeriod("3791")
	a.True(ok)
	a.Equal("20191009", r.Date())
}

func TestGames(t *testing.T) {
	a := assert.New(t)
	e, _ := newTestRouter("")

	_, env := do(e, http.MethodGet, "/api/v1/games", "", nil)
	var games []GameInfo
	a.NoError(json.Unmarshal(env.Data, &games))
	a.Len(games, 3)
	a.Equal("ftn", games[0].Code)
	a.Nil(games[0].Special)
	a.True(games[1].Special.SamePool)
	a.Equal(8, games[2].Special.MaxNumber)
	a.True(games[2].Special.UserPick)

	_, env = do(e, http.MethodGet, "/api/v1/bl", "", nil)
	var g GameInfo
	a.NoError(json.Unmarshal(env.Data, &g))
	a.Equal("大樂透", g.Name)
}

func TestListDraws(t *testing.T) {
	a := assert.New(t)
	e, _ := newTestRouter("")

	status, env := do(e, http.MethodGet, "/api/v1/ftn/draws?limit=2", "", nil)
	a.Equal(http.StatusOK, status)
	var page DrawPage
	a.NoError(json.Unmarshal(env.Data, &page))
	a.Equal(3, page.Total)
	a.Len(page.Items, 2)
	a.Equal("5990", page.Items[0].Period)
	a.Equal("2026-10-10", page.Items[0].Date)
	a.Equal([]int{14, 21, 23, 32, 34}, page.Items[0].Numbers)
	a.Nil(page.Items[0].Special)
	a.Equal(124, page.Items[0].Features.Sum)
	a.Len(page.Items[0].Features.TenGroup, 4)

	_, env = do(e, http.MethodGet, "/api/v1/ftn/draws?from=2026-10-09&to=2026-10-09", "", nil)
	a.NoError(json.Unmarshal(env.Data, &page))
	a.Equal(1, page.Total)
	a.Equal("5989", page.Items[0].Period)

	_, env = do(e, http.MethodGet, "/api/v1/bl/draws", "", nil)
	a.NoError(json.Unmarshal(env.Data, &page))
	a.Equal(27, *page.Items[0].Special)
	a.Len(page.Items[0].Features.TenGroup, 5)
	a.Equal(3, page.Items[0].Features.Big) // >= 25 : 25, 34, 41

	status, env = do(e, http.MethodGet, "/api/v1/ftn/draws?limit=999", "", nil)
	a.Equal(http.StatusBadRequest, status)
	a.Equal(api.CodeInvalidParam, env.Code)

	status, _ = do(e, http.MethodGet, "/api/v1/ftn/draws?from=20261009", "", nil)
	a.Equal(http.StatusBadRequest, status)
}

func TestGetDraw(t *testing.T) {
	a := assert.New(t)
	e, _ := newTestRouter("")

	_, env := do(e, http.MethodGet, "/api/v1/pw/draws/latest", "", nil)
	var d Draw
	a.NoError(json.Unmarshal(env.Data, &d))
	a.Equal("1953", d.Period)
	a.Equal(5, *d.Special)
	a.Equal("3", d.Features.Consecutive) // 27 28 29

	_, env = do(e, http.MethodGet, "/api/v1/ftn/draws/5988", "", nil)
	a.NoError(json.Unmarshal(env.Data, &d))
	a.Equal([]int{3, 7, 15, 17, 30}, d.Numbers)
	a.Equal("none", d.Features.Consecutive)

	status, env := do(e, http.MethodGet, "/api/v1/bl/draws/47", "", nil)
	a.Equal(http.StatusNotFound, status)
	a.Equal(api.CodeNotFound, env.Code)
}

func TestNumberStats(t *testing.T) {
	a := assert.New(t)
	e, _ := newTestRouter("")

	_, env := do(e, http.MethodGet, "/api/v1/ftn/numbers/stats?window=2", "", nil)
	var s NumberStats
	a.NoError(json.Unmarshal(env.Data, &s))
	a.Len(s.Numbers, 39)
	a.Nil(s.Special)
	n23 := s.Numbers[22]
	a.Equal(23, n23.Number)
	a.Equal(2, n23.WindowCount)
	a.Equal(0, n23.Omission)
	n3 := s.Numbers[2]
	a.Equal(0, n3.WindowCount)
	a.Equal(1, n3.TotalCount)
	a.Equal(2, n3.Omission)
	a.Equal(3, s.Numbers[38].Omission)

	_, env = do(e, http.MethodGet, "/api/v1/pw/numbers/stats?window=1", "", nil)
	a.NoError(json.Unmarshal(env.Data, &s))
	a.Len(s.Numbers, 38)
	a.Len(s.Special, 8)
	a.Equal(1, s.Special[4].WindowCount) // 第二區 5
	a.Equal(1, s.Special[5].Omission)    // 第二區 6 上一期開出
	a.InDelta(0.125, s.SpecialExpectedWindow, 1e-9)

	status, _ := do(e, http.MethodGet, "/api/v1/ftn/numbers/stats?window=4", "", nil)
	a.Equal(http.StatusBadRequest, status)
}

func findMatch(h []MatchCount, main int, special bool) MatchCount {
	for _, m := range h {
		if m.Main == main && m.Special == special {
			return m
		}
	}
	return MatchCount{Main: -1}
}

func TestAnalyzeFTN(t *testing.T) {
	a := assert.New(t)
	e, _ := newTestRouter("")

	_, env := do(e, http.MethodPost, "/api/v1/ftn/analyze", `{"numbers":[34,14,23,21,1]}`, nil)
	var r AnalyzeResult
	a.NoError(json.Unmarshal(env.Data, &r))
	a.Equal([]int{1, 14, 21, 23, 34}, r.Numbers)
	a.Len(r.History, 5)
	m4 := findMatch(r.History, 4, false)
	a.Equal(1, m4.Count)
	a.Equal("貳獎", m4.Prize)
	a.Equal(1, findMatch(r.History, 1, false).Count) // 5989 只中 23
	a.Equal("", findMatch(r.History, 1, false).Prize)
	a.Equal("5990", r.LastMatch3.Period)

	for _, body := range []string{`{"numbers":[1,2,3,4]}`, `{"numbers":[1,2,3,4,40]}`, `{"numbers":[1,1,3,4,5]}`, `{}`, `{"numbers":[1,2,3,4,5],"special":1}`} {
		status, env := do(e, http.MethodPost, "/api/v1/ftn/analyze", body, nil)
		a.Equal(http.StatusBadRequest, status, body)
		a.Equal(api.CodeInvalidParam, env.Code, body)
	}
}

func TestAnalyzeBL(t *testing.T) {
	a := assert.New(t)
	e, _ := newTestRouter("")

	// 2489 : 03 09 11 25 34 41 特別號 27
	_, env := do(e, http.MethodPost, "/api/v1/bl/analyze", `{"numbers":[3,9,11,25,34,27]}`, nil)
	var r AnalyzeResult
	a.NoError(json.Unmarshal(env.Data, &r))
	m := findMatch(r.History, 5, true)
	a.Equal(1, m.Count)
	a.Equal("貳獎", m.Prize)
	a.Equal(-1, findMatch(r.History, 6, true).Main) // 不可能的組合不列出
	a.Len(r.History, 11)

	status, env := do(e, http.MethodPost, "/api/v1/bl/analyze", `{"numbers":[3,9,11,25,34,27],"special":5}`, nil)
	a.Equal(http.StatusBadRequest, status)
	a.Equal("special is not allowed for bl", env.Message)
}

func TestAnalyzePW(t *testing.T) {
	a := assert.New(t)
	e, _ := newTestRouter("")

	// 1953 : 04 23 27 28 29 38 第二區 05
	_, env := do(e, http.MethodPost, "/api/v1/pw/analyze", `{"numbers":[4,23,27,28,29,38],"special":5}`, nil)
	var r AnalyzeResult
	a.NoError(json.Unmarshal(env.Data, &r))
	a.Equal(5, *r.Special)
	top := findMatch(r.History, 6, true)
	a.Equal(1, top.Count)
	a.Equal("頭獎", top.Prize)
	a.Equal("普獎", findMatch(r.History, 1, true).Prize)
	a.Len(r.History, 12)

	for _, body := range []string{`{"numbers":[4,23,27,28,29,38]}`, `{"numbers":[4,23,27,28,29,38],"special":9}`, `{"numbers":[4,23,27,28,29,39],"special":1}`} {
		status, env := do(e, http.MethodPost, "/api/v1/pw/analyze", body, nil)
		a.Equal(http.StatusBadRequest, status, body)
		a.Equal(api.CodeInvalidParam, env.Code, body)
	}
}

func TestAPIKey(t *testing.T) {
	a := assert.New(t)
	e, _ := newTestRouter("secret")

	status, env := do(e, http.MethodGet, "/api/v1/bl/draws/latest", "", nil)
	a.Equal(http.StatusUnauthorized, status)
	a.Equal(api.CodeUnauthorized, env.Code)

	status, _ = do(e, http.MethodGet, "/api/v1/bl/draws/latest", "", map[string]string{api.HeaderAPIKey: "secret"})
	a.Equal(http.StatusOK, status)

	status, env = do(e, http.MethodGet, "/api/v1/health", "", nil)
	a.Equal(http.StatusOK, status)
	var h Health
	a.NoError(json.Unmarshal(env.Data, &h))
	a.Equal("ok", h.Status)
	a.Equal("2489", h.Games["bl"].LatestPeriod)
}

func TestDataNotReady(t *testing.T) {
	a := assert.New(t)
	gin.SetMode(gin.TestMode)
	stores := []*Store{NewStore(PW()), NewStoreWithRecords(FTN(), []Record{rec("1", "20260101", 0, 1, 2, 3, 4, 5)})}
	e := gin.New()
	e.GET("/health", HealthHandler(stores))
	Register(e.Group("/pw"), stores[0])

	status, env := do(e, http.MethodGet, "/pw/draws/latest", "", nil)
	a.Equal(http.StatusServiceUnavailable, status)
	a.Equal(api.CodeDataNotReady, env.Code)

	_, env = do(e, http.MethodGet, "/health", "", nil)
	var h Health
	a.NoError(json.Unmarshal(env.Data, &h))
	a.Equal("partial", h.Status)
}
