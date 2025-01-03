package bl

import (
	"encoding/csv"
	"os"
	"sort"
	"strconv"
)

type BigLotteryList []BigLottery

func (fa BigLotteryList) Len() int {
	return len(fa)
}

func (fa BigLotteryList) Presentation() {
	for _, bl := range fa {
		bl.formRow()
	}
}

// Less ...
func (fa BigLotteryList) Less(i, j int) bool {
	ii, _ := strconv.Atoi(fa[i].TIdx)
	jj, _ := strconv.Atoi(fa[j].TIdx)
	return ii > jj
}

// Swap swaps the elements with indexes i and j.
func (fa BigLotteryList) Swap(i, j int) {
	fa[i], fa[j] = fa[j], fa[i]
}

func (fa BigLotteryList) WithRange(r int) BigLotteryList {
	al := len(fa)
	if r > 0 {
		return fa[al-r : al]
	}
	return fa
}

func (fa BigLotteryList) CSVExport1() {
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

func (fa BigLotteryList) CSVExport2() {
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

func (fa BigLotteryList) Reverse() BigLotteryList {
	sort.Sort(sort.Reverse(fa))
	return fa
}

func (fa BigLotteryList) CSVExport(fn string) {
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
