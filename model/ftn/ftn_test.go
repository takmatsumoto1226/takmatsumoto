package ftn

import (
	"encoding/json"
	"fmt"
	"lottery/algorithm"
	"lottery/config"
	"lottery/model/common"
	"lottery/model/df"
	"lottery/model/interf"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func Test_initNumberToIndex(t *testing.T) {
	// initNumberToIndex()
	// logrus.Info(numberToIndex)
	// fmt.Println(df.FeatureTenGroup1)
	// n := 14
	// fmt.Println(n / 10)
	// fmt.Println(df.Primes)
	// fmt.Println(bytes.IndexByte(df.Primes, 31))
	// fmt.Println(bytes.IndexByte(df.Primes, 30))
	// fmt.Println(df.FilterContinue3)

}

func Test_loadFTNs(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = FTNsManager{}
	as.Prepare()
	// aipicks := as.List.FeatureRange(*interf.SmartPureIntervalTH(0, 20))
	// sort.Sort(aipicks)
	// aipicks.ShowAll()
	// fmt.Println("")
	// fmt.Println("")
	// fmt.Println("")
	// distinct := aipicks.Distinct()
	// distinct.ShowAll()
	normalRange := as.List.WithRange(0, 20).Reverse()
	normalRange.ShowAll()

}

func Test_newFTNTest(t *testing.T) {
	df.DisableFilters([]int{df.FilterTailDigit})
	elems := strings.Split("2023,1230,312,04,11,17,20,32,5114", ",")
	ftn := NewFTN(elems)
	fmt.Println(ftn)

	elems2 := strings.Split("2023,1230,312,04,11,17,20,34,5114", ",")
	ftn2 := NewFTN(elems2)
	fmt.Println(ftn2)
	fmt.Println(ftn2.MatchFeature(ftn))
}

func Test_listLikeExecl(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = FTNsManager{}
	as.Prepare()
	arr := as.List.WithRange(2, 1)
	arr.ShowAll()
}

func Test_findnumbers(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = FTNsManager{}
	as.Prepare()
	fmt.Println("")
	fmt.Println("")
	p := PickParam{SortType: df.Descending, Interval: 20, Whichfront: df.Normal}
	list := as.List.WithRange(0, int(p.Interval)).Reverse()
	fmt.Println(list.Presentation())
	fmt.Println(list.BallsCountStatic().Presentation(false))

	result := algorithm.Combinations(as.List[0].toStringArray(), 3)
	for _, v := range result {
		fmt.Println("")
		fmt.Println("")
		fmt.Printf("=================== %s ================\n", v)
		as.List.findNumbers(v, df.Next).ShowAll()
	}
}

func Test_combination(t *testing.T) {
	// fmt.Println(algorithm.All([]string{"09", "14", "30"}))
	// fmt.Println(Ball39())
	balls := 1
	combarr := algorithm.Combinations(Ball39(), balls)
	for i, comb := range combarr {
		fmt.Println(comb)
		fmt.Println(i + 1)
	}
	bytes, err := json.Marshal(combarr)
	if err != nil {
		logrus.Error(err)
		return
	}

	err = os.WriteFile(fmt.Sprintf("./combination%d.json", balls), bytes, 0777)
	if err != nil {
		logrus.Error(err)
		return
	}
}

func Test_combination2(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = FTNsManager{}
	as.Prepare()

}

func Test_findDTree(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = FTNsManager{}
	as.Prepare()
	fmt.Println("")
	fmt.Println("")

	p := PickParam{SortType: df.Descending, Interval: 60, Whichfront: df.Normal}
	list := as.List.WithRange(0, int(p.Interval))
	fmt.Println(list.Presentation())
	fmt.Println(list.BallsCountStatic().Presentation(true))
	fmt.Println("")
	fmt.Println("")
	as.List.DTree(p).ShowAll()
}

func Test_findUTree(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = FTNsManager{}
	as.Prepare()
	fmt.Println("")
	fmt.Println("")

	p := PickParam{SortType: df.Descending, Interval: 60, Whichfront: df.Normal}
	as.List.ShowWithRange(int(p.Interval))
	params := PickParams{
		p,
	}
	as.IntervalBallsCountStatic(params)
	as.Picknumber(params)[p.GetKey()].ShowAll()
	fmt.Println("")
	fmt.Println("")
	as.List.UTree(p).ShowAll()
}

func Test_readwordFile(t *testing.T) {
	var i I

	i = &T{"Hello"}                // 把 type T 的值賦予給變數 i
	fmt.Printf("(%v, %T)\n", i, i) // i 的 dynamic value 是 &{Hello}、 dynamic type 是 *main.T
	i.M()                          // 意思是將 type T 對應的 value （&{Hello}） 來執行 type T 對應的 Ｍ 方法

	i = F(math.Pi)                 // 把 type F 的值賦予給變數 i
	fmt.Printf("(%v, %T)\n", i, i) // i 的 dynamic value 是 3.141、dynamic type 是 main.F
	i.M()                          // 意思是將 type F 對應的 value （3.1415） 去執行 type F 對應的 Ｍ 方法
}

type I interface {
	M()
}

// Type T 實作了 I interface
type T struct {
	S string
}

func (t *T) M() {
	fmt.Println(t.S)
}

// Type F 實作了 I interface
type F float64

func (f F) M() {
	fmt.Println(f)
}

func Test_loadCombination(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	loadCombination()
}

func Test_findDate(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = FTNsManager{}
	as.Prepare()
	fmt.Println("")
	fmt.Println("")
	p := PickParam{SortType: df.Descending, Interval: 20, Whichfront: df.Normal}
	as.List.ShowWithRange(int(p.Interval))
	params := PickParams{
		p,
	}
	as.IntervalBallsCountStatic(params)
	as.Picknumber(params)[p.GetKey()].ShowAll()
	fmt.Println("")
	fmt.Println("")
	fmt.Printf("=================== %s ================\n", "0322")
	as.findDate("0322", df.None).ShowAll()
}

func Test_GenerateTopPriceNumberJSON(t *testing.T) {
	fmt.Println("Start GenerateTopPriceNumberJSON....")
	defer common.TimeTaken(time.Now(), "Top Price Taken Time")
	config.LoadConfig("../../config.yaml")

	var ar = FTNsManager{}
	ar.Prepare()

	start := 0
	//
	df.DisableFilters([]int{df.FilterOddCount, df.FilterEvenCount, df.FilterTailDigit})
	// df.DisableFilters([]int{df.FilterTailDigit})
	th := interf.Threshold{
		Round:      20,
		Value:      7,
		SampleTime: 4,
		Sample:     len(ar.Combinations),
		Interval: interf.Interval{
			Index:  start,
			Length: len(ar.List)/5 + start,
		},
		Smart: interf.Smart{
			Enable: true,
			Type:   interf.RangeTypeLatestRange,
		},
		Randomer: 1,
		Match:    false,
	}

	ar.JSONGenerateTopPriceNumber(th)
	ar.SaveBTs()

}

// prediction
func Test_DoPrediction(t *testing.T) {
	defer common.TimeTaken(time.Now(), "DoPrediction")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	ar.Predictions(FileNames())

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

// func Test_listPredictionTops(t *testing.T) {
// 	defer common.TimeTaken(time.Now(), "listPredictionTops")
// 	config.LoadConfig("../../config.yaml")
// 	var ar = FTNsManager{}
// 	ar.Prepare()
// 	ar.ReadJSON(FileNames())

// 	for _, bt := range ar.BackTests {
// 		if len(bt.PickNumbers.PredictionTops) > 0 {
// 			fmt.Printf("\n\n\n%s", bt.Summery())
// 			fmt.Printf("PickNumbers.PredictionTops : %d\n", len(bt.PickNumbers.PredictionTops))
// 			bt.PickNumbers.PredictionTops.ShowAll()
// 		} else {
// 			if len(bt.ThresholdNumbers.PredictionTops) > 0 {
// 				fmt.Printf("\n\n\n%s", bt.Summery())
// 				fmt.Printf("ThresholdNumbers.PredictionTops : %d\n", len(bt.ThresholdNumbers.PredictionTops))
// 				bt.ThresholdNumbers.PredictionTops.ShowAll()
// 			}
// 		}

// 		if len(bt.ExcludeTops.PredictionTops) > 0 {
// 			fmt.Printf("\n\n\nExcludeTops:")
// 			fmt.Printf("ExcludeTops.PredictionTops : %d\n", len(bt.ExcludeTops.PredictionTops))
// 			bt.ExcludeTops.PredictionTops.ShowAll()
// 		}

// 	}
// }

func Test_DoBackTesting(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Back Test")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	ar.DoBackTesting(FileNames(), targetsub)

	fmt.Printf("PickNumbers:\n")
	for _, bt := range ar.BackTests {
		if len(bt.PickNumbers.TopMatch) > 0 {
			fmt.Println(bt.Summery())
			bt.PickNumbers.TopMatch.ShowAll()
		}
	}

	fmt.Printf("\n\n\n\nThresholdNumbers:\n")
	for _, bt := range ar.BackTests {
		if len(bt.ThresholdNumbers.TopMatch) > 0 {
			fmt.Println(bt.Summery())
			bt.ThresholdNumbers.TopMatch.ShowAll()
		}
	}

	fmt.Printf("\n\n\n\nExcludeTops:\n")
	for _, bt := range ar.BackTests {
		if len(bt.ExcludeTops.TopMatch) > 0 {
			fmt.Println(bt.Summery())
			bt.ThresholdNumbers.TopMatch.ShowAll()
		}
	}

}

// func Test_FindTop(t *testing.T) {
// 	defer common.TimeTaken(time.Now(), "FindTop")
// 	config.LoadConfig("../../config.yaml")
// 	var ar = FTNsManager{}
// 	ar.Prepare()
// 	ar.ReadJSON(FileNames())

// 	fmt.Printf("PickNumbers:\n")
// 	for _, bt := range ar.BackTests {
// 		if len(bt.PickNumbers.TopMatch) > 0 {
// 			fmt.Println(bt.Summery())
// 			bt.PickNumbers.TopMatch.ShowAll()
// 		}
// 	}

// 	fmt.Printf("\n\n\n\nThresholdNumbers:\n")
// 	for _, bt := range ar.BackTests {
// 		if len(bt.ThresholdNumbers.TopMatch) > 0 {
// 			fmt.Println(bt.Summery())
// 			bt.ThresholdNumbers.TopMatch.ShowAll()
// 		}
// 	}
// }

func Test_backtestreport(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	ar.BackTestingReports(FileNames())
}

func Test_repick(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Back Test")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	ar.ReadJSON(FileNames())
	df.DisableFilters([]int{df.FilterTailDigit})
	featuresFTNs := FTNArray{}
	for _, bt := range ar.BackTests {
		features := bt.Features.Balls
		for _, tn := range bt.ThresholdNumbers.Balls {
			for _, l := range features {
				if tn.MatchFeature(&l) {
					featuresFTNs = append(featuresFTNs, tn)
					break
				}
			}
		}
	}

	fmt.Println(len(featuresFTNs))

	tops := ar.List.WithRange(0, 1)
	for _, top := range tops {
		featuresFTNs.findNumbers(top.toStringArray(), df.None).ShowAll()
	}
}

func Test_pickupSum(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Back Test")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	ar.ReadJSON(FileNames())
	result := FTNArray{}
	for _, bt := range ar.BackTests {
		result = append(result, bt.PickNumbers.Balls...)
	}

	result.Distinct().ShowAll()
	fmt.Println(len(result))
	fmt.Println(len(result.Distinct()))
}

func Test_groupNumbers(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Group Index")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	ar.ReadJSON(FileNames())
	top := ar.List.GetNode(0)
	group := NewGroup(100, ar.Combinations, ar.List)
	p := PickParam{SortType: df.Descending, Interval: 20, Whichfront: df.Normal, Freq: 666}
	infl1 := ar.List.FragmentRange([]int{1})
	exfl2 := ar.List.FragmentRange([]int{0, 2})
	filterPick := ar.FilterByGroupIndex(group, 0).FilterHighFreqNumber(ar.List, p).FilterPickBySpecConfition().FilterIncludes(infl1, []int{}).FilterExcludes(exfl2, []int{}).FilterExcludeNode(ar.List).FilterNeighberNumber(&top, 2).FilterByTebGroup([]int{df.FeatureTenGroup1, df.FeatureTenGroup4}, []int{2, 1}).findNumbers([]string{}, df.None).Distinct()
	// filterPick := ar.FilterByGroupIndex(group, 1).FilterPickBySpecConfition().findNumbers([]string{"01", "04", "34", "35"}, df.None).Distinct()
	filterPick.ShowAll()
	fmt.Println(len(filterPick))
	fmt.Println(filterPick.IntervalBallsCountStatic(p).AppearBalls.Presentation(true))
	fmt.Println("got top")

	for _, f := range filterPick {
		if f.IsSame(&top) {
			fmt.Println("Oooooohhhhh My God!!!  it's " + f.formRow())
		}
	}

	fmt.Printf("\n\n\nGod Pick....\n")
	ar.GodPick(filterPick, 1)

	ar.List.WithRange(0, 20).Reverse().ShowAll()
}

func Test_FTNGroup(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Group Index")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	GroupCount := 100
	gftn := NewGroup(GroupCount, ar.Combinations, ar.List)
	fmt.Println(gftn)
}

func Test_FindGroupIndex(t *testing.T) {

	defer common.TimeTaken(time.Now(), "Find Group Index")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	GroupCount := 100
	group := NewGroup(GroupCount, ar.Combinations, ar.List)
	ftns := ar.List.WithRange(0, 5)
	for _, ftn := range ftns {
		v, k := group.FindGroupIndex(ftn)
		fmt.Printf("%4d:%2d => %s\n", v, k, ftn.formRow())
	}

}

func Find(slice interface{}, f func(value interface{}) bool) int {
	s := reflect.ValueOf(slice)
	if s.Kind() == reflect.Slice {
		for index := 0; index < s.Len(); index++ {
			if f(s.Index(index).Interface()) {
				return index
			}
		}
	}
	return -1
}

func Test_montecarlo(t *testing.T) {
	samplesExponent := 10

	var r1 float64
	var r2 float64
	var heads float64
	samples := math.Pow(10, float64(samplesExponent))
	heads = 0
	for range make([]struct{}, int(samples)) {
		r1 = rand.Float64()
		r2 = rand.Float64()
		toss := math.Pow(r1-0.5, 2) + math.Pow(r2-0.5, 2)
		if toss < 0.25 {
			heads++
		}
	}

	area := samples * 0.25

	pi := heads / area

	fmt.Printf("pi estimation - %f\n", pi)
}

func MultiPI(samples int, threads int) float64 {
	threadSamples := samples / threads
	results := make(chan float64, threads)

	for j := 0; j < threads; j++ {
		go func() {
			var inside int
			for i := 0; i < threadSamples; i++ {
				x, y := rand.Float64(), rand.Float64()

				if x*x+y*y <= 1 {
					inside++
				}
			}
			results <- float64(inside) / float64(threadSamples) * 4
		}()
	}

	var total float64
	for i := 0; i < threads; i++ {
		total += <-results
	}

	return total / float64(threads)
}

func Test_compareTest(t *testing.T) {

}

var targetsub = "20240612"

func FileNames() []string {

	fmt.Println("date : " + targetsub)
	return []string{
		filepath.Join(RootDir, targetsub, "content_07_4.0_20240612085440.json"), // p1
		filepath.Join(RootDir, targetsub, "content_07_4.0_20240612085612.json"), // p1
		filepath.Join(RootDir, targetsub, "content_07_4.0_20240612085742.json"), // p1
		filepath.Join(RootDir, targetsub, "content_07_4.0_20240612085915.json"), // h1e1
		filepath.Join(RootDir, targetsub, "content_07_4.0_20240612090046.json"), // h1
		// filepath.Join(RootDir, targetsub, "content_07_4.0_20240612090214.json"),
		// filepath.Join(RootDir, targetsub, "content_07_4.0_20240612090344.json"),
		filepath.Join(RootDir, targetsub, "content_07_4.0_20240612090513.json"), // p1
		// filepath.Join(RootDir, targetsub, "content_07_4.0_20240612090643.json"),
		// filepath.Join(RootDir, targetsub, "content_07_4.0_20240612090815.json"),
		// filepath.Join(RootDir, targetsub, "content_07_4.0_20240612101651.json"),
		filepath.Join(RootDir, targetsub, "content_07_4.0_20240612101823.json"), // p1
		filepath.Join(RootDir, targetsub, "content_07_4.0_20240612101955.json"), // h1
		filepath.Join(RootDir, targetsub, "content_07_4.0_20240612102130.json"), // h1
		filepath.Join(RootDir, targetsub, "content_07_4.0_20240612102300.json"), // p2
		filepath.Join(RootDir, targetsub, "content_07_4.0_20240612102432.json"), // h2e1
		filepath.Join(RootDir, targetsub, "content_07_4.0_20240612102609.json"), // p1
		// filepath.Join(RootDir, targetsub, "content_07_4.0_20240612102742.json"),
		filepath.Join(RootDir, targetsub, "content_07_4.0_20240612102911.json"), // h2e1
		filepath.Join(RootDir, targetsub, "content_07_4.0_20240612103041.json"), // h1
		filepath.Join(RootDir, targetsub, "content_07_4.0_20240612103212.json"), // p1
		// filepath.Join(RootDir, targetsub, "content_07_4.0_20240612103341.json"),
		filepath.Join(RootDir, targetsub, "content_07_4.0_20240612103511.json"), // p1
		// filepath.Join(RootDir, targetsub, "content_07_4.0_20240612103643.json"),
		filepath.Join(RootDir, targetsub, "content_07_4.0_20240612103816.json"), // p2
		// filepath.Join(RootDir, targetsub, "content_07_4.0_20240612103947.json"),
		filepath.Join(RootDir, targetsub, "content_07_4.0_20240612104111.json"), // p1
		filepath.Join(RootDir, targetsub, "content_07_4.0_20240612104235.json"), // h2
		// filepath.Join(RootDir, targetsub, "content_07_4.0_20240612104358.json"),
		// filepath.Join(RootDir, targetsub, "content_07_4.0_20240612104523.json"),
	}

	// files, _ := os.ReadDir(filepath.Join(RootDir, targetsub))
	// filenames := []string{}
	// for _, f := range files {
	// 	if strings.Contains(f.Name(), ".json") {
	// 		filenames = append(filenames, filepath.Join(RootDir, targetsub, f.Name()))
	// 	}
	// }
	// return filenames
}

func Test_CompareLatestAndHistoryFeature(t *testing.T) {
	defer common.TimeTaken(time.Now(), "CompareLatestAndHistoryFeature")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	df.DisableFilters([]int{df.FilterEvenCount, df.FilterOddCount})
	ar.CompareLatestAndHistoryFeature()
}

func Test_B2i(t *testing.T) {
	fmt.Println(B2i(false) & B2i(true))
}

func Test_static(t *testing.T) {
	defer common.TimeTaken(time.Now(), "static")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	fmt.Printf("%.2f%%\n", ar.List.WithRange(0, 20).StaticTotalNotInclude(1)*100)
}

func Test_ShowTwiceUP(t *testing.T) {
	defer common.TimeTaken(time.Now(), "static")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	fmt.Printf("%.2f%%\n", ar.List.StaticNumberShowTwiceup(2)*100)
}

func Test_ShowContinue2(t *testing.T) {
	defer common.TimeTaken(time.Now(), "static")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	r := interf.NewIntervalR(7)
	fmt.Printf("%.4f%%\n", ar.List.StaticContinue2Percent(r))
	ar.List.WithInterval(r).ShowAll()
}

func Test_ShowContinue3(t *testing.T) {
	defer common.TimeTaken(time.Now(), "static")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	r := interf.NewInterval(0, 0)
	fmt.Printf("%.4f%%\n", ar.List.StaticContinue3Percent(r))
	// ar.List.WithRange(0, r).Reverse().ShowAll()
}

func Test_FilterNeighberNumberTest(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_FilterNeighberNumberTest")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	top := ar.List.GetNode(0)
	tf := NewFTNWithInts([]int{17, 20, 28, 35, 36})
	fmt.Println(tf.haveNeighber(&top, 2))
}

func Test_StaticTenGroup(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_StaticTenGroup")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	r := interf.NewInterval(0, 0)
	tt := []int{
		df.FeatureTenGroup1,
		df.FeatureTenGroup2,
		df.FeatureTenGroup3,
		df.FeatureTenGroup4,
	}
	for _, t := range tt {
		fmt.Printf("%02d:\n", t)
		for _, v := range []int{0, 1, 2, 3, 4, 5} {
			fmt.Printf("%d : %.2f%%\n", v, ar.List.StaticGroupTen(r, t, v))
		}
		fmt.Print("\n\n")
	}

}
