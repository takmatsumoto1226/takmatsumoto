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
	List    PowerList
	RevList PowerList
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
	ar.RevList = make(PowerList, len(ftns))
	copy(ar.RevList, ftns)
	ar.List = ftns
	sort.Sort(ar.RevList)
}

func (mgr *PowerManager) JSONGenerateTopPriceNumber(th interf.Threshold) []BackTest {
	common.SetRandomGenerator(th.Randomer)
	bts := []BackTest{}

	for r := 0; r < th.Round; r++ {
		bt := BackTest{}
		result := map[string]int{}
		bt.Threshold = th
		for _, v := range mgr.Combinations {
			balls := NewPowerWithInts(v)
			result[balls.Key()] = 0
		}

		featureMatchs := PowerList{}
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
						featureMatchs = append(featureMatchs, *pw)
						break
					}
				}
			}
		}

		bt.ThreadHoldCount = len(bt.ThresholdNumbers.Balls)

		bt.PickNumbers.Title = "Pickup Numbers"
		bt.PickNumbers.Balls = featureMatchs
		bt.PickupCount = len(featureMatchs)

		bt.HistoryTopCount = len(bt.HistoryTopsMatch.Balls)

		bt.ID = time.Now().Format("20060102150405")
		bts = append(bts, bt)
	}
	return bts
}

func (mgr *PowerManager) Predictions() {
	interval := interf.Interval{Index: 0, Length: 1}
	count := 0

	for _, bt := range mgr.BackTests {
		for i := interval.Index; i < interval.Length; i++ {
			tops := mgr.List.WithRange(i, 1)
			total := 0
			testRows := bt.PickNumbers
			for _, ftn := range tops {
				for _, pn := range testRows.Balls {
					currentPrice := ftn.AdariPrice(&pn)
					total = total + currentPrice
					if currentPrice >= PriceTop {
						fmt.Println(ftn.formRow())
					}
				}
			}
			if total >= PriceTop {
				fmt.Printf("Limit: %5d ID: %s, %d : %d, 第 %04d : %d\n\n\n", i, bt.ID, len(testRows.Balls), len(testRows.Balls)*50, i, total)
				count++
			}
		}
	}
	fmt.Println(count)
}

func (mgr *PowerManager) ReadJSON(filenames []string) {
	for _, filename := range filenames {
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
	latest := mgr.RevList[0]
	i := interf.Interval{Index: 1, Length: len(mgr.RevList) - 1}
	histories := mgr.List.WithRange(i.Index, i.Length)
	for _, his := range histories {
		if his.MatchFeature(&latest) {
			fmt.Println(his.formRow())
		}
	}
}
