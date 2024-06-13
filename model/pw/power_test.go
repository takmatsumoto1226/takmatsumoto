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
	df.DisableFilters([]int{df.FilterOddCount, df.FilterTenGroup})
	start := 0
	th := interf.Threshold{
		Round:      1,
		Value:      3,
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

	p := PickParam{SortType: df.Descending, Interval: 30, Whichfront: df.Normal, Freq: 200}
	filter1 := pwm.List.FragmentRange([]int{})
	filter2 := pwm.List.WithRange(1, 0)
	filterPick := pwm.ListByGroupIndex(pwg, 0).FilterHighFreqNumber(pwm.List, p).FilterPickBySpecConfition().FilterIncludes(filter2, []int{}).FilterExcludes(filter1, []int{}).FilterExcludeNode(pwm.List).findNumbers([]string{}, df.None).Distinct()
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
	GodPick(filterPick, 1)
	// fmt.Println(pwm.List.WithRange(0, 20).Presentation())
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
	defer common.TimeTaken(time.Now(), "Test_CompareLatestAndHistoryFeature")
	config.LoadConfig("../../config.yaml")
	var ar = PowerManager{}
	ar.Prepare()
	df.DisableFilters([]int{df.FilterEvenCount, df.FilterOddCount})
	ar.CompareLatestAndHistoryFeature()
}

func Test_Continue2TypeStatic(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_Continue2TypeStatic")
	config.LoadConfig("../../config.yaml")
	var pwm = PowerManager{}
	pwm.Prepare()
	fmt.Printf("%.2f%%\n", pwm.List.StaticContinue2Percent(interf.Interval{Index: 0, Length: len(pwm.List)}))
}

func Test_Continue22TypeStatic(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_Continue22TypeStatic")
	config.LoadConfig("../../config.yaml")
	var pwm = PowerManager{}
	pwm.Prepare()
	fmt.Printf("%.2f%%\n", pwm.List.StaticContinue22Percent(interf.Interval{Index: 0, Length: len(pwm.List)}))
}

func Test_Continue3TypeStatic(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_Continue22TypeStatic")
	config.LoadConfig("../../config.yaml")
	var pwm = PowerManager{}
	pwm.Prepare()
	fmt.Printf("%.2f%%\n", pwm.List.StaticContinue3Percent(interf.Interval{Index: 0, Length: len(pwm.List)}))
}

func Test_Continue4TypeStatic(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_Continue22TypeStatic")
	config.LoadConfig("../../config.yaml")
	var pwm = PowerManager{}
	pwm.Prepare()
	fmt.Printf("%.2f%%\n", pwm.List.StaticContinue4Percent(interf.Interval{Index: 0, Length: len(pwm.List)}))
}

func Test_Continue4InAllCombinations(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_Continue4InAllCombinations")
	config.LoadConfig("../../config.yaml")
	var pwm = PowerManager{}
	pwm.Prepare()
	list := PowerList{}
	for _, v := range pwm.Combinations {
		balls := NewPowerWithInts(v)
		list = append(list, *balls)
	}
	fmt.Println(len(list.FilterPickBySpecConfition().Distinct()))
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
