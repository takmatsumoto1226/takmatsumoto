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
	p := PickParam{SortType: df.Descending, Interval: 25, Whichfront: df.Normal}
	as.List.PresentationWithRange(int(p.Interval))
	params := PickParams{
		p,
	}
	as.intervalBallsCountStatic(params)
	as.Picknumber(params)[p.GetKey()].Presentation()
	result := algorithm.Combinations(as.RevList[0].toStringArray(), 3)
	// result := algorithm.Combinations([]string{"01", "30", "31", "38", "39"}, 3)
	for _, v := range result {
		fmt.Println("")
		fmt.Println("")
		fmt.Printf("=================== %s ================\n", v)
		as.findNumbers(v, df.Next).Presentation()
	}
	// as.findNumbers([]string{"02", "16"}, true).Presentation()
	// as.findNumbers([]string{"16", "17"}, true).Presentation()
}

func Test_combination(t *testing.T) {
	// fmt.Println(algorithm.All([]string{"09", "14", "30"}))
	// fmt.Println(Ball39())
	combarr := algorithm.Combinations(Ball39(), 3)
	for i, comb := range combarr {
		fmt.Println(comb)
		fmt.Println(i + 1)
	}
}

func Test_combination2(t *testing.T) {
	// fmt.Println(algorithm.Combinations([]string{"09", "11", "14", "30", "35"}, 3))
	config.LoadConfig("../../config.yaml")
	var as = FTNsManager{}
	as.Prepare()
	as.findNumbers([]string{"12", "23", "26"}, df.Next).Presentation()
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
	as.Picknumber(params)[p.GetKey()].Presentation()
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
	as.RevList.DTree(p).Presentation()
}
