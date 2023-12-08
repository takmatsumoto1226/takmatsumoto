package ftn

import (
	"fmt"
	"lottery/algorithm"
	"lottery/config"
	"lottery/model/df"
	"sort"
	"testing"

	"github.com/sirupsen/logrus"
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
	// as.findNumbers([]string{"01", "38"}, true).Presentation()
	// as.findNumbers([]string{"02", "36", "37"}, true).Presentation()
	// as.List.Head()

	// as.findNumbers([]string{"11", "22"}, true).Presentation()
	// fmt.Println("")
	// fmt.Println("")
	// as.findNumbers([]string{"10", "20", "30"}, true).Presentation()
	// fmt.Println("")
	// fmt.Println("")
	// as.findNumbers([]string{"11", "22", "33"}, true).Presentation()
	// fmt.Println("")
	// fmt.Println("")
	p := PickParam{SortType: df.Descending, Interval: 20, Whichfront: df.Normal}
	as.List.ListWithRange(int(p.Interval))
	// fmt.Println("===============================================================================================================================")
	params := PickParams{
		// {SortType: descending, Interval: 360, Whichfront: biggerfront},
		// {SortType: descending, Interval: 180, Whichfront: biggerfront},
		// {SortType: descending, Interval: 60, Whichfront: biggerfront},
		// {SortType: descending, Interval: 48, Whichfront: biggerfront},
		// {SortType: descending, Interval: 36, Whichfront: biggerfront},
		// {SortType: df.Descending, Interval: 24, Whichfront: df.Biggerfront},
		p,
		// {SortType: descending, Interval: 5, Whichfront: biggerfront},
		// {SortType: descending, Interval: 2, Whichfront: biggerfront},
		// {SortType: descending, Interval: 1, Whichfront: biggerfront},
	}
	as.intervalBallsCountStatic(params)
	as.Picknumber(params)[p.GetKey()].Presentation()
	result := algorithm.Combinations(as.RevList[0].toStringArray(), 3)
	for _, v := range result {
		fmt.Println("")
		fmt.Println("")
		fmt.Printf("=================== %s ================\n", v)
		as.findNumbers(v, true).Presentation()
	}
	fmt.Println("")
	fmt.Println("")
	as.findNumbers([]string{"06", "10", "19", "30", "35"}, true).Presentation()
}

func Test_combination(t *testing.T) {
	fmt.Println(algorithm.All([]string{"09", "14", "30"}))
}

func Test_combination2(t *testing.T) {
	fmt.Println(algorithm.Combinations([]string{"09", "11", "14", "30", "35"}, 3))
}
