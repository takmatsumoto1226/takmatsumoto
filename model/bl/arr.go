package bl

import (
	"encoding/csv"
	"fmt"
	"lottery/model/common"
	"os"
	"sort"
	"strconv"
	"time"
)

type BLList []BL

func (fa BLList) Len() int {
	return len(fa)
}

func (fa BLList) Presentation(b_optional ...bool) string {
	result := ""
	for _, bl := range fa {
		result = result + bl.formRow(b_optional...) + "\n"
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
		if err := writer.Write(f.ToStringArr6()); err != nil {
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
	for _, pn := range fa {
		currentPrice := pn.AdariPrice(adari.IBalls[:6], adari.IBalls[6:])
		if currentPrice > 0 {
			fmt.Println(pn.simpleFormRow())
		}
		total = total + currentPrice
	}
	return total
}

func (fa BLList) FilterORIncludes(tops BLList, sb []int) BLList {
	defer common.TimeTaken(time.Now(), fmt.Sprintf("FilterIncludes : %d\n", len(fa)))

	// if len(tops) == 0 && len(sb) == 0 {
	// 	return fa
	// }
	result := BLList{}
	search := common.LIMap{}
	for _, t := range tops {
		for _, i := range t.Feature.IBalls {
			search[i] = true
		}
	}

	if len(sb) > 0 {
		for _, i := range sb {
			search[i] = true
		}
	}

	if len(search) == 0 {
		return fa
	}
	fmt.Println(search.Presentation())

	for _, ftn := range fa {
		for _, n := range ftn.Feature.IBalls {
			if _, ok := search[n]; ok {
				result = append(result, ftn)
				break
			}
		}
	}
	return result
}

func (fa BLList) FilterANDIncludes(tops BLList, sb []int) BLList {
	defer common.TimeTaken(time.Now(), fmt.Sprintf("FilterIncludes : %d\n", len(fa)))

	// if len(tops) == 0 && len(sb) == 0 {
	// 	return fa
	// }
	matchCount := len(sb)
	result := BLList{}
	search := common.LIMap{}
	for _, t := range tops {
		for _, i := range t.Feature.IBalls {
			search[i] = true
		}
	}

	if len(sb) > 0 {
		for _, i := range sb {
			search[i] = true
		}
	}

	if len(search) == 0 {
		return fa
	}
	fmt.Println(search.Presentation())

	for _, ftn := range fa {
		count := 0
		for _, n := range ftn.Feature.IBalls {
			if _, ok := search[n]; ok {
				count++
			}

			if count == matchCount {
				result = append(result, ftn)
				break
			}
		}
	}
	return result
}

func (fa BLList) FragmentRange(indexs []int) BLList {
	result := BLList{}
	for _, i := range indexs {
		result = append(result, fa[i])
	}
	return result
}
