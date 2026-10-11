package lotteryapi

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Store 單一彩種歷史資料, 讀多寫少, 重新載入時整批替換
type Store struct {
	Game *Game

	loadMu   sync.Mutex // model 載入會用到套件層級的累計變數, 同時間只能有一個載入
	mu       sync.RWMutex
	records  []Record // 新到舊
	byPeriod map[string]int
	loadedAt time.Time
}

func NewStore(g *Game) *Store {
	return &Store{Game: g, byPeriod: map[string]int{}}
}

// NewStoreWithRecords 直接指定資料, 測試用
func NewStoreWithRecords(g *Game, records []Record) *Store {
	s := NewStore(g)
	s.set(s.clean(records))
	return s
}

// Reload 從 csv 重新載入
func (s *Store) Reload() error {
	s.loadMu.Lock()
	defer s.loadMu.Unlock()

	records := s.clean(s.Game.load())
	if len(records) == 0 {
		return errors.New("no " + s.Game.Code + " draw data")
	}
	s.set(records)
	return nil
}

// clean 排除格式錯誤的資料, 處理期數重複, 主號排序, 依期數由新到舊
func (s *Store) clean(records []Record) []Record {
	g := s.Game
	valid := make([]Record, 0, len(records))
	for _, r := range records {
		if err := g.validRecord(&r); err != nil {
			continue
		}
		p, _ := strconv.Atoi(r.Period)
		r.Period = strconv.Itoa(p) // 早期資料期數有前導零, 例 "037"
		r.Numbers = append([]int(nil), r.Numbers...)
		sort.Ints(r.Numbers)
		valid = append(valid, r)
	}

	result := dedupPeriods(valid)
	if skipped := len(records) - len(result); skipped > 0 {
		logrus.Warnf("%s: 略過 %d 筆格式錯誤或期數重複的資料", g.Code, skipped)
	}

	sort.SliceStable(result, func(i, j int) bool {
		pi, _ := strconv.Atoi(result[i].Period)
		pj, _ := strconv.Atoi(result[j].Period)
		return pi > pj
	})
	return result
}

// dedupPeriods 期數重複時(來源資料抄錯期數), 只保留日期落在前後期之間的那一筆;
// 完全相同的資料保留一筆, 無法判斷時全部捨棄
func dedupPeriods(records []Record) []Record {
	byPeriod := map[int][]Record{}
	for _, r := range records {
		p, _ := strconv.Atoi(r.Period)
		byPeriod[p] = append(byPeriod[p], r)
	}

	uniques := []int{}
	for p, rs := range byPeriod {
		if len(rs) == 1 {
			uniques = append(uniques, p)
		}
	}
	sort.Ints(uniques)

	result := make([]Record, 0, len(records))
	for p, rs := range byPeriod {
		if len(rs) == 1 {
			result = append(result, rs[0])
			continue
		}

		// 最接近的前一期與後一期(期數不重複者)
		i := sort.SearchInts(uniques, p)
		lo, hi := "", ""
		if i > 0 {
			lo = byPeriod[uniques[i-1]][0].Date()
		}
		if i < len(uniques) {
			hi = byPeriod[uniques[i]][0].Date()
		}

		fits := []Record{}
		for _, r := range rs {
			if (lo == "" || r.Date() >= lo) && (hi == "" || r.Date() <= hi) {
				fits = append(fits, r)
			}
		}
		if len(fits) > 0 && allSame(fits) {
			result = append(result, fits[0])
		}
	}
	return result
}

func allSame(rs []Record) bool {
	for _, r := range rs[1:] {
		if r.Date() != rs[0].Date() || r.Special != rs[0].Special || fmt.Sprint(r.Numbers) != fmt.Sprint(rs[0].Numbers) {
			return false
		}
	}
	return true
}

func (g *Game) validRecord(r *Record) error {
	if _, err := strconv.Atoi(r.Period); err != nil {
		return fmt.Errorf("invalid period %q", r.Period)
	}
	if len(r.Year) != 4 || len(r.MonthDay) != 4 {
		return fmt.Errorf("invalid date %s%s", r.Year, r.MonthDay)
	}
	if err := g.validNumbers(r.Numbers); err != nil {
		return err
	}
	if g.Special != nil && (r.Special < 1 || r.Special > g.Special.MaxNumber) {
		return fmt.Errorf("invalid special %d", r.Special)
	}
	return nil
}

func (s *Store) set(records []Record) {
	byPeriod := make(map[string]int, len(records))
	for i, r := range records {
		byPeriod[r.Period] = i
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.records = records
	s.byPeriod = byPeriod
	s.loadedAt = time.Now()
}

// snapshot 取得目前資料, 回傳的 slice 不可修改
func (s *Store) snapshot() ([]Record, time.Time) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.records, s.loadedAt
}

// findPeriod 依期數查詢, 回傳資料與其在新到舊列表中的位置
func (s *Store) findPeriod(period string) (Record, int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	i, ok := s.byPeriod[period]
	if !ok {
		return Record{}, 0, false
	}
	return s.records[i], i, true
}
