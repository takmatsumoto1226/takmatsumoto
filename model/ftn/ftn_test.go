package ftn

import (
	"lottery/config"
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
	counts := PickParams{
		// {SortType: descending, Interval: 360, Whichfront: biggerfront},
		// {SortType: descending, Interval: 180, Whichfront: biggerfront},
		// {SortType: descending, Interval: 60, Whichfront: biggerfront},
		// {SortType: descending, Interval: 48, Whichfront: biggerfront},
		// {SortType: descending, Interval: 36, Whichfront: biggerfront},
		// {SortType: descending, Interval: 24, Whichfront: biggerfront},
		// {SortType: descending, Interval: 12, Whichfront: biggerfront},
		// {SortType: descending, Interval: 5, Whichfront: biggerfront},
		// {SortType: descending, Interval: 2, Whichfront: biggerfront},
		// {SortType: descending, Interval: 1, Whichfront: biggerfront},
	}
	as.intervalBallsCountStatic(counts)
	as.picknumber(counts)
	// for k, v := range ballPools {
	// 	logrus.Infof("%s:%v", k, v)
	// }
}

func Test_listLikeExecl(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = FTNsManager{}
	as.Prepare()
	as.list()
}

func Test_findnumbers(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = FTNsManager{}
	as.Prepare()
	// as.findNumbers([]string{"11", "22", "23", "38"}, true).List()
	// as.findNumbers([]string{"01", "38"}, true).List()
	// as.findNumbers([]string{"02", "36", "37"}, true).List()
	as.list()
	// fmt.Println("===============================================================================================================================")
	// fmt.Println("")
	// as.findNumbers([]string{"15", "17", "28"}, true).List()
	// fmt.Println("")
	// fmt.Println("")
	// as.findNumbers([]string{"28", "35", "38"}, true).List()
	// fmt.Println("")
	// fmt.Println("")
	// as.findNumbers([]string{"15", "35", "38"}, true).List()
	// // as.findNumbers([]string{"04", "20", "22"}, true).List()
	// fmt.Println("")
	// fmt.Println("")
	// as.findNumbers([]string{"15", "17", "35"}, true).List()

}
