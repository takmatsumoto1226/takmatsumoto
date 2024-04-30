package pw

import (
	"fmt"
	"lottery/algorithm"
	"lottery/config"
	"lottery/model/common"
	"lottery/model/df"
	"lottery/model/interf"
	"math"
	"strings"
	"testing"
	"time"

	crypto_rand "crypto/rand"

	"gonum.org/v1/gonum/stat/combin"
)

func Test_listLikeExecl(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = PowerManager{numberToIndex: map[string]int{}}
	as.Prepare()
	as.List.WithRange(0, 20).Presentation()
	// fmt.Println(mt.New(mt19937.New(19650218)).Uint64())
}

func Test_findnumber(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = PowerManager{numberToIndex: map[string]int{}}
	as.Prepare()
	result := algorithm.Combinations(as.RevList[0].toStringArray(), 3)
	for _, v := range result {
		fmt.Println("")
		fmt.Println("")
		fmt.Printf("=================== %s ================\n", v)
		as.List.findNumbers(v, df.Next).Presentation()
	}
}

func maxmin(a uint32, diff, sample int) int {
	result := math.Max(float64(a), float64(diff))
	return int(math.Min(result, float64(sample-diff)))
}

func isOverRange(a uint32, diff, sample int) bool {
	min := math.Max(float64(a), float64(diff))
	max := math.Min(min, float64(sample-diff))
	return a <= uint32(min) || a >= uint32(max)
}

func Test_minmax(t *testing.T) {
	fmt.Println(maxmin(32, 10, 100))
	fmt.Println(maxmin(9, 10, 100))
	fmt.Println(maxmin(94, 10, 100))
}

func Test_newPowerTest(t *testing.T) {
	elems := strings.Split("2023,0330,26,01,12,17,26,33,38,05,1585", ",")
	power := NewPower(elems)
	fmt.Println(power)

	elems2 := strings.Split("2021,0520,40,03,14,15,22,25,36,04,1391", ",")
	// elems2 := strings.Split("2023,0330,26,01,12,17,26,33,38,05,1585", ",")
	power2 := NewPower(elems2)
	fmt.Println(power2)
	fmt.Println(power2.CompareFeature(power))
}
func Test_random(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Top Price Taken Time")
	// Create and seed the generator.
	// Typically a non-fixed seed should be used, such as time.Now().UnixNano().
	// Using a fixed seed will produce the same output on every run.
	filestr := ""
	config.LoadConfig("../../config.yaml") // 17591400
	var as = PowerManager{}
	as.Prepare()

	balls := 6
	combarr := combin.Combinations(38, balls)
	// lens := len(combarr)
	// th := interf.Threshold{Round: 1, Value: 7, SampleTime: 2, Sample: len(combarr), RefRange: 5}
	th := interf.Threshold{Round: 1, Value: 16, SampleTime: 8, Sample: len(combarr), Interval: interf.Interval{Index: 0, Length: 20}}
	// th := interf.Threshold{Round: 1, Value: 1, SampleTime: 2, Sample: 2000, RefRange: 5}
	topss := []PowerList{}

	for r := 0; r < th.Round; r++ {
		result := map[string]int{}

		var b [8]byte
		_, err := crypto_rand.Read(b[:])
		if err != nil {
			panic("cannot seed math/rand package with cryptographically secure random number generator")
		}

		// rnumber := rand.New(rand.NewSource(int64(binary.LittleEndian.Uint64(b[:]))))
		common.SetRandomGenerator(0)

		for _, v := range combarr {
			balls := NewPowerWithInts(v)
			result[balls.Key()] = 0
		}

		filestr = filestr + "result len : \n"
		filestr = filestr + fmt.Sprintln(len(result))
		total := int(float32(th.Sample) * th.SampleTime)
		for i := 0; i < total; i++ {
			index := uint32(common.RandomNuber() % uint64(len(result)))
			// for {
			// 	if isOverRange(index, 0, th.Sample) {
			// 		break
			// 	}
			// 	index = uint32(rnumber.Uint32() % uint32(len(result)))
			// }
			// balls := NewBalls(combarr[maxmin(index, 10000, th.Sample)])
			balls := NewPowerWithInts(combarr[index])
			if v, ok := result[balls.Key()]; ok {
				result[balls.Key()] = v + 1
			}
		}

		lottos := as.List.FeatureRange(th)

		zeros := PowerList{}
		zerosFeature := PowerList{}
		count := 0
		tops := PowerList{}
		featuresList := PowerList{}
		for k, v := range result {
			arr := strings.Split(k, "_")
			if v > th.Value {
				filestr = filestr + fmt.Sprintf("%v:%v\n", k, v)
				powarr := as.List.findNumbers(arr, df.None)
				if len(powarr) > 0 {
					filestr = filestr + powarr.Presentation()
					tops = append(tops, powarr...)
				}

				pwr := NewPowerWithString(arr)
				for _, l := range lottos {
					if l.CompareFeature(pwr) {
						featuresList = append(featuresList, *pwr)
					}
				}
				count++
			}

			if v == 0 {
				zeros = append(zeros, *NewPowerWithString(arr))
			}
		}

		lottos = as.List.WithRange(th.Interval.Index, th.Interval.Length)
		for _, z := range zeros {
			for _, l := range lottos {
				if l.CompareFeature(&z) {
					zerosFeature = append(zerosFeature, z)
				}
			}
		}

		topss = append(topss, tops)

		filestr = filestr + fmt.Sprintf("%d 元, %d \n", count*100, count)
		filestr = filestr + fmt.Sprintf("%d tops\n", len(tops))
		filestr = filestr + fmt.Sprintf("%.9f tops\n", float32(len(tops))/float32(count))
		filestr = filestr + fmt.Sprintln("")
		filestr = filestr + fmt.Sprintln("featuresPowers")
		filestr = filestr + featuresList.Presentation() + "\n"
		filestr = filestr + fmt.Sprintln("")
		filestr = filestr + fmt.Sprintln("")
		// filestr = filestr + fmt.Sprintln("zeros")
		// filestr = filestr + zeros.Presentation()
		filestr = filestr + fmt.Sprintln("zeros feature")
		filestr = filestr + zerosFeature.Presentation() + "\n"
		filestr = filestr + fmt.Sprintf("done : %02d\n", r+1)
		filestr = filestr + fmt.Sprintln("")
		filestr = filestr + fmt.Sprintln("")
	}

	miss := 0
	for i, tops := range topss {
		if len(tops) == 0 {
			miss++
		}
		filestr = filestr + fmt.Sprintf("群 %02d, 有 %d Top\n", i+1, len(tops))
	}
	filestr = filestr + fmt.Sprintf("Top Percent %.3f\n", (float32(th.Round-miss)/float32(th.Round))*100)
	common.Save(filestr, fmt.Sprintf("content%s.txt", time.Now().Format(time.RFC3339)), 0)
}
