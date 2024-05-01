package pw

import (
	"encoding/binary"
	"fmt"
	"lottery/algorithm"
	"lottery/config"
	"lottery/model/common"
	"lottery/model/df"
	"lottery/model/interf"
	"math/rand"
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
	as.List.Presentation()
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
		as.findNumbers(v, df.Next).Presentation()
	}
}

func Test_random(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Top Price Taken Time")
	// Create and seed the generator.
	// Typically a non-fixed seed should be used, such as time.Now().UnixNano().
	// Using a fixed seed will produce the same output on every run.
	config.LoadConfig("../../config.yaml") // 17591400
	var as = PowerManager{}
	as.Prepare()

	balls := 6
	combarr := combin.Combinations(38, balls)
	// lens := len(combarr)

	th := interf.Threshold{Round: 1, Value: 26, SampleTime: 10, Sample: len(combarr)}

	for r := 0; r < th.Round; r++ {

	rnumber := rand.New(rand.NewSource(int64(binary.LittleEndian.Uint64(b[:]))))

	for _, v := range combarr {
		balls := NewBalls(v)
		result[balls.Key()] = 0
	}
	fmt.Println(len(result))
	total := 2324784 * 16

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

		result := map[string]int{}

	count := 0
	for k, v := range result {
		if v > 28 {
			fmt.Printf("%v:%v\n", k, v)
			arr := strings.Split(k, "_")
			as.findNumbers(arr, df.Next).Presentation()
			count++
		}
		fmt.Println(len(result))
		total := th.Sample * int(th.SampleTime)
		for i := 0; i < total; i++ {

			index := uint32(rnumber.Uint32() % uint32(len(result)))
			balls := NewBalls(combarr[index])
			if v, ok := result[balls.Key()]; ok {
				result[balls.Key()] = v + 1
			}
		}

		count := 0
		for k, v := range result {
			if v > th.Value {
				fmt.Printf("%v:%v\n", k, v)
				arr := strings.Split(k, "_")
				as.findNumbers(arr, df.Next).Presentation()
				count++
			}
		}
		fmt.Printf("%d 元, %d \n", count*100, count)
		fmt.Printf("done : %02d\n", r)
	}
}
