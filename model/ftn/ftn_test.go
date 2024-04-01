package ftn

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"lottery/algorithm"
	"lottery/config"
	"lottery/model/common"
	"lottery/model/df"
	"math"
	"math/rand"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	crypto_rand "crypto/rand"

	"github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/stat/combin"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func Test_initNumberToIndex(t *testing.T) {
	initNumberToIndex()
	logrus.Info(numberToIndex)
}

func Test_loadFTNs(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = FTNsManager{}
	as.loadAllData()
	sort.Sort(as.List)
	logrus.Info(as.List)
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
	as.Picknumber(params)["0_12_2"].Presentation()
	// for k, v := range ballPools {
	// 	logrus.Infof("%s:%v", k, v)
	// }
}

func Test_listLikeExecl(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = FTNsManager{}
	as.Prepare()
	as.List.Presentation()
}

func Test_findnumbers(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = FTNsManager{}
	as.Prepare()
	fmt.Println("")
	fmt.Println("")
	p := PickParam{SortType: df.Descending, Interval: 20, Whichfront: df.Normal}
	as.List.PresentationWithRange(int(p.Interval))
	params := PickParams{
		p,
	}
	as.intervalBallsCountStatic(params)
	as.Picknumber(params)[p.GetKey()].Presentation()
	result := algorithm.Combinations(as.RevList[0].toStringArray(), 3)
	for _, v := range result {
		fmt.Println("")
		fmt.Println("")
		fmt.Printf("=================== %s ================\n", v)
		as.findNumbers(v, df.Both).Presentation()
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
	// fmt.Println(algorithm.Combinations([]string{"09", "11", "14", "30", "35"}, 3))
	// 1_7_11_24_32:38
	// 3_18_19_20_33:38
	config.LoadConfig("../../config.yaml")
	var as = FTNsManager{}
	as.Prepare()
	// as.findNumbers([]string{"01", "07", "11", "24", "32"}, df.Both).Presentation()
	// as.findNumbers([]string{"03", "18", "19", "20", "33"}, df.Both).Presentation()
	as.findNumbers([]string{"07", "30", "32", "33", "34"}, df.Both).Presentation()

}

func Test_continue(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = FTNsManager{}
	as.Prepare()
	fmt.Println("")
	fmt.Println("")

	p := PickParam{SortType: df.Descending, Interval: 60, Whichfront: df.Normal}
	as.List.PresentationWithRange(int(p.Interval))
	params := PickParams{
		p,
	}
	as.intervalBallsCountStatic(params)
	// as.Picknumber(params)[p.GetKey()].Presentation()
	fmt.Println("")
	fmt.Println("")
	as.RevList.Continue4(p).Presentation()
}

func Test_findDTree(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = FTNsManager{}
	as.Prepare()
	fmt.Println("")
	fmt.Println("")

	p := PickParam{SortType: df.Descending, Interval: 60, Whichfront: df.Normal}
	as.List.PresentationWithRange(int(p.Interval))
	params := PickParams{
		p,
	}
	as.intervalBallsCountStatic(params)
	as.Picknumber(params)[p.GetKey()].Presentation()
	fmt.Println("")
	fmt.Println("")
	as.RevList.DTree(p).Presentation()
}

func Test_findUTree(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = FTNsManager{}
	as.Prepare()
	fmt.Println("")
	fmt.Println("")

	p := PickParam{SortType: df.Descending, Interval: 60, Whichfront: df.Normal}
	as.List.PresentationWithRange(int(p.Interval))
	params := PickParams{
		p,
	}
	as.intervalBallsCountStatic(params)
	as.Picknumber(params)[p.GetKey()].Presentation()
	fmt.Println("")
	fmt.Println("")
	as.RevList.UTree(p).Presentation()
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
	as.List.PresentationWithRange(int(p.Interval))
	params := PickParams{
		p,
	}
	as.intervalBallsCountStatic(params)
	as.Picknumber(params)[p.GetKey()].Presentation()
	fmt.Println("")
	fmt.Println("")
	fmt.Printf("=================== %s ================\n", "0322")
	as.findDate("0322", df.None).Presentation()
}

func Test_random(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Top Price Taken Time")
	// Create and seed the generator.
	// Typically a non-fixed seed should be used, such as time.Now().UnixNano().
	// Using a fixed seed will produce the same output on every run.
	config.LoadConfig("../../config.yaml")
	var as = FTNsManager{}
	as.Prepare()

	balls := 5
	combarr := combin.Combinations(39, balls)
	// lens := len(combarr)

	result := map[string]int{}

	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}

	rnumber := rand.New(rand.NewSource(int64(binary.LittleEndian.Uint64(b[:]))))

	for _, v := range combarr {
		balls := NewBalls(v)
		result[balls.Key()] = 0
	}
	fmt.Println(len(result))
	total := 15545439 * 3

	// for i := 0; i < 575757000; i++ {
	for i := 0; i < total; i++ {

		index := uint32(rnumber.Uint32() % uint32(len(result)))
		// index := int(rnumber.Int31() / int32(len(combarr)))
		// index := int(rnumber.Uint32() / uint32(len(combarr)))
		// fmt.Println(index)
		// time.Sleep(time.Second)
		balls := NewBalls(combarr[index])
		if v, ok := result[balls.Key()]; ok {
			result[balls.Key()] = v + 1
		}

		// fmt.Println(combarr[index])
	}

	count := 0
	for k, v := range result {
		if v > 135 {
			fmt.Printf("%v:%v\n", k, v)
			arr := strings.Split(k, "_")
			as.findNumbers(arr, df.Both).Presentation()
			count++
		}
	}
	fmt.Printf("%d 元, %d 注\n", count*50, count)
	fmt.Println("done")
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
