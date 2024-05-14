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
	"gonum.org/v1/gonum/stat/combin"
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
	as.Picknumber(params)[p.GetKey()].ShowAll()
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

func Test_findAdariPrice(t *testing.T) {
	balls := 5
	combarr := combin.Combinations(10, balls)
	arr := FTNArray{}
	for _, comb := range combarr {
		arr = append(arr, *NewFTNWithInts(comb))
	}
	arr.AdariPrice(NewFTNWithInts([]int{2, 4, 5, 6, 7}))
}

func Test_GenerateTopPriceNumberJSON(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Top Price Taken Time")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()

	start := 1
	//
	df.DisableFilters([]int{df.FilterTenGroupEvenCount, df.FilterTenGroupOddCount})
	// df.DisableFilters([]int{df.FilterTailDigit})
	th := interf.Threshold{
		Round:      5,
		Value:      9,
		SampleTime: 5,
		Sample:     len(ar.Combinations),
		Interval: interf.Interval{
			Index:  start,
			Length: len(ar.List)/5 - start,
		},
		Smart: interf.Smart{
			Enable: true,
			Type:   interf.RangeTypeLatestRange,
		},
		Randomer: 1,
	}

	bts := ar.JSONGenerateTopPriceNumber(th)

	err := os.MkdirAll(filepath.Join(RootDir, SubDir), 0755)
	if err != nil {
		logrus.Error(err)
		return
	}
	for r, bt := range bts {
		fn := fmt.Sprintf("content_%02d_%02.1f_%s.json", bt.Threshold.Value, bt.Threshold.SampleTime, bt.ID)
		fmt.Println(fn)
		filename := filepath.Join(RootDir, SubDir, fn)
		common.SaveJSON(bt, filename, r+1)
	}
}

func FileNames() []string {

	return []string{
		// filepath.Join(RootDir, SubDir, "content_10_6.0_20240514091854.json"),
		// filepath.Join(RootDir, SubDir, "content_10_6.0_20240514092011.json"),
		// filepath.Join(RootDir, SubDir, "content_10_6.0_20240514092128.json"),
		// filepath.Join(RootDir, SubDir, "content_10_6.0_20240514092244.json"),
		// filepath.Join(RootDir, SubDir, "content_10_6.0_20240514092402.json"),
		// filepath.Join(RootDir, SubDir, "content_10_6.0_20240514095110.json"),
		// filepath.Join(RootDir, SubDir, "content_10_6.0_20240514095228.json"),
		// filepath.Join(RootDir, SubDir, "content_10_6.0_20240514095349.json"),
		// filepath.Join(RootDir, SubDir, "content_10_6.0_20240514095509.json"),
		// filepath.Join(RootDir, SubDir, "content_10_6.0_20240514095628.json"),
		// filepath.Join(RootDir, SubDir, "content_10_6.0_20240514100447.json"),
		// filepath.Join(RootDir, SubDir, "content_10_6.0_20240514100604.json"),
		// filepath.Join(RootDir, SubDir, "content_10_6.0_20240514100719.json"),
		// filepath.Join(RootDir, SubDir, "content_10_6.0_20240514100834.json"),
		// filepath.Join(RootDir, SubDir, "content_10_6.0_20240514100948.json"),
		// filepath.Join(RootDir, SubDir, "content_10_6.0_20240514101101.json"),
		// filepath.Join(RootDir, SubDir, "content_10_6.0_20240514101214.json"),
		// filepath.Join(RootDir, SubDir, "content_10_6.0_20240514101327.json"),
		// filepath.Join(RootDir, SubDir, "content_10_6.0_20240514101440.json"),
		// filepath.Join(RootDir, SubDir, "content_10_6.0_20240514101554.json"),
		// filepath.Join(RootDir, SubDir, "content_10_6.0_20240514102938.json"),
		// filepath.Join(RootDir, SubDir, "content_10_6.0_20240514103056.json"),
		// filepath.Join(RootDir, SubDir, "content_10_6.0_20240514103214.json"),
		// filepath.Join(RootDir, SubDir, "content_10_6.0_20240514103332.json"),
		// filepath.Join(RootDir, SubDir, "content_10_6.0_20240514103450.json"),
		// filepath.Join(RootDir, SubDir, "content_10_6.0_20240514103932.json"),
		// filepath.Join(RootDir, SubDir, "content_10_6.0_20240514104048.json"),
		// filepath.Join(RootDir, SubDir, "content_10_6.0_20240514104203.json"),
		// filepath.Join(RootDir, SubDir, "content_10_6.0_20240514104319.json"),
		// filepath.Join(RootDir, SubDir, "content_10_6.0_20240514104436.json"),
		// filepath.Join(RootDir, SubDir, "content_11_7.0_20240514104930.json"),
		// filepath.Join(RootDir, SubDir, "content_11_7.0_20240514105102.json"),
		// filepath.Join(RootDir, SubDir, "content_11_7.0_20240514105235.json"),
		// filepath.Join(RootDir, SubDir, "content_11_7.0_20240514105410.json"),
		// filepath.Join(RootDir, SubDir, "content_11_7.0_20240514105542.json"),
		// filepath.Join(RootDir, SubDir, "content_11_7.0_20240514105659.json"),
		// filepath.Join(RootDir, SubDir, "content_11_7.0_20240514105835.json"),
		// filepath.Join(RootDir, SubDir, "content_11_7.0_20240514110011.json"),
		// filepath.Join(RootDir, SubDir, "content_11_7.0_20240514110147.json"),
		// filepath.Join(RootDir, SubDir, "content_11_7.0_20240514110325.json"),
		// filepath.Join(RootDir, SubDir, "content_08_5.0_20240514111423.json"),
		// filepath.Join(RootDir, SubDir, "content_08_5.0_20240514111615.json"),
		// filepath.Join(RootDir, SubDir, "content_08_5.0_20240514111807.json"),
		// filepath.Join(RootDir, SubDir, "content_08_5.0_20240514112000.json"),
		// filepath.Join(RootDir, SubDir, "content_08_5.0_20240514112152.json"),
		// filepath.Join(RootDir, SubDir, "content_08_5.0_20240514112705.json"),
		// filepath.Join(RootDir, SubDir, "content_08_5.0_20240514113044.json"),
		// filepath.Join(RootDir, SubDir, "content_08_5.0_20240514113044.json"),
		// filepath.Join(RootDir, SubDir, "content_08_5.0_20240514113231.json"),
		// filepath.Join(RootDir, SubDir, "content_08_5.0_20240514113418.json"),
		// filepath.Join(RootDir, SubDir, "content_08_5.0_20240514114548.json"),
		// filepath.Join(RootDir, SubDir, "content_08_5.0_20240514114739.json"),
		// filepath.Join(RootDir, SubDir, "content_08_5.0_20240514114929.json"),
		// filepath.Join(RootDir, SubDir, "content_08_5.0_20240514115118.json"),
		// filepath.Join(RootDir, SubDir, "content_08_5.0_20240514115305.json"),
		// filepath.Join(RootDir, SubDir, "content_09_5.0_20240514141354.json"),
		// filepath.Join(RootDir, SubDir, "content_09_5.0_20240514141452.json"),
		// filepath.Join(RootDir, SubDir, "content_09_5.0_20240514141548.json"),
		// filepath.Join(RootDir, SubDir, "content_09_5.0_20240514141644.json"),
		filepath.Join(RootDir, SubDir, "content_09_5.0_20240514141739.json"),
	}

	// files, _ := os.ReadDir(filepath.Join(RootDir, SubDir))
	// filenames := []string{}
	// for _, f := range files {
	// 	if strings.Contains(f.Name(), ".json") {
	// 		filenames = append(filenames, filepath.Join(RootDir, SubDir, f.Name()))
	// 	}
	// }
	// return filenames
}

func Test_readBackTesting(t *testing.T) {
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

	tops := ar.List.WithRange(1, 1)
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
}

func Test_groupNumbers(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Group Index")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()

	GroupCount := 200

	groupMapping := ar.GroupIndexMapping(GroupCount)

	result := map[int]FTN{}
	for _, v := range ar.RevList {
		gidx := groupMapping[v.Key()]
		result[gidx] = v
	}

	bytes, _ := json.Marshal(result)

	common.Save(string(bytes), fmt.Sprintf("./gendata/Group/topGroupedStatic_%d_%s.json", GroupCount, time.Now().Format("20060102150405")), 0)

}

func Test_FTNGroup(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Group Index")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	GroupCount := 500
	gftn := NewFTNGroup(GroupCount, ar.Combinations, ar.RevList)
	common.Save(string(gftn.Presentation()), fmt.Sprintf("./gendata/Group/ReportTopGroupedStatic_%d_%s.json", GroupCount, time.Now().Format("20060102150405")), 0)
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

// func Test_findTest(t *testing.T) {
// 	samples := [][]int{{1, 2, 3, 4, 5}, {1, 2, 3, 4, 6}, {1, 2, 3, 4, 7}, {1, 2, 3, 4, 8}, {1, 2, 3, 4, 9}}
// 	t := []int{1, 2, 3, 4, 5}
// 	if Find(sample,)

// }

func Test_backTesting(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Top Price Taken Time")
	pickupsFile := ""
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	result := ar.List.featureBackTesting()
	for k, v := range result {
		pickupsFile = pickupsFile + fmt.Sprintf("%v\n:%v\n", k, v.Presentation())
	}
	pickupsFile = pickupsFile + fmt.Sprintln("")
	pickupsFile = pickupsFile + fmt.Sprintln("")
	pickupsFile = pickupsFile + fmt.Sprintln("list out")
	count := 0
	for k, v := range result {
		if len(v) > 0 {
			count++
			pickupsFile = pickupsFile + fmt.Sprintf("%v\n:%v\n", k, v.Presentation())
		}
	}
	pickupsFile = pickupsFile + fmt.Sprintf("match : %d\n", count)
	pickupsFile = pickupsFile + fmt.Sprintf("match percent : %.2f\n", float32(count)/float32(len(ar.List))*float32(100))
	lp := filepath.Join(RootDir, SubDir, fmt.Sprintf("backtesting%s.txt", time.Now().Format("20060102150405")))
	common.Save(pickupsFile, lp, 0)
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
