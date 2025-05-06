package cache

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Record struct {
	Year       string
	MonthDay   string
	YearIndex  int
	Num1       string
	Num2       string
	Num3       string
	Num4       string
	Num5       string
	TotalIndex int
}

type Cache struct {
	db  *sql.DB
	ttl time.Duration
}

func NewLotteryDB(dbPath string) (*Cache, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	ldb := &Cache{db: db}
	if err := ldb.migrate(); err != nil {
		return nil, err
	}
	return ldb, nil
}

// New 建立新的快取物件，指定檔案路徑與TTL時間
func (c *Cache) migrate() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS record (
		year VARCHAR(4),
		month_day VARCHAR(4),
		year_index INTEGER,
		num1 VARCHAR(3),
		num2 VARCHAR(3),
		num3 VARCHAR(3),
		num4 VARCHAR(3),
		num5 VARCHAR(3),
		total_index INTEGER,
		PRIMARY KEY (year, month_day, year_index)
	);`
	_, err := c.db.Exec(createTableSQL)
	if err != nil {
		log.Println("migrate error:", err)
	}
	return err
}

// Set 寫入快取資料
func (c *Cache) InsertRecord(r Record) error {
	_, err := c.db.Exec(`INSERT OR REPLACE INTO record
		(year, month_day, year_index, num1, num2, num3, num4, num5, total_index)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		r.Year, r.MonthDay, r.YearIndex, r.Num1, r.Num2, r.Num3, r.Num4, r.Num5, r.TotalIndex)
	return err
}

// 查資料
func (c *Cache) GetRecord(year, monthDay string, yearIndex int) (Record, bool) {
	var r Record
	err := c.db.QueryRow(`SELECT year, month_day, year_index, num1, num2, num3, num4, num5, total_index
		FROM record WHERE year = ? AND month_day = ? AND year_index = ?`,
		year, monthDay, yearIndex).
		Scan(&r.Year, &r.MonthDay, &r.YearIndex, &r.Num1, &r.Num2, &r.Num3, &r.Num4, &r.Num5, &r.TotalIndex)
	if err != nil {
		if err == sql.ErrNoRows {
			return Record{}, false
		}
		log.Println("查詢錯誤:", err)
		return Record{}, false
	}
	return r, true
}

// ClearExpired 移除過期快取資料
func (c *Cache) ClearExpired() error {
	threshold := time.Now().Add(-c.ttl)
	_, err := c.db.Exec("DELETE FROM cache WHERE created_at < ?", threshold)
	return err
}

// Close 關閉資料庫連線
func (c *Cache) Close() error {
	return c.db.Close()
}
