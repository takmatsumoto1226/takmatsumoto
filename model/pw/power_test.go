package pw

import (
	"fmt"
	"lottery/algorithm"
	"lottery/config"
	"lottery/model/common"
	"lottery/model/df"
	"lottery/model/interf"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func Test_listLikeExecl(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = PowerManager{numberToIndex: map[string]int{}}
	as.Prepare()
	fmt.Println(as.List.WithRange(0, 20).Presentation())
}

func Test_findnumber(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = PowerManager{numberToIndex: map[string]int{}}
	as.Prepare()
	result := algorithm.Combinations(as.RevList[0].toStringArray(), 3)
	for _, v := range result {
		fmt.Println("")
		fmt.Println("")
		fmt.Printf("=================== %s ================\n", v)
		as.List.findNumbers(v, df.Next).Presentation()
	}
}

func Test_random(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Top Price Taken Time")
	// Create and seed the generator.
	// Typically a non-fixed seed should be used, such as time.Now().UnixNano().
	// Using a fixed seed will produce the same output on every run.
	config.LoadConfig("../../config.yaml") // 17591400
	var as = PowerManager{}
	as.Prepare()
	// lens := len(combarr)
	df.DisableFilters([]int{df.FilterOddCount, df.FilterTenGroup, df.FilterTailDigit})

	th := interf.Threshold{
		Round:      1,
		Value:      14,
		SampleTime: 8,
		Sample:     len(as.Combinations),
		Interval: interf.Interval{
			Index:  2,
			Length: 20,
		},
		Smart: interf.Smart{
			Enable: true,
			Type:   interf.RangeTypeLatestRange,
		},
		Randomer: 1,
	}

	// th := interf.Threshold{Round: 1, Value: 26, SampleTime: 10, Sample: len(combarr)}

	bts := as.JSONGenerateTopPriceNumber(th)

	for r, bt := range bts {

		fn := filepath.Join(RootDir, SubDir, fmt.Sprintf("powercontent%s.json", bt.ID))
		common.SaveJSON(bt, fn, r+1)
	}

}

func Test_backtesting(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Back Test")
	config.LoadConfig("../../config.yaml")
	var ar = PowerManager{}
	ar.Prepare()
	ar.ReadJSON(FileNames())
	ar.Predictions()
}

func Test_backtestReport(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Back Test")
	config.LoadConfig("../../config.yaml")
	var ar = PowerManager{}
	ar.Prepare()
	ar.BackTestingReports(FileNames())
}

// func FileNames() []string {
// 	return []string{
// 		filepath.Join(RootDir, SubDir, "powercontent20240513154337"),
// 		filepath.Join(RootDir, SubDir, "powercontent20240513155407"),
// 		filepath.Join(RootDir, SubDir, "powercontent20240513160317"),
// 	}
// }

func FileNames() []string {
	files, _ := os.ReadDir(filepath.Join(RootDir, SubDir))
	filenames := []string{}
	for _, f := range files {
		if strings.Contains(f.Name(), ".json") {
			filenames = append(filenames, filepath.Join(RootDir, SubDir, f.Name()))
		}
	}
	return filenames
}
