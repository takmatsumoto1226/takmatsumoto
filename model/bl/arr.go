package bl

import (
	"encoding/csv"
	"os"
	"sort"
	"strconv"
)

type BLList []BL

func (fa BLList) Len() int {
	return len(fa)
}

func (fa BLList) Presentation() string {
	result := ""
	for _, bl := range fa {
		result = result + bl.formRow() + "\n"
	}
	return result
}

// Less ...
func (fa BLList) Less(i, j int) bool {
	ii, _ := strconv.Atoi(fa[i].TIdx)
	jj, _ := strconv.Atoi(fa[j].TIdx)
	return ii > jj
}

// Swap swaps the elements with indexes i and j.
func (fa BLList) Swap(i, j int) {
	fa[i], fa[j] = fa[j], fa[i]
}

func (fa BLList) WithRange(r int) BLList {
	al := len(fa)
	if r > 0 {
		return fa[al-r : al]
	}
	return fa
}

func (fa BLList) CSVExport1() {
	// 建立 CSV 檔案
	file, err := os.Create("output.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 建立 CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 將資料寫入 CSV
	for _, f := range fa {
		if err := writer.Write(f.toStringArray()); err != nil {
			panic(err)
		}
	}

	// 檢查是否有錯誤
	if err := writer.Error(); err != nil {
		panic(err)
	}

	println("CSV 檔案輸出完成！")
}

func (fa BLList) CSVExport2() {
	// 建立 CSV 檔案
	file, err := os.Create("output.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 建立 CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 將資料寫入 CSV
	for _, f := range fa {
		if err := writer.Write(f.toStringArray2()); err != nil {
			panic(err)
		}
	}

	// 檢查是否有錯誤
	if err := writer.Error(); err != nil {
		panic(err)
	}

	println("CSV 檔案輸出完成！")
}

func (fa BLList) Reverse() BLList {
	sort.Sort(sort.Reverse(fa))
	return fa
}

func (fa BLList) CSVExport(fn string) {
	// 建立 CSV 檔案
	file, err := os.Create(fn)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 建立 CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 將資料寫入 CSV
	for _, f := range fa {
		if err := writer.Write(f.ToStringArr()); err != nil {
			panic(err)
		}
	}

	// 檢查是否有錯誤
	if err := writer.Error(); err != nil {
		panic(err)
	}

	println("CSV 檔案輸出完成！")

}

func (fa BLList) AdariPrice(adari *BL) int {
	total := 0
	// for _, pn := range fa {
	// 	currentPrice := pn.AdariPrice(adari)
	// 	total = total + currentPrice
	// }
	return total
}
