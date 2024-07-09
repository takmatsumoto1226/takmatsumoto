package pw

import (
	"encoding/json"
	"fmt"
	"lottery/config"
	"lottery/csv"
	"lottery/model/common"
	"lottery/model/df"
	"lottery/model/interf"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/stat/combin"
)

const PriceTop = 200000000
const PriceSecond = 1000000
const PriceThird = 150000
const PriceFourth = 20000
const PriceFifth = 4000
const PriceSixth = 800
const PriceSeventh = 400
const PriceEigth = 200
const PriceNinth = 100

type PowerManager struct {
	List PowerList
	// ballsCount    map[uint]NormalizeInfo
	numberToIndex map[string]int
	BackTests     []BackTest
	Combinations  [][]int
}

func (ar *PowerManager) Prepare() error {

	initNumberToIndex()

	// LoadAllData
	ar.loadAllData()

	ar.Combinations = combin.Combinations(ballsCountPower, 6)
	return nil
}

// LoadAllData ...
func (ar *PowerManager) loadAllData() {
	info := config.Config.HTTP.Infos[df.InfoPOWER]
	now := time.Now()

	iyear, err := strconv.Atoi(info.BaseYear)
	if err != nil {
		logrus.Error(err)
		return
	}
	var ftns PowerList
	for year := iyear; year <= now.Year(); year++ {
		fpath, err := csv.GetPath(&info, year)
		if err != nil {
			logrus.Error(err)
		}
		yearDatas, err := common.ReadCSV(fpath)
		if err != nil {
			logrus.Error(err)
			break
		}
		for _, yd := range yearDatas {
			ftn := NewPower(yd)
			if ftn.B1.Number == "" || ftn.B1.Number == "00" {
				continue
			}
			ftns = append(ftns, *ftn)
		}
	}
	sort.Sort(ftns)
	ar.List = ftns

}

func (ar *PowerManager) FullCombination() PowerList {
	result := PowerList{}
	for _, v := range ar.Combinations {
		balls := NewPowerWithInts(v)
		result = append(result, *balls)
	}
	return result
}

func (mgr *PowerManager) JSONGenerateTopPriceNumber(th interf.Threshold) {
	common.SetRandomGenerator(th.Randomer)
	bts := []BackTest{}

	for r := 0; r < th.Round; r++ {
		bt := NewBackTest(time.Now(), th)
		result := map[string]int{}
		bt.Threshold = th
		for _, v := range mgr.Combinations {
			balls := NewPowerWithInts(v)
			result[balls.Key()] = 0
		}

		bt.Features.Title = "Power Feature List"
		bt.Features.Balls = mgr.List.FeatureRange(th)
		total := len(mgr.Combinations) * int(th.SampleTime)
		for i := 0; i < total; i++ {
			index := uint64(uint64(common.RandomNuber() % uint64(len(result))))
			balls := NewPowerWithInts(mgr.Combinations[index])
			if v, ok := result[balls.Key()]; ok {
				result[balls.Key()] = v + 1
			}
		}

		for k, v := range result {
			arr := strings.Split(k, "_")
			if v > th.Value {
				pw := NewPowerWithString(arr)
				bt.ThresholdNumbers.Balls = append(bt.ThresholdNumbers.Balls, *pw)
				findarr := mgr.List.findNumbers(arr, df.None)
				if len(findarr) > 0 {
					bt.HistoryTopsMatch.Balls = append(bt.HistoryTopsMatch.Balls, findarr...)
				}

				for _, f := range bt.Features.Balls {
					if f.MatchFeature(pw) {
						bt.PickNumbers.Balls = append(bt.PickNumbers.Balls, *pw)
						break
					}
				}
			}
		}

		bt.ThresholdNumbers.Title = "Thread Hold Numbers"
		bt.ThreadHoldCount = len(bt.ThresholdNumbers.Balls)

		bt.HistoryTopCount = len(bt.HistoryTopsMatch.Balls)

		bt.PickNumbers.Title = "Pickup Numbers"
		bt.PickupCount = len(bt.PickNumbers.Balls)

		bt.HistoryTopCount = len(bt.HistoryTopsMatch.Balls)

		pures := bt.ThresholdNumbers.Balls.FilterFeatureExcludes(mgr.List)

		bt.ExcludeTops.Title = "Pures"
		bt.ExcludeTops.Balls = pures

		bt.ID = time.Now().Format("20060102150405")
		bts = append(bts, *bt)
	}
	mgr.BackTests = bts
}

func (mgr *PowerManager) Predictions() {
	fmt.Println("Predictions")
	interval := interf.Interval{Index: 0, Length: 20}
	tops := mgr.List.WithRange(interval.Index, interval.Length)
	for _, bt := range mgr.BackTests {
		bt.ThresholdNumbers.DoPrediction(tops)
		bt.PickNumbers.DoPrediction(tops)
		bt.ExcludeTops.DoPrediction(tops)
		bt.Save()
	}
}

func (mgr *PowerManager) ReadJSON(filenames []string) {
	for i, filename := range filenames {
		fmt.Println(fmt.Sprintf("%03d : ", i+1) + "reading....." + filename)
		if !strings.Contains(filename, ".json") {
			filename = filename + ".json"
		}
		bt := BackTest{}
		file, err := os.ReadFile(filename)
		if err != nil {
			logrus.Error(err)
			continue
		}
		if err := json.Unmarshal(file, &bt); err != nil {
			continue
		}
		mgr.BackTests = append(mgr.BackTests, bt)
	}
}

func (mgr *PowerManager) BackTestingReports(filenames []string) {
	mgr.ReadJSON(filenames)
	for _, bt := range mgr.BackTests {
		bt.Report()
	}
}

func (mgr *PowerManager) CompareLatestAndHistoryFeature() {
	latest := mgr.List[0]
	i := interf.Interval{Index: 1, Length: len(mgr.List) - 1}
	histories := mgr.List.WithRange(i.Index, i.Length)
	for _, his := range histories {
		if his.MatchFeature(&latest) {
			fmt.Println(his.formRow())
		}
	}
}

func (ar *PowerManager) ListByGroupIndex(group *PWGroup, c int) PowerList {
	arr := PowerList{}
	for _, bt := range ar.BackTests {
		arr = append(arr, bt.ExcludeTops.Balls.FilterByGroupIndex(group, []int{c})...)
	}
	return arr.Distinct()
}

func GodPick(arr PowerList, c int) {
	if len(arr) == 0 {
		return
	}
	common.SetRandomGenerator(1)
	picks := PowerList{}
	for i := 0; i < c; i++ {
		a := arr[common.RandomNuber()%uint64(len(arr))]
		picks = append(picks, a)
	}
	fmt.Println(picks.Distinct().Presentation())
}

func (ar *PowerManager) SaveBTs() {
	err := os.MkdirAll(filepath.Join(RootDir, time.Now().Format("20060102")), 0755)
	if err != nil {
		logrus.Error(err)
		return
	}
	for _, bt := range ar.BackTests {
		fmt.Println(bt.Save())
	}
}
