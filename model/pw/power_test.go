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
	p := PickParam{SortType: df.Descending, Interval: 20, Whichfront: df.Normal, Freq: 200}
	list := as.List.WithRange(0, int(p.Interval)).Reverse()
	fmt.Println(list.Presentation())
	ballsCount := list.IntervalBallsCountStatic(p)
	fmt.Println(ballsCount.Presentation(false))
}

func Test_findnumber(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = PowerManager{numberToIndex: map[string]int{}}
	as.Prepare()
	result := algorithm.Combinations(as.List[0].toStringArray(), 3)
	for _, v := range result {
		fmt.Println("")
		fmt.Println("")
		fmt.Printf("=================== %s ================\n", v)
		fmt.Println(as.List.Reverse().findNumbers(v, df.Next).Presentation())
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
		Value:      2,
		SampleTime: 1,
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

	pwm.JSONGenerateTopPriceNumber(th)
	pwm.SaveBTs()

}

func Test_Predictions(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_Predictions")
	config.LoadConfig("../../config.yaml")
	var ar = PowerManager{}
	ar.Prepare()
	ar.ReadJSON(FileNames())
	ar.Predictions()

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

func Test_PickupNumber(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_PickupNumber")
	config.LoadConfig("../../config.yaml")
	var pwm = PowerManager{}
	pwm.Prepare()
	pwm.ReadJSON(FileNames())

	GroupCount := 500
	pwg := NewPWGroup(GroupCount, pwm.Combinations, pwm.List)

	p := PickParam{SortType: df.Descending, Interval: 20, Whichfront: df.Normal, Freq: 200}
	filter1 := pwm.List.FragmentRange([]int{0})
	filter2 := pwm.List.WithRange(1, 1)
	filterPick := pwm.ListByGroupIndex(pwg, 0).FilterHighFreqNumber(pwm.List, p).FilterPickBySpecConfition().FilterIncludes(filter2, []int{1, 38}).FilterExcludes(filter1, []int{15}).FilterExcludeNode(pwm.List).findNumbers([]string{"16"}, df.None).Distinct()
	filterPick.ShowAll()
	fmt.Println(len(filterPick))
	fmt.Println(filterPick.IntervalBallsCountStatic(p).Presentation(false))
	fmt.Println("got top")
	top := pwm.List.GetNode(0)
	for _, f := range filterPick {
		if f.IsSame(&top) {
			fmt.Println(f.formRow())
		}
	}

	fmt.Printf("\n\n\nGod Pick....\n")
	GodPick(filterPick, 10)
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
	GroupCount := 500
	gpw := NewPWGroup(GroupCount, ar.Combinations, ar.List)
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
// 		filepath.Join(RootDir, SubDir, "content_07_4.0_20240613132929.json"),
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
