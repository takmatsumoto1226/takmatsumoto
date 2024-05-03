package pw

import (
	"fmt"
	"lottery/algorithm"
	"lottery/config"
	"lottery/model/common"
	"lottery/model/df"
	"lottery/model/interf"
	"testing"
	"time"

	"gonum.org/v1/gonum/stat/combin"
)

func Test_listLikeExecl(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = PowerManager{numberToIndex: map[string]int{}}
	as.Prepare()
	fmt.Println(as.List.Presentation())
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

	th := interf.Threshold{
		Round:      1,
		Value:      15,
		SampleTime: 8,
		Sample:     len(combarr),
		Interval: interf.Interval{
			Index:  1057,
			Length: 20,
		},
		Combinations: combarr,
		Smart: interf.Smart{
			Enable: true,
			Type:   interf.RangeTypeLatestRange,
		},
		Randomer: 1,
	}

	// th := interf.Threshold{Round: 1, Value: 26, SampleTime: 10, Sample: len(combarr)}
	common.SetRandomGenerator(1)

	as.GenerateTopPriceNumber(th)
}
