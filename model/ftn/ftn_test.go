package ftn

import (
	"bytes"
	"encoding/binary"
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

	crypto_rand "crypto/rand"

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
	fmt.Println(df.Primes)
	fmt.Println(bytes.IndexByte(df.Primes, 31))
	fmt.Println(bytes.IndexByte(df.Primes, 30))
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

func Test_newFTNTest(t *testing.T) {
	elems := strings.Split("2023,1230,312,04,11,17,20,32,5114", ",")
	ftn := NewFTN(elems)
	fmt.Println(ftn)

	elems2 := strings.Split("2023,1230,312,04,11,17,20,33,5114", ",")
	ftn2 := NewFTN(elems2)
	fmt.Println(ftn2)
	fmt.Println(ftn2.CompareFeature(ftn))
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
	p := PickParam{SortType: df.Descending, Interval: 30, Whichfront: df.Normal}
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
		as.List.findNumbers(v, df.Both).Presentation()
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
	as.List.findNumbers([]string{"05", "08", "16", "24", "31"}, df.Both).Presentation()

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
	var as = FTNsManager{}
	as.Prepare()

	balls := 5
	combarr := combin.Combinations(39, balls)
	// lens := len(combarr)
	// total_sell := 50 / 50Y

	topss := []FTNArray{}
	counts := []int{}
	resultss := map[string]int{} // n round result
	// realsale := uint32(501942)
	// realsale := uint32(531942)
	realsale := uint32(len(combarr))

	// th := Threshold{Round: 20, Value: 11, SampleTime: 3, Sample: len(combarr)}
	// th := Threshold{Round: 100, Value: 11, SampleTime: 3, Sample: len(combarr)}
	th := interf.Threshold{Round: 1, Value: 14, SampleTime: 6, Sample: len(combarr), RefRange: 20}

	// lottos := as.List.WithRange(th.RefRange)
	lottos := FTNArray{}
	result := algorithm.Combinations(as.RevList[0].toStringArray(), 3)
	for _, v := range result {
		lottos = append(lottos, as.List.findNumbers(v, df.NextOnly)...)
	}

	filestr := ""
	for i := 0; i < th.Round; i++ {
		result := map[string]int{}

		var b [8]byte
		_, err := crypto_rand.Read(b[:])
		if err != nil {
			panic("cannot seed math/rand package with cryptographically secure random number generator")
		}

		rnumber := rand.New(rand.NewSource(int64(binary.LittleEndian.Uint64(b[:]))))

		for _, v := range combarr {
			balls := NewFTNWithInts(v)
			result[balls.Key()] = 0
		}
		fmt.Println(len(result))
		total := int(float32(th.Sample) * th.SampleTime)

		for i := 0; i < total; i++ {
			index := uint32(rnumber.Uint32() % realsale)
			// fmt.Println(index)
			balls := NewFTNWithInts(combarr[index])
			bK := balls.Key()
			if v, ok := result[bK]; ok {
				result[bK] = v + 1
			}

			// fmt.Println(combarr[index])
		}
		count := 0
		tops := FTNArray{}
		featuresFTNs := FTNArray{}
		for k, v := range result {
			if v > th.Value {
				filestr = filestr + fmt.Sprintf("%v:%v\n", k, v)
				arr := strings.Split(k, "_")
				ftnarr := as.List.findNumbers(arr, df.None)
				if len(ftnarr) > 0 {
					filestr = filestr + ftnarr.Presentation()
					tops = append(tops, ftnarr...)
				}

				ftn := NewFTNWithStrings(arr)
				for _, l := range lottos {
					if l.CompareFeature(ftn) {
						featuresFTNs = append(featuresFTNs, *ftn)
					}
				}

				if v2, ok := resultss[k]; ok {
					resultss[k] = v2 + v
				} else {
					resultss[k] = v
				}
				count++
			}
		}

		topss = append(topss, tops)
		counts = append(counts, count)
		filestr = filestr + fmt.Sprintf("%d TWD, %d\n", count*45, count)
		filestr = filestr + fmt.Sprintf("群 %02d, get %d Top\n", i+1, len(tops))
		filestr = filestr + fmt.Sprintf("done %02d\n", i+1)
		filestr = filestr + "\n"
		filestr = filestr + "\n"
		filestr = filestr + "\n"

		filestr = filestr + "Feature Close\n"
		if len(featuresFTNs) > 0 {
			for _, fftn := range featuresFTNs {
				filestr = filestr + fftn.formRow() + "\n"
			}
		}
	}

	miss := 0
	for i, tops := range topss {
		if len(tops) == 0 {
			miss++
		}
		filestr = filestr + fmt.Sprintf("群 %02d, 有 %d Top, D:%d TWD\n", i+1, len(tops), counts[i]*50)
	}
	filestr = filestr + fmt.Sprintf("Top Percent %.3f\n", (float32(th.Round-miss)/float32(th.Round))*100)

	filestr = filestr + "\n"
	filestr = filestr + "\n"
	filestr = filestr + "\n"
	for k, v := range resultss {
		if v >= 45 {
			filestr = filestr + fmt.Sprintf("%v:%v\n", k, v)
		}
	}

	common.Save(filestr, "content.txt")
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
