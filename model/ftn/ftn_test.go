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
	"sort"
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
	ftn := NewFTNWithStrings([]string{"10", "14", "25", "31", "36"})
	ftn.ShowRow()
}

func Test_loadFTNs(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = FTNsManager{}
	as.Prepare()
	aipicks := as.List.FeatureRange(*interf.SmartPureIntervalTH(0, 20))
	sort.Sort(aipicks)
	aipicks.ShowAll()
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	distinct := aipicks.Distinct()
	distinct.ShowAll()

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
	as.List.findNumbers([]string{"01", "10", "20", "30", "36"}, df.None).ShowAll()

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
	arr.adariPrice(NewFTNWithInts([]int{2, 4, 5, 6, 7}))
}

func Test_random(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Top Price Taken Time")
	config.LoadConfig("../../config.yaml")
	var ar = FTNsManager{}
	ar.Prepare()
	combarr := combin.Combinations(39, BallsOfFTN)

	df.DistableFilters([]int{df.FilterTenGroupOddCount, df.FilterTenGroupEvenCount, df.FilterTailDigit})
	th := interf.Threshold{
		Round:      1,
		Value:      30,
		SampleTime: 20,
		Sample:     len(combarr),
		Interval: interf.Interval{
			Index:  0,
			Length: 10,
		},
		Combinations: combarr,
		Smart: interf.Smart{
			Enable: true,
			Type:   interf.RangeTypeLatest,
		},
		Randomer: 1,
	}

	ar.GenerateTopPriceNumber(th)

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
	// 	F:2007|0501|01|  |  |  |  |  |  |  |  |10|  |  |  |  |  |  |  |  |  |  |  |22|  |  |  |  |  |  |  |  |  |  |  |34|  |  |  |  |39|
	// M:||01|  |  |  |  |  |  |  |  |  |  |12|  |  |  |  |  |  |  |  |  |  |  |  |  |  |27|  |  |  |  |  |  |34|  |36|  |  |  |
	// TenGroup : [1 1 1 2 0], Odd:Even==2:3, OddTen:EvenTen===[1 0 0 1 0]:[0 1 1 1 0], DigitTail : [1 1 1 0 1 0 0 0 0 1], PrimeCount:0
	// TenGroup : [1 1 1 2 0], Odd:Even==2:3, OddTen:EvenTen===[1 0 1 0 0]:[0 1 0 2 0], DigitTail : [0 1 1 0 1 0 1 1 0 0], PrimeCount:0

	// F:2024|0427|  |  |  |  |  |  |  |  |09|  |  |  |  |  |15|  |  |  |  |  |  |22|  |  |  |  |  |  |  |  |  |  |  |34|  |36|  |  |  |
	// M:||  |  |  |  |  |  |  |08|  |10|  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |27|  |  |  |  |  |33|34|  |  |  |  |  |

	// F:2024|0426|  |  |  |  |  |  |  |  |  |10|  |  |  |14|  |  |  |  |  |  |  |  |  |  |25|  |  |  |  |  |31|  |  |  |  |36|  |  |  |
	// M:||  |  |  |  |  |  |  |  |  |  |  |  |13|  |15|  |  |  |  |  |  |22|  |  |  |  |  |  |  |  |  |  |  |34|  |36|  |  |  |
	ftn1 := NewFTNWithInts([]int{10, 14, 25, 31, 36})
	ftn2 := NewFTNWithInts([]int{13, 15, 22, 34, 36})
	if ftn1.MatchFeature(ftn2) {
		fmt.Println("一樣")
	} else {
		fmt.Println("不一樣")
	}
	fmt.Println(ftn1.Feature.String())
	fmt.Println(ftn2.Feature.String())

}
