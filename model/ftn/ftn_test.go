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
	normalRange := as.List.WithRange(0, 20)
	normalRange.ShowAll()

}

func Test_calculateTotalCount(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = FTNsManager{}
	as.Prepare()
	// interval := uint(len(as.FTNs) - 1)
	params := PickParams{
		// {SortType: descending, Interval: 360, Whichfront: biggerfront},
		// {SortType: descending, Interval: 180, Whichfront: biggerfront},
		// {SortType: descending, Interval: 60, Whichfront: biggerfront},
		// {SortType: descending, Interval: 48, Whichfront: biggerfront},
		// {SortType: descending, Interval: 36, Whichfront: biggerfront},
		// {SortType: df.Descending, Interval: 24, Whichfront: df.Biggerfront},
		{SortType: df.Descending, Interval: 12, Whichfront: df.Normal},
		// {SortType: descending, Interval: 5, Whichfront: biggerfront},
		// {SortType: descending, Interval: 2, Whichfront: biggerfront},
		// {SortType: descending, Interval: 1, Whichfront: biggerfront},
	}
	as.intervalBallsCountStatic(params)
	as.Picknumber(params)["0_12_2"].ShowAll()
	// for k, v := range ballPools {
	// 	logrus.Infof("%s:%v", k, v)
	// }
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
	list := as.List.WithRange(0, int(p.Interval))
	fmt.Println(list.Presentation())
	fmt.Println(list.BallsCountStatic().Presentation(true))

	result := algorithm.Combinations(as.RevList[0].toStringArray(), 3)
	for _, v := range result {
		fmt.Println("")
		fmt.Println("")
		fmt.Printf("=================== %s ================\n", v)
		as.List.findNumbers(v, df.Both).ShowAll()
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
	as.RevList.DTree(p).ShowAll()
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
	as.intervalBallsCountStatic(params)
	as.Picknumber(params)[p.GetKey()].ShowAll()
	fmt.Println("")
	fmt.Println("")
	as.RevList.UTree(p).ShowAll()
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
	as.intervalBallsCountStatic(params)
	as.Picknumber(params)[p.GetKey()].ShowAll()
	fmt.Println("")
	fmt.Println("")
	fmt.Printf("=================== %s ================\n", "0322")
	as.findDate("0322", df.None).ShowAll()
}

func Test_GenerateTopPriceNumberJSON(t *testing.T) {
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
		Value:      9,
		SampleTime: 5,
		Sample:     len(ar.Combinations),
		Interval: interf.Interval{
			Index:  start,
			Length: len(ar.List)/3 + start,
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
}

func Test_listPredictionTops(t *testing.T) {
	defer common.TimeTaken(time.Now(), "listPredictionTops")
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
		featuresFTNs.findNumbers(top.toStringArray(), df.None).ShowAll()
	}
}

func Test_DoBackTesting(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Back Test")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	ar.DoBackTesting(FileNames(), "20240528")

}

func Test_FindTop(t *testing.T) {
	defer common.TimeTaken(time.Now(), "FindTop")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	ar.ReadJSON(FileNames())

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
	group := NewFTNGroup(100, ar.Combinations, ar.RevList)
	filterPick := ar.RevList.FilterByGroupIndex(group, ar.BackTests)
	filterPick.Distinct().ShowAll()
	fmt.Println(len(filterPick.Distinct()))
	fmt.Println("這是分隔線-------")

	// p := PickParam{SortType: df.Descending, Interval: 20, Whichfront: df.Normal}
	// params := PickParams{
	// 	p,
	// }
	// ar.intervalBallsCountStatic(params)
	// static := ar.Picknumber(params)[p.GetKey()]
	// numbers := []string{}
	// for _, b := range static {
	// 	if b.Count > 3 {
	// 		numbers = append(numbers, b.Ball.Number)
	// 	}
	// }

	// findnumbers := FTNArray{}
	// for _, b := range numbers {
	// 	findnumbers = append(findnumbers, filterPick.findNumbers([]string{b}, df.None)...)
	// }
	// findnumbers.Distinct().ShowAll()
	// fmt.Println(len(findnumbers.Distinct()))
	// fmt.Println("這是分隔線-------2")

	filteragain := FTNArray{}
	for _, ftn := range filterPick.Distinct() {
		if ftn.Feature.IsContinue2() {
			filteragain = append(filteragain, ftn)
		}
	}
	// filteragain.Distinct().ShowAll()
	fmt.Println(len(filteragain.Distinct()))
	// fmt.Println("尋找-------")
	// filteragain.Distinct().findNumbers([]string{"16", "23", "29"}, df.None).ShowAll()d
	fmt.Println("這是分隔線-------3")

	picklen := 20
	filterExclide := FTNArray{}
	tops := ar.List.WithRange(0, picklen)
	for _, ftn := range filteragain {
		ecount := 0
		for _, t := range tops {
			if !ftn.IsExclude(t) {
				ecount++
			}
			if ecount == picklen {
				filterExclide = append(filterExclide, ftn)
			}
		}
	}

	diFileterExclude := filterExclide.Distinct()
	diFileterExclude.ShowAll()
	fmt.Println(len(diFileterExclude))
	fmt.Println(diFileterExclude.BallsCountStatic().Presentation(false))

}

func Test_FTNGroup(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Group Index")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	GroupCount := 500
	gftn := NewFTNGroup(GroupCount, ar.Combinations, ar.RevList)
	fmt.Println(gftn)
}

func Test_FindGroupIndex(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Find Group Index")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	GroupCount := 100
	group := NewFTNGroup(GroupCount, ar.Combinations, ar.RevList)
	ftns := ar.List.WithRange(2, 1)
	for _, ftn := range ftns {
		v, k := group.FindGroupIndex(ftn)
		fmt.Printf("%4d:%2d\n%s", v, k, ftn.formRow())
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

func FileNames() []string {
	targetsub := "20240529"
	fmt.Println("date : " + targetsub)
	return []string{
		filepath.Join(RootDir, targetsub, "content_09_5.0_20240529151540.json"), // ++
		// filepath.Join(RootDir, targetsub, "content_09_5.0_20240529151643.json"),
		filepath.Join(RootDir, targetsub, "content_09_5.0_20240529151746.json"), // ++
		filepath.Join(RootDir, targetsub, "content_09_5.0_20240529151848.json"), // ++2
		// filepath.Join(RootDir, targetsub, "content_09_5.0_20240529151950.json"),
		// filepath.Join(RootDir, targetsub, "content_09_5.0_20240529152055.json"),
		filepath.Join(RootDir, targetsub, "content_09_5.0_20240529152204.json"), // ++
		filepath.Join(RootDir, targetsub, "content_09_5.0_20240529152306.json"), // ++
		filepath.Join(RootDir, targetsub, "content_09_5.0_20240529152407.json"), // ++
		filepath.Join(RootDir, targetsub, "content_09_5.0_20240529152507.json"), // ++
		// filepath.Join(RootDir, targetsub, "content_09_5.0_20240529152606.json"),
		// filepath.Join(RootDir, targetsub, "content_09_5.0_20240529152707.json"),
		filepath.Join(RootDir, targetsub, "content_09_5.0_20240529152809.json"), // ++
		filepath.Join(RootDir, targetsub, "content_09_5.0_20240529152908.json"), //
		filepath.Join(RootDir, targetsub, "content_09_5.0_20240529153007.json"), // ++
		filepath.Join(RootDir, targetsub, "content_09_5.0_20240529153106.json"), // ++2
		filepath.Join(RootDir, targetsub, "content_09_5.0_20240529153205.json"), // ++
		// filepath.Join(RootDir, targetsub, "content_09_5.0_20240529153303.json"),
		// filepath.Join(RootDir, targetsub, "content_09_5.0_20240529153403.json"),
		filepath.Join(RootDir, targetsub, "content_09_5.0_20240529153501.json"), // ++
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
	df.DisableFilters([]int{df.FilterTailDigit})
	ar.CompareLatestAndHistoryFeature(interf.Interval{Index: 1, Length: len(ar.List) - 1})
}
