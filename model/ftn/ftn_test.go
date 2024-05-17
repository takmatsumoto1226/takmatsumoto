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
	elems := strings.Split("2023,1230,312,04,11,17,20,32,5114", ",")
	ftn := NewFTN(elems)
	fmt.Println(ftn)

	elems2 := strings.Split("2023,1230,312,04,11,17,20,33,5114", ",")
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
	as.List.ShowWithRange(int(p.Interval))
	params := PickParams{
		p,
	}
	as.intervalBallsCountStatic(params)
	static := as.Picknumber(params)[p.GetKey()]
	static.ShowAll()
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

func Test_continue(t *testing.T) {
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
	// as.Picknumber(params)[p.GetKey()].ShowAll()
	fmt.Println("")
	fmt.Println("")
	as.RevList.Continue4(p).ShowAll()
}

func Test_findDTree(t *testing.T) {
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
	df.DisableFilters([]int{df.FilterTenGroupOddCount, df.FilterTenGroupEvenCount})
	// df.DisableFilters([]int{df.FilterTailDigit})
	th := interf.Threshold{
		Round:      10,
		Value:      9,
		SampleTime: 6,
		Sample:     len(ar.Combinations),
		Interval: interf.Interval{
			Index:  start,
			Length: len(ar.List)/2 + start,
		},
		Smart: interf.Smart{
			Enable: true,
			Type:   interf.RangeTypeLatestRange,
		},
		Randomer: 1,
	}

	bts := ar.JSONGenerateTopPriceNumber(th)

	err := os.MkdirAll(filepath.Join(RootDir, time.Now().Format("20060102")), 0755)
	if err != nil {
		logrus.Error(err)
		return
	}
	for _, bt := range bts {
		fmt.Println(bt.Save())
	}
}

// prediction
func Test_DoPrediction(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Back Test")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	ar.Predictions(FileNames())
}

func Test_backtestreport(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	ar.BackTestingReports(FileNames())
}

func Test_listPridictionTops(t *testing.T) {
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
	ar.DoBackTesting(FileNames())
}

func Test_FindPickTop(t *testing.T) {
	defer common.TimeTaken(time.Now(), "FindPickTop")
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

	group := NewFTNGroup(200, ar.Combinations, ar.RevList)

	ar.ReadJSON(FileNames())
	filterPick := FTNArray{}

	for _, bt := range ar.BackTests {
		for _, ftn := range bt.PickNumbers.Balls {

			if gi, gcount := group.FindGroupIndex(ftn); gcount == 0 {
				fmt.Println(gi)
				fmt.Println(gcount)
				filterPick = append(filterPick, ftn)
			}
		}
	}

	// filterPick.ShowAll()
	// fmt.Println(len(filterPick))
	fmt.Println("這是分隔線-------")
	p := PickParam{SortType: df.Descending, Interval: 20, Whichfront: df.Normal}
	params := PickParams{
		p,
	}
	ar.intervalBallsCountStatic(params)
	static := ar.Picknumber(params)[p.GetKey()]
	numbers := []string{}
	for _, b := range static {
		if b.Count == 0 {
			numbers = append(numbers, b.Ball.Number)
		}
	}

	findnumbers := FTNArray{}
	for _, b := range numbers {
		findnumbers = append(findnumbers, filterPick.findNumbers([]string{b}, 1)...)
	}
	// findnumbers.Distinct().ShowAll()
	fmt.Println("這是分隔線-------2")
	fmt.Println(len(findnumbers.Distinct()))

	filteragain := FTNArray{}
	for _, ftn := range findnumbers {
		if ftn.Feature.IsContinue2() {
			filteragain = append(filteragain, ftn)
		}
	}
	filteragain.Distinct().ShowAll()
	fmt.Println(len(filteragain.Distinct()))

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
	GroupCount := 200
	gftn := NewFTNGroup(GroupCount, ar.Combinations, ar.RevList)
	ftn := NewFTNWithInts([]int{8, 18, 24, 29, 37})
	fmt.Println(gftn.FindGroupIndex(*ftn))
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
	targetsub := "20240517"
	fmt.Println("date : " + targetsub)
	return []string{
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517093355.json"), //
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517093511.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517093628.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517093744.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517093900.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517094018.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517094139.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517094256.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517094413.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517094529.json"), // ++
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517095758.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517095918.json"), //
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517100041.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517100200.json"), // ++
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517100327.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517100448.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517100610.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517100735.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517100856.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517101015.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517102213.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517102338.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517102500.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517102621.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517102742.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517102901.json"), //
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517103023.json"), //
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517103143.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517103303.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517103424.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517112920.json"), // ++
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517113045.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517113209.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517113332.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517113454.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517113616.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517113738.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517113901.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517114024.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517114146.json"), //
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517114721.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517114846.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517115011.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517115137.json"), // ++
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517115303.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517115429.json"), //
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517115554.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517115721.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517115843.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517120006.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517125757.json"), // ++
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517125920.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517130047.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517130212.json"), //
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517130340.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517130505.json"), // ++
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517130626.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517130748.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517130910.json"), //
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517131033.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517131155.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517131316.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517131438.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517131559.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517131718.json"), // ++++
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517131835.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517131955.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517132111.json"), // ++
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517133512.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517133636.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517133758.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517133919.json"), //
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517134042.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517134203.json"), // ++
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517134322.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517134443.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517134602.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517134722.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517134839.json"), //
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517134959.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517135116.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517135234.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517135352.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517135510.json"), // ++
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517135630.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517135747.json"), // ++
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517140925.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517141053.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517141215.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517141338.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517141501.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517141624.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517141753.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517141917.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517142042.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517142205.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517142328.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517142454.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517142618.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517142739.json"),
		filepath.Join(RootDir, targetsub, "content_10_6.0_20240517142904.json"), // ++
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517143029.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517143153.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517143315.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517143441.json"),
		filepath.Join(RootDir, targetsub, "content_10_6.0_20240517143603.json"), // ++
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517143726.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517143851.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517144012.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517144135.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517144255.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517144422.json"),
		filepath.Join(RootDir, targetsub, "content_10_6.0_20240517144545.json"), // ++
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517144707.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517144830.json"),
		filepath.Join(RootDir, targetsub, "content_10_6.0_20240517144954.json"), // ++
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517145116.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517145239.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517145404.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517145527.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517145650.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517145814.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517145936.json"),
		filepath.Join(RootDir, targetsub, "content_10_6.0_20240517150108.json"), // ++
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517150229.json"),
		// filepath.Join(RootDir, targetsub, "content_10_6.0_20240517150353.json"),
		// filepath.Join(RootDir, targetsub, "content_09_6.0_20240517153850.json"),
		filepath.Join(RootDir, targetsub, "content_09_6.0_20240517154120.json"), // ++
		// filepath.Join(RootDir, targetsub, "content_09_6.0_20240517154345.json"),
		// filepath.Join(RootDir, targetsub, "content_09_6.0_20240517154611.json"),
		filepath.Join(RootDir, targetsub, "content_09_6.0_20240517154843.json"), // ++
		// filepath.Join(RootDir, targetsub, "content_09_6.0_20240517155110.json"),
		filepath.Join(RootDir, targetsub, "content_09_6.0_20240517155349.json"), // ++
		// filepath.Join(RootDir, targetsub, "content_09_6.0_20240517155614.json"),
		// filepath.Join(RootDir, targetsub, "content_09_6.0_20240517155837.json"),
		// filepath.Join(RootDir, targetsub, "content_09_6.0_20240517160110.json"),
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
