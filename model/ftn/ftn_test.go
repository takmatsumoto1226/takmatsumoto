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

	//
	// df.DistableFilters([]int{df.FilterOddCount, df.FilterEvenCount})
	df.DistableFilters([]int{df.FilterTenGroupOddCount, df.FilterTenGroupEvenCount})
	th := interf.Threshold{
		Round:      5,
		Value:      11,
		SampleTime: 6,
		Sample:     len(ar.Combinations),
		Interval: interf.Interval{
			Index:  1,
			Length: 100,
		},
		Smart: interf.Smart{
			Enable: true,
			Type:   interf.RangeTypeSpecStartRange,
		},
		Randomer: 1,
	}

	ar.JSONGenerateTopPriceNumber(th)

	// ar.BackTest()
}

func FileNames() []string {

	return []string{
		// filepath.Join(RootDir, SubDir, "content20240510090648"),
		// filepath.Join(RootDir, SubDir, "content20240510090852"),
		// filepath.Join(RootDir, SubDir, "content20240510091041"),
		// filepath.Join(RootDir, SubDir, "content20240510091434"),
		// filepath.Join(RootDir, SubDir, "content20240510094727"),
		// filepath.Join(RootDir, SubDir, "content20240510094812"),
		// filepath.Join(RootDir, SubDir, "content20240510094903"),
		// filepath.Join(RootDir, SubDir, "content20240510094950"),
		// filepath.Join(RootDir, SubDir, "content20240510095035"),
		// filepath.Join(RootDir, SubDir, "content20240510095121"),
		// filepath.Join(RootDir, SubDir, "content20240510095207"),
		// filepath.Join(RootDir, SubDir, "content20240510095253"),
		// filepath.Join(RootDir, SubDir, "content20240510095341"),
		// filepath.Join(RootDir, SubDir, "content20240510095426"),
		// filepath.Join(RootDir, SubDir, "content20240510095511"),
		// filepath.Join(RootDir, SubDir, "content20240510095555"),
		// filepath.Join(RootDir, SubDir, "content20240510095640"),
		// filepath.Join(RootDir, SubDir, "content20240510095724"),
		// filepath.Join(RootDir, SubDir, "content20240510095809"),
		// filepath.Join(RootDir, SubDir, "content20240510095853"),
		// filepath.Join(RootDir, SubDir, "content20240510095939"),
		// filepath.Join(RootDir, SubDir, "content20240510100024"),
		// filepath.Join(RootDir, SubDir, "content20240510100110"),
		// filepath.Join(RootDir, SubDir, "content20240510100156"),
		// filepath.Join(RootDir, SubDir, "content20240510103617"),
		// filepath.Join(RootDir, SubDir, "content20240510103703"),
		// filepath.Join(RootDir, SubDir, "content20240510103753"),
		// filepath.Join(RootDir, SubDir, "content20240510103840"),
		// filepath.Join(RootDir, SubDir, "content20240510103925"),
		// filepath.Join(RootDir, SubDir, "content20240510104613"),
		// filepath.Join(RootDir, SubDir, "content20240510104658"),
		// filepath.Join(RootDir, SubDir, "content20240510104743"),
		// filepath.Join(RootDir, SubDir, "content20240510104833"),
		// filepath.Join(RootDir, SubDir, "content20240510104917"), top
		// filepath.Join(RootDir, SubDir, "content20240510105211"),
		// filepath.Join(RootDir, SubDir, "content20240510105255"),
		// filepath.Join(RootDir, SubDir, "content20240510105338"),
		// filepath.Join(RootDir, SubDir, "content20240510105423"),
		// filepath.Join(RootDir, SubDir, "content20240510105514"),
		// filepath.Join(RootDir, SubDir, "content20240510105600"),
		// filepath.Join(RootDir, SubDir, "content20240510105645"),
		// filepath.Join(RootDir, SubDir, "content20240510105732"),
		// filepath.Join(RootDir, SubDir, "content20240510105817"),
		// filepath.Join(RootDir, SubDir, "content20240510105902"),
		// filepath.Join(RootDir, SubDir, "content20240510114548"),
		// filepath.Join(RootDir, SubDir, "content20240510114634"),
		// filepath.Join(RootDir, SubDir, "content20240510114721"),
		// filepath.Join(RootDir, SubDir, "content20240510114809"), top
		// filepath.Join(RootDir, SubDir, "content20240510114858"),
		// filepath.Join(RootDir, SubDir, "content20240510120223"),
		// filepath.Join(RootDir, SubDir, "content20240510120308"),
		// filepath.Join(RootDir, SubDir, "content20240510120353"),
		// filepath.Join(RootDir, SubDir, "content20240510120437"),
		// filepath.Join(RootDir, SubDir, "content20240510120522"), // top
		// filepath.Join(RootDir, SubDir, "content20240510130224"),
		// filepath.Join(RootDir, SubDir, "content20240510130307"), //top
		// filepath.Join(RootDir, SubDir, "content20240510130352"),
		// filepath.Join(RootDir, SubDir, "content20240510130436"),
		// filepath.Join(RootDir, SubDir, "content20240510130519"),
		// filepath.Join(RootDir, SubDir, "content20240510130603"),
		// filepath.Join(RootDir, SubDir, "content20240510130647"),
		// filepath.Join(RootDir, SubDir, "content20240510130730"),
		// filepath.Join(RootDir, SubDir, "content20240510130812"),
		// filepath.Join(RootDir, SubDir, "content20240510130855"),
		// filepath.Join(RootDir, SubDir, "content20240510130938"),
		// filepath.Join(RootDir, SubDir, "content20240510131020"),
		// filepath.Join(RootDir, SubDir, "content20240510131103"),
		// filepath.Join(RootDir, SubDir, "content20240510131146"),
		// filepath.Join(RootDir, SubDir, "content20240510131229"),
		// filepath.Join(RootDir, SubDir, "content20240510131312"),
		// filepath.Join(RootDir, SubDir, "content20240510131354"),
		// filepath.Join(RootDir, SubDir, "content20240510131438"),
		// filepath.Join(RootDir, SubDir, "content20240510131521"),
		// filepath.Join(RootDir, SubDir, "content20240510131604"),
		// filepath.Join(RootDir, SubDir, "content20240510152220"),
		// filepath.Join(RootDir, SubDir, "content20240510152306"),
		// filepath.Join(RootDir, SubDir, "content20240510152353"),
		// filepath.Join(RootDir, SubDir, "content20240510152439"),
		// filepath.Join(RootDir, SubDir, "content20240510152524"),
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

	ar.ReadJSON(FileNames()) // 20240509111846, 20240509112008
	interval := interf.Interval{Index: 0, Length: 5}
	count := 0
	// totalPickupNumber := FTNArray{}

	for _, bt := range ar.BackTest {
		fmt.Println(len(bt.ThresholdNumbers.Balls))
		for i := interval.Index; i < interval.Length; i++ {
			tops := ar.List.WithRange(i, 1)
			total := 0
			testRows := bt.PickNumbers
			for _, ftn := range tops {
				for _, pn := range testRows.Balls {
					currentPrice := ftn.AdariPrice(&pn)
					total = total + currentPrice
					if currentPrice >= df.PriceTop {
						fmt.Println(ftn.formRow())
					}
				}
			}
			// totalPickupNumber = append(totalPickupNumber, testRows.Balls...)
			if total >= df.PriceTop {
				fmt.Printf("Limit: %5d ID: %s, %d : %d, 第 %04d : %d\n\n\n", i, bt.ID, len(testRows.Balls), len(testRows.Balls)*50, i, total)
				count++
			}
		}
	}
	fmt.Println(count)

	// totalPickupNumber.ShowAll()
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
	filenames := []string{}

	ar.ReadJSON(filenames)
	featuresFTNs := FTNArray{}
	for _, bt := range ar.BackTest {
		bt.Threshold.Interval = interf.Interval{Index: 0, Length: 20}
		features := ar.List.FeatureRange(bt.Threshold)
		for _, tn := range bt.ThresholdNumbers.Balls {
			for _, l := range features {
				if l.MatchFeature(&tn) {
					featuresFTNs = append(featuresFTNs, tn)
					break
				}
			}
		}
	}
	featuresFTNs.ShowAll()
	fmt.Println(len(featuresFTNs))
}

func Test_groupNumbers(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Group Index")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	combarr := combin.Combinations(39, BallsOfFTN)
	GroupCount := 500

	result := map[int]FTN{}
	for _, ftn := range ar.List {
		for idx, v := range combarr {
			nftn := NewFTNWithInts(v)
			if nftn.IsSame(&ftn) {
				result[idx] = *nftn
				break
			}
		}
	}

	bytes, _ := json.Marshal(result)

	common.Save(string(bytes), fmt.Sprintf("./gendata/Group/topGroupedStatic_%d_%s.json", GroupCount, time.Now().Format("20060102150405")), 0)

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
