package ftnapi

import (
	"errors"
	"lottery/model/ftn"
	"sync"
	"time"
)

// Store 今彩539 歷史資料, 讀多寫少, 重新載入時整批替換
type Store struct {
	loadMu   sync.Mutex // NewFTN 使用全域累計變數, 同時間只能有一個載入
	mu       sync.RWMutex
	draws    ftn.FTNArray // 新到舊
	byPeriod map[string]int
	loadedAt time.Time
}

func NewStore() *Store {
	return &Store{byPeriod: map[string]int{}}
}

// NewStoreWithDraws 直接指定資料(新到舊), 測試用
func NewStoreWithDraws(draws ftn.FTNArray) *Store {
	s := NewStore()
	s.set(draws)
	return s
}

// Reload 從 csv 重新載入
func (s *Store) Reload() error {
	s.loadMu.Lock()
	defer s.loadMu.Unlock()

	draws := ftn.LoadHistory()
	if len(draws) == 0 {
		return errors.New("no ftn draw data")
	}
	s.set(draws)
	return nil
}

func (s *Store) set(draws ftn.FTNArray) {
	byPeriod := make(map[string]int, len(draws))
	for i, d := range draws {
		byPeriod[d.TIdx] = i
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.draws = draws
	s.byPeriod = byPeriod
	s.loadedAt = time.Now()
}

// snapshot 取得目前資料, 回傳的 slice 不可修改
func (s *Store) snapshot() (ftn.FTNArray, time.Time) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.draws, s.loadedAt
}

func (s *Store) findPeriod(period string) (ftn.FTN, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	i, ok := s.byPeriod[period]
	if !ok {
		return ftn.FTN{}, false
	}
	return s.draws[i], true
}
