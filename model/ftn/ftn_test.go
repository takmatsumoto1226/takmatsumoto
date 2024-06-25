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
	arr := as.List.WithRange(0, 21).Reverse()
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

	result := algorithm.Combinations(as.List[0].ToStringArr(), 3)
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
	date := "0614"
	fmt.Printf("=================== %s ================\n", date)
	as.List.FindDate(date, df.None).ShowAll()
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
		Round:      10,
		Value:      8,
		SampleTime: 5,
		Sample:     len(ar.Combinations),
		Interval: interf.Interval{
			Index:  start,
			Length: len(ar.List) / 3,
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
	defer common.TimeTaken(time.Now(), "Test_DoPrediction")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	interval := interf.Interval{Index: 0, Length: 20}
	tops := ar.List.WithRange(interval.Index, interval.Length)
	ar.Predictions(FileNames(), tops)
}

func Test_ListPredictions(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_ListPredictions")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
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

func Test_DoBackTesting(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Back Test")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	ar.DoBackTesting(FileNames(), "20240622")

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
		featuresFTNs.findNumbers(top.ToStringArr(), df.None).ShowAll()
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
	// ar.ReadJSON(FileNames())
	top := ar.List.GetNode(0)
	p := PickParam{SortType: df.Descending, Interval: 20, Whichfront: df.Normal, Freq: 655}
	GroupCount := 100
	group := NewGroup(GroupCount, ar.Combinations, ar.List)

	filterPick := ar.
		// FilterByGroupIndex(NewGroup(100, ar.Combinations, ar.List), []int{0}).
		FullCombination().
		// FilterHighFreqNumber(ar.List, p).
		FilterPickBySpecConfition().
		FilterIncludes(ar.List.FragmentRange([]int{0}), []int{}).
		FilterExcludes(ar.List.FragmentRange([]int{}), []int{}).
		FilterExcludeNode(ar.List).
		FilterCol(&top, 1).
		FilterNeighber(&top, 2).
		FilterByTebGroup([]int{df.FeatureTenGroup3}, []int{2, 2}).
		FilterFeatureExcludes(ar.List).
		findNumbers([]string{}, df.None).
		FilterByGroupIndex(group, []int{0}).
		FilterOddEvenList(2).
		Distinct()

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
	ar.GodPick(filterPick, 10)

	ar.List.WithRange(0, 20).Reverse().ShowAll()
}

func Test_FTNGroup(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Group Index")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	GroupCount := 50
	gftn := NewGroup(GroupCount, ar.Combinations, ar.List)
	gftn.StaticCounts()
	fmt.Println(gftn.Presentation())
	// ar.ReadJSON(FileNames())
	// filterPick := ar.FilterByGroupIndex(gftn, []int{0})
	// filterPick.ShowAll()
	// fmt.Println(len(filterPick))
	// p := PickParam{SortType: df.Descending, Interval: 20, Whichfront: df.Normal, Freq: 666}
	// fmt.Println(filterPick.IntervalBallsCountStatic(p).AppearBalls.Presentation(true))
}

func Test_FindGroupIndex(t *testing.T) {

	defer common.TimeTaken(time.Now(), "Find Group Index")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	GroupCount := 100
	group := NewGroup(GroupCount, ar.Combinations, ar.List)
	ftns := ar.List.WithRange(0, 1)
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

func Test_compareTest(t *testing.T) {

}

var targetsub = "20240625"

func FileNames() []string {

	fmt.Println("date : " + targetsub)
	return []string{
		filepath.Join(RootDir, targetsub, "content_08_5.0_20240625093523.json"),
		filepath.Join(RootDir, targetsub, "content_08_5.0_20240625093615.json"),
		filepath.Join(RootDir, targetsub, "content_08_5.0_20240625093707.json"),
		filepath.Join(RootDir, targetsub, "content_08_5.0_20240625093754.json"),
		filepath.Join(RootDir, targetsub, "content_08_5.0_20240625093842.json"),
		filepath.Join(RootDir, targetsub, "content_08_5.0_20240625093932.json"),
		filepath.Join(RootDir, targetsub, "content_08_5.0_20240625094020.json"),
		filepath.Join(RootDir, targetsub, "content_08_5.0_20240625094108.json"),
		filepath.Join(RootDir, targetsub, "content_08_5.0_20240625094155.json"),
		filepath.Join(RootDir, targetsub, "content_08_5.0_20240625094243.json"),
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
	df.DisableFilters([]int{df.FilterEvenCount, df.FilterOddCount, df.FilterTailDigit})
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
	r := interf.NewIntervalR(30)
	list := ar.List.WithInterval(r)
	list.ShowAll()
	fmt.Printf("%.4f%%\n", list.StaticContinue2Percent())
	list.Continue2s().ShowAll()

}

func Test_ShowContinue2Avg(t *testing.T) {
	defer common.TimeTaken(time.Now(), "static")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	r := interf.NewIntervalR(20)
	list := ar.List.WithInterval(r)
	list.ShowAll()
	for i := 0; i < 1000; i++ {
		rl := ar.List.WithInterval(interf.NewInterval(i, 20))
		fmt.Printf("%s => %.2f%%\n", rl[0].Date(), rl.StaticContinue2Percent())
	}
	fmt.Println(len(list.Continue2s()))
	list.Continue2s().ShowAll()

}
func Test_ShowContinue3(t *testing.T) {
	defer common.TimeTaken(time.Now(), "static")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	r := interf.NewInterval(0, 0)
	ar.List.WithInterval(r)
	fmt.Printf("%.4f%%\n", ar.List.StaticContinue3Percent(true))
	// ar.List.WithRange(0, r).Reverse().ShowAll()
}

func Test_ShowContinue3Avg(t *testing.T) {
	defer common.TimeTaken(time.Now(), "static")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	ir := 1000
	r := interf.NewIntervalR(ir)
	list := ar.List.WithInterval(r).Reverse()
	list.ShowAll()
	for i := 0; i < ir; i++ {
		rl := ar.List.WithInterval(interf.NewInterval(i, 20))
		fmt.Printf("%.4f%%\n", rl.StaticContinue3Percent(false))
	}
	fmt.Println(len(list.Continue3s()))
	list.Continue3s().ShowAll()

}

func Test_FilterNeighberNumberTest(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_FilterNeighberNumberTest")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	top := ar.List.GetNode(0)
	tf := NewFTNWithStrings([]string{})
	fmt.Println(tf.haveNeighber(&top, 0))
}

func Test_StaticTenGroup(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_StaticTenGroup")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	r := interf.NewInterval(0, 100)
	tt := []int{
		df.FeatureTenGroup1,
		df.FeatureTenGroup2,
		df.FeatureTenGroup3,
		df.FeatureTenGroup4,
	}
	report := ""
	for _, t := range tt {
		report = report + fmt.Sprintf("%02d:\n", t+1)
		for _, v := range []int{0, 1, 2, 3, 4, 5} {
			report = report + fmt.Sprintf("%d : %.2f%%  ", v, ar.List.WithInterval(r).StaticGroupTen(t, v))
		}
		report = report + "\n\n"
	}
	fmt.Println(report)
}

func Test_StaticTenGroupAvg(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_StaticTenGroup")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	tt := []int{
		df.FeatureTenGroup1,
		df.FeatureTenGroup2,
		df.FeatureTenGroup3,
		df.FeatureTenGroup4,
	}
	report := ""
	for _, t := range tt {
		report = report + fmt.Sprintf("%02d:\n", t+1)
		for i := 0; i < 100; i++ {
			r := interf.NewInterval(i, 30)
			report = report + fmt.Sprintf(" %03d => ", i)
			for _, v := range []int{0, 1, 2, 3, 4, 5} {
				report = report + fmt.Sprintf("%d : %4.2f%%   ", v, ar.List.WithInterval(r).StaticGroupTen(t, v))
			}
			report = report + "\n"
		}
		report = report + "\n\n"
	}
	report = report + "\n\n"

	fmt.Println(report)
}

func Test_ListTenGroup(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_StaticTenGroup")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	r := interf.NewInterval(0, 30)

	report := ""
	report = report + ar.List.WithInterval(r).Reverse().PresentationGroupTenWithRange(0)
	fmt.Println(report)
}

func Test_ListTenGroupByKey(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_StaticTenGroup")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	r := interf.NewInterval(0, 0)

	report := ""
	report = report + ar.List.WithInterval(r).StaticTenGroupByTKey().Presentation()
	fmt.Println(report)
}

func Test_ListSameWithTopTenGroups(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_StaticTenGroup")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	df.DisableFilters([]int{
		df.FilterOddCount,
		df.FilterEvenCount,
		df.FilterTailDigit,
		df.FilterPrimeCount,
		df.FilterContinueRowType,
	})
	for i := 0; i < 100; i++ {
		tops := ar.List.WithRange(i, 1)
		list := ar.List.FindFeature(tops)
		fmt.Printf("%s%d:%.2f%%\n", tops.Presentation(), len(list), float64(len(list))/float64(len(ar.List))*100)
		fmt.Printf("\n")
	}

}

func Test_FeatureStatic(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_FeatureStatic")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	mapping := map[string]int{}
	for _, f := range ar.List {
		if v, ok := mapping[f.Feature.Key]; ok {
			mapping[f.Feature.Key] = v + 1
		} else {
			mapping[f.Feature.Key] = 1
		}
	}

	for k, m := range mapping {
		fmt.Printf("%s:%d\n", k, m)
	}
}

func Test_allcombinationFeatureGroup(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_allcombination")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	top := ar.List.GetNode(0)
	fmt.Println(top)
	mapping := map[string]int{}
	arr := FTNArray{}

	for _, c := range ar.Combinations {
		arr = append(arr, *NewFTNWithInts(c))
	}

	for _, f := range arr {
		if v, ok := mapping[f.Feature.Key]; ok {
			mapping[f.Feature.Key] = v + 1
		} else {
			mapping[f.Feature.Key] = 1
		}
	}

	arrc := map[int]int{}
	for k, m := range mapping {
		fmt.Printf("%s:%d\n", k, m)
		if v, ok := arrc[m]; ok {
			arrc[m] = v + 1
		} else {
			arrc[m] = 1
		}
	}
	fmt.Println(len(mapping))
	for idx, c := range arrc {
		fmt.Printf("%5d個的有 %d\n", idx, c)
	}
}

func Test_StaticColPercent(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_StaticColPercent")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	n := 1
	for i := 0; i <= 100; i++ {
		r := interf.NewInterval(i, 20)
		rl := ar.List.WithInterval(r)
		fmt.Printf("%s : %.4f%%\n", rl[0].Date(), rl.StaticColPercent(n))
	}
	ar.List.WithInterval(interf.NewInterval(0, 100)).Cols(n).ShowAll()
}

func Test_StaticColPercentAll(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_StaticColPercent")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	fmt.Printf("%.4f%%\n", ar.List.StaticColPercent(1))
}

func Test_Cols(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_StaticColPercent")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	ar.List.Cols(1).ShowAll()
}

func Test_StaticHaveNeighberPercent(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_StaticHaveNeighberPercent")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	for i := 0; i <= 100; i++ {
		r := interf.NewInterval(i, 20)
		rl := ar.List.WithInterval(r)
		fmt.Printf("%s : %.2f%%\n", rl[0].Date(), rl.StaticHaveNeighberPercent(3))
	}
}

func Test_Neighbers(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_StaticColPercent")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	n := 2
	ar.Prepare()
	sl := ar.List.WithRange(0, 0).Neighbers(n)
	sl.ShowAll()
	fmt.Println(len(sl) / 3)
	fmt.Printf("%.1f%%\n", sl.StaticHaveNeighberPercent(n))
}

func Test_FoundGroups(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Group Index")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	GroupCount := 50
	group := NewGroup(GroupCount, ar.Combinations, ar.List)
	gftns := group.FindGroupNotes(6)
	gftns.ShowAll()
	result := FTNArray{}
	for _, top := range ar.List {
		for _, gftn := range gftns {
			if gftn.IsSame(&top) {
				result = append(result, top)
				break
			}
		}
	}
	result.ShowAll()
}

func Test_StaticFullTenGrouopPercent(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_StaticFullTenGrouopPercent")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	length := 20
	for i := 0; i < 200; i++ {
		r := interf.NewInterval(i, length)
		arr := ar.List.WithInterval(r)
		fmt.Printf("%04d : %.3f%%\n", i, arr.StaticFullTenGroupPercent())
	}

	r := interf.NewInterval(0, length)
	result := FTNArray{}
	for _, f := range ar.List.WithInterval(r) {
		if f.Feature.IsFullTenGrouop() {
			result = append(result, f)
		}
	}
	result.Reverse().ShowAll()
}

func Test_0ShowTest(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_0ShowTest")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	length := 5
	for y := 1; y <= 20; y++ {
		count := 0
		for i := 0; i < length*y; i++ {
			r := interf.Interval{Index: i, Length: 1}
			tops := ar.List.WithInterval(r)
			exfl2 := ar.List.FragmentRange([]int{i + 1})
			result := tops.FilterExcludes(exfl2, []int{})
			// result.ShowAll()
			if len(result) > 0 {
				count++
			}
		}
		fmt.Printf("%2d:%.2f\n", length*y, float64(count)/float64(length*y)*100)
	}
}

func Test_StaticExclude(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_NeightberStatic")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	for i := 0; i < 20; i++ {
		r := interf.NewInterval(i, 21)
		sl := ar.List.WithInterval(r)
		fmt.Printf("%s : %.2f%%\n", sl[0].Date(), sl.StaticExclude(1, false))
	}
	excludes := ar.List.WithRange(0, 20).Exclude(1)
	excludes.Reverse().ShowAll()
	fmt.Println(len(excludes))
	ar.List.WithRange(0, 20).Reverse().ShowAll()

}

func Test_StaticHaveNeighberAndColsPercent(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_NeightberStatic")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	ar.List.WithRange(0, 200).NeighberAndCols(2, 1).Reverse().ShowAll()
}

func Test_StaticOddList(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_NeightberStatic")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	ar.List.WithRange(0, 100).FilterOddEvenList(2).Reverse().ShowAll()
}
