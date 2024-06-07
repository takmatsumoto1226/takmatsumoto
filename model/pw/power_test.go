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

	config.LoadConfig("../../config.yaml") // 17591400
	var pwm = PowerManager{}
	pwm.Prepare()
	// lens := len(combarr)
	df.DisableFilters([]int{df.FilterOddCount, df.FilterTenGroup, df.FilterTailDigit})
	start := 0
	th := interf.Threshold{
		Round:      1,
		Value:      14,
		SampleTime: 8,
		Sample:     len(pwm.Combinations),
		Interval: interf.Interval{
			Index:  start,
			Length: len(pwm.List)/3 + start,
		},
		Smart: interf.Smart{
			Enable: true,
			Type:   interf.RangeTypeLatestRange,
		},
		Randomer: 1,
	}

	bts := pwm.JSONGenerateTopPriceNumber(th)

	for _, bt := range bts {
		fn := filepath.Join(RootDir, SubDir, fmt.Sprintf("powercontent%s.json", bt.ID))
		common.SaveJSON(bt, fn)
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

func Test_listPredictionTops(t *testing.T) {
	defer common.TimeTaken(time.Now(), "listPredictionTops")
	config.LoadConfig("../../config.yaml")
	var ar = PowerManager{}
	ar.Prepare()
	ar.ReadJSON(FileNames())

	for _, bt := range ar.BackTests {
		if len(bt.PickNumbers.PredictionTops) > 0 {
			fmt.Printf("\n\n\n%s", bt.Summery())
			fmt.Printf("PickNumbers.PredictionTops : %d\n", len(bt.PickNumbers.PredictionTops))
			bt.PickNumbers.PredictionTops.ShowAll()
		} else {
			if len(bt.ThresholdNumbers.PredictionTops) > 0 {
				fmt.Printf("\n\n\n%s", bt.Summery())
				fmt.Printf("ThresholdNumbers.PredictionTops : %d\n", len(bt.ThresholdNumbers.PredictionTops))
				bt.ThresholdNumbers.PredictionTops.ShowAll()
			}
		}

		if len(bt.ExcludeTops.PredictionTops) > 0 {
			fmt.Printf("\n\n\nExcludeTops:")
			fmt.Printf("ExcludeTops.PredictionTops : %d\n", len(bt.ExcludeTops.PredictionTops))
			bt.ExcludeTops.PredictionTops.ShowAll()
		}

	}
}

func Test_backtestReport(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Back Test")
	config.LoadConfig("../../config.yaml")
	var ar = PowerManager{}
	ar.Prepare()
	ar.BackTestingReports(FileNames())
}

func Test_NewPowerGroupTest(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Back Test")
	config.LoadConfig("../../config.yaml")
	var ar = PowerManager{}
	ar.Prepare()
	GroupCount := 200
	gpw := NewPWGroup(GroupCount, ar.Combinations, ar.RevList)
	fmt.Println(gpw.Presentation())
}

func Test_CompareLatestAndHistoryFeature(t *testing.T) {
	defer common.TimeTaken(time.Now(), "CompareLatestAndHistoryFeature")
	config.LoadConfig("../../config.yaml")
	var ar = PowerManager{}
	ar.Prepare()
	df.DisableFilters([]int{df.FilterTailDigit, df.FilterEvenCount, df.FilterOddCount})
	ar.CompareLatestAndHistoryFeature()
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
