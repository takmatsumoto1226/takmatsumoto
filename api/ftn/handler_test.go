package ftnapi

import (
	"bytes"
	"encoding/json"
	"lottery/api"
	"lottery/model/ftn"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func newTestRouter(key string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	rows := [][]string{
		{"2026", "1008", "244", "03", "07", "15", "17", "30", "5988"},
		{"2026", "1009", "245", "20", "23", "25", "27", "31", "5989"},
		{"2026", "1010", "246", "14", "21", "23", "32", "34", "5990"},
	}
	draws := ftn.FTNArray{}
	for i := len(rows) - 1; i >= 0; i-- {
		draws = append(draws, *ftn.NewFTN(rows[i]))
	}
	store := NewStoreWithDraws(draws)

	e := gin.New()
	v1 := e.Group("/api/v1")
	v1.GET("/health", HealthHandler(store))
	Register(v1.Group("/ftn", api.APIKeyAuth(key)), store)
	return e
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

func TestListDraws(t *testing.T) {
	a := assert.New(t)
	e := newTestRouter("")

	status, env := do(e, http.MethodGet, "/api/v1/ftn/draws?limit=2", "", nil)
	a.Equal(http.StatusOK, status)
	var page DrawPage
	a.NoError(json.Unmarshal(env.Data, &page))
	a.Equal(3, page.Total)
	a.Len(page.Items, 2)
	a.Equal("5990", page.Items[0].Period)
	a.Equal("2026-10-10", page.Items[0].Date)
	a.Equal([]int{14, 21, 23, 32, 34}, page.Items[0].Numbers)
	a.Equal(124, page.Items[0].Features.Sum)

	_, env = do(e, http.MethodGet, "/api/v1/ftn/draws?from=2026-10-09&to=2026-10-09", "", nil)
	a.NoError(json.Unmarshal(env.Data, &page))
	a.Equal(1, page.Total)
	a.Equal("5989", page.Items[0].Period)

	status, env = do(e, http.MethodGet, "/api/v1/ftn/draws?limit=999", "", nil)
	a.Equal(http.StatusBadRequest, status)
	a.Equal(api.CodeInvalidParam, env.Code)

	status, _ = do(e, http.MethodGet, "/api/v1/ftn/draws?from=20261009", "", nil)
	a.Equal(http.StatusBadRequest, status)
}

func TestGetDraw(t *testing.T) {
	a := assert.New(t)
	e := newTestRouter("")

	_, env := do(e, http.MethodGet, "/api/v1/ftn/draws/latest", "", nil)
	var d Draw
	a.NoError(json.Unmarshal(env.Data, &d))
	a.Equal("5990", d.Period)

	_, env = do(e, http.MethodGet, "/api/v1/ftn/draws/5988", "", nil)
	a.NoError(json.Unmarshal(env.Data, &d))
	a.Equal([]int{3, 7, 15, 17, 30}, d.Numbers)
	a.Equal(244, d.YearIndex)
	a.Equal("none", d.Features.Consecutive)

	status, env := do(e, http.MethodGet, "/api/v1/ftn/draws/1", "", nil)
	a.Equal(http.StatusNotFound, status)
	a.Equal(api.CodeNotFound, env.Code)
}

func TestNumberStats(t *testing.T) {
	a := assert.New(t)
	e := newTestRouter("")

	_, env := do(e, http.MethodGet, "/api/v1/ftn/numbers/stats?window=2", "", nil)
	var s NumberStats
	a.NoError(json.Unmarshal(env.Data, &s))
	a.Len(s.Numbers, 39)
	n23 := s.Numbers[22]
	a.Equal(23, n23.Number)
	a.Equal(2, n23.WindowCount)
	a.Equal(2, n23.TotalCount)
	a.Equal(0, n23.Omission)
	n3 := s.Numbers[2]
	a.Equal(0, n3.WindowCount)
	a.Equal(1, n3.TotalCount)
	a.Equal(2, n3.Omission)
	a.Equal(39, s.Numbers[38].Number)
	a.Equal(3, s.Numbers[38].Omission)

	status, _ := do(e, http.MethodGet, "/api/v1/ftn/numbers/stats?window=4", "", nil)
	a.Equal(http.StatusBadRequest, status)
}

func TestAnalyze(t *testing.T) {
	a := assert.New(t)
	e := newTestRouter("")

	_, env := do(e, http.MethodPost, "/api/v1/ftn/analyze", `{"numbers":[34,14,23,21,1]}`, nil)
	var r AnalyzeResult
	a.NoError(json.Unmarshal(env.Data, &r))
	a.Equal([]int{1, 14, 21, 23, 34}, r.Numbers)
	a.Equal(1, r.History.Match4)
	a.Equal(0, r.History.Match2) // 5989 只中 23, 5988 沒中
	a.Equal(3, r.TotalDraws)
	a.NotNil(r.LastMatch3)
	a.Equal("5990", r.LastMatch3.Period)

	for _, body := range []string{`{"numbers":[1,2,3,4]}`, `{"numbers":[1,2,3,4,40]}`, `{"numbers":[1,1,3,4,5]}`, `{}`} {
		status, env := do(e, http.MethodPost, "/api/v1/ftn/analyze", body, nil)
		a.Equal(http.StatusBadRequest, status, body)
		a.Equal(api.CodeInvalidParam, env.Code)
	}
}

func TestAPIKey(t *testing.T) {
	a := assert.New(t)
	e := newTestRouter("secret")

	status, env := do(e, http.MethodGet, "/api/v1/ftn/draws/latest", "", nil)
	a.Equal(http.StatusUnauthorized, status)
	a.Equal(api.CodeUnauthorized, env.Code)

	status, _ = do(e, http.MethodGet, "/api/v1/ftn/draws/latest", "", map[string]string{api.HeaderAPIKey: "secret"})
	a.Equal(http.StatusOK, status)

	status, _ = do(e, http.MethodGet, "/api/v1/health", "", nil)
	a.Equal(http.StatusOK, status)
}

func TestDataNotReady(t *testing.T) {
	a := assert.New(t)
	gin.SetMode(gin.TestMode)
	store := NewStore()
	e := gin.New()
	Register(e.Group("/ftn"), store)

	status, env := do(e, http.MethodGet, "/ftn/draws/latest", "", nil)
	a.Equal(http.StatusServiceUnavailable, status)
	a.Equal(api.CodeDataNotReady, env.Code)
}
