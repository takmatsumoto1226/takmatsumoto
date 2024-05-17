package ftn

import (
	"encoding/json"
	"errors"
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

// FTNsManager ...
type FTNsManager struct {
	List         FTNArray
	RevList      FTNArray
	ballsCount   map[uint]NormalizeInfo
	BackTests    []FTNBT
	Combinations [][]int
}

// LoadAllData ...
func (ar *FTNsManager) loadAllData() {
	info := config.Config.HTTP.Infos[df.InfoFTN]
	now := time.Now()

	iyear, err := strconv.Atoi(info.BaseYear)
	if err != nil {
		logrus.Error(err)
		return
	}
	var ftns FTNArray
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
			ftn := NewFTN(yd)
			if ftn.B1.Illegal() {
				continue
			}
			ftns = append(ftns, *ftn)
		}
	}
	ar.RevList = make(FTNArray, len(ftns))
	copy(ar.RevList, ftns)
	ar.List = ftns
	sort.Sort(ar.RevList)
}

func (ar *FTNsManager) Prepare() error {

	initNumberToIndex()

	// LoadAllData
	ar.loadAllData()

	ar.Combinations = combin.Combinations(ballsCountFTN, BallsOfFTN)
	ar.BackTests = []FTNBT{}
	return nil
}

func (ar *FTNsManager) GroupIndexMapping(gc int) map[string]int {
	groupMapping := map[string]int{}
	for idx, v := range ar.Combinations {
		nftn := NewFTNWithInts(v)
		groupMapping[nftn.Key()] = idx / gc
	}
	return groupMapping
}

func (ar *FTNsManager) intervalBallsCountStatic(params PickParams) {
	if len(params) == 0 {
		logrus.Error(errors.New("不可沒設定interval"))
		return
	}
	for _, p := range params {
		if p.SortType == df.Descending {
			ar.ballsCount = ar.RevList[:p.Interval].intervalBallsCountStatic(p)
		} else if p.SortType == df.Ascending {
			ar.ballsCount = ar.List[:p.Interval].intervalBallsCountStatic(p)
		} else {
			logrus.Error("必須指定型態")
			break
		}
	}
}

func (ar *FTNsManager) Picknumber(params PickParams) map[string]BallsInfo {
	for _, p := range params {
		norball := ar.ballsCount[p.Interval]
		if p.Whichfront == df.Biggerfront {
			sort.Sort(sort.Reverse(norball.AppearBalls))
		} else if p.Whichfront == df.Smallfront {
			sort.Sort(norball.AppearBalls)
		} else {

		}

		if len(norball.AppearBalls) > 5 {
			pool := BallsInfo{}
			// logrus.Infof("%s %s %s %s %s", balls[blIdx].Ball.Number, balls[blIdx-1].Ball.Number, balls[blIdx-2].Ball.Number, balls[blIdx-3].Ball.Number, balls[blIdx-4].Ball.Number)
			for _, ball := range norball.AppearBalls {

				pool = append(pool, ball)
			}
			ballPools[p.GetKey()] = pool
		} else {
			logrus.Error("數字不足")
			return nil
		}
	}
	return ballPools
}

func (ar *FTNsManager) JSONGenerateTopPriceNumber(th interf.Threshold) []FTNBT {
	common.SetRandomGenerator(th.Randomer)
	bts := []FTNBT{}
	for r := 0; r < th.Round; r++ {
		bt := NewBackTest(time.Now(), th)
		result := map[string]int{}

		for _, v := range ar.Combinations {
			balls := NewFTNWithInts(v)
			result[balls.Key()] = 0
		}
		total := int(float32(th.Sample) * th.SampleTime)

		for i := 0; i < total; i++ {
			index := common.RandomNuber() % uint64(th.Sample)
			balls := NewFTNWithInts(ar.Combinations[index])
			if v, ok := result[balls.Key()]; ok {
				result[balls.Key()] = v + 1
			}
		}

		features := ar.List.FeatureRange(th)
		bt.Features.Title = "features row\n\n\n"
		bt.Features.Balls = features

		count := 0
		tops := FTNArray{}
		threadholdNumbers := FTNArray{}
		featuresFTNs := FTNArray{}
		for k, v := range result {
			if v > th.Value {
				thNumb := NewFTNWithStrings(strings.Split(k, "_"))
				threadholdNumbers = append(threadholdNumbers, *thNumb)
				ftnarr := ar.List.findNumbers(thNumb.toStringArray(), df.None)
				if len(ftnarr) > 0 {
					tops = append(tops, ftnarr...)
				}

				for _, l := range features {
					if l.MatchFeature(thNumb) {
						featuresFTNs = append(featuresFTNs, *thNumb)
						break
					}
				}
				count++
			}
		}
		bt.ThresholdNumbers.Title = "Thread Hold Numbers"
		bt.ThresholdNumbers.Balls = threadholdNumbers
		bt.ThreadHoldCount = len(threadholdNumbers)

		bt.HistoryTopCount = len(tops)

		bt.PickNumbers.Title = "Feature Close"
		bt.PickNumbers.Balls = featuresFTNs
		bt.PickupCount = len(featuresFTNs)
		bt.NumbersHistoryTopsPercent = float32(len(tops)) / float32(count) * 100.0

		// exclude tops
		pures := FTNArray{}
		for _, fFTN := range featuresFTNs {
			for _, f := range ar.List {
				if !fFTN.MatchFeature(&f) {
					pures = append(pures, fFTN)
					break
				}
			}
		}

		bt.ExcludeTops.Title = "Pures"
		bt.ExcludeTops.Balls = pures

		bts = append(bts, *bt)
	}

	return bts
}

func (ar *FTNsManager) ReadJSON(filenames []string) {
	for _, filename := range filenames {
		if !strings.Contains(filename, ".json") {
			filename = filename + ".json"
		}
		bt := FTNBT{}
		file, err := os.ReadFile(filename)
		if err != nil {
			logrus.Error(err)
			continue
		}
		if err := json.Unmarshal(file, &bt); err != nil {
			logrus.Error(err)
			continue
		}
		ar.BackTests = append(ar.BackTests, bt)
	}
}

func (ar *FTNsManager) BackTestingReports(filenames []string) {
	ar.ReadJSON(filenames)
	for _, bt := range ar.BackTests {
		bt.Report()
	}
}

func (ar *FTNsManager) DoBackTesting(filenames []string) {
	ar.ReadJSON(filenames)
	top := ar.RevList.GetFTNWithDate("20240516")
	for _, bt := range ar.BackTests {
		bt.DoBacktesting(top)
	}
}

func (ar *FTNsManager) Predictions(filenames []string) {
	ar.ReadJSON(filenames)

	interval := interf.Interval{Index: 0, Length: 5}
	count := 0

	for _, bt := range ar.BackTests {
		for i := interval.Index; i < interval.Length; i++ {
			tops := ar.List.WithRange(i, 1)

			backtests := bt.PickNumbers
			bt.ThresholdNumbers.DoPrediction(tops)

			total := bt.PickNumbers.DoPrediction(tops)
			if total >= PriceTop {
				fmt.Printf("ID: %s, %d : %d, 第 %04d : %d\n\n\n", bt.ID, len(backtests.Balls), len(backtests.Balls)*50, i, total)
				if len(tops) > 0 {
					fmt.Println(tops[0].formRow())
				}

				count++
			}
		}
		bt.Save()
	}
	fmt.Println(count)
}

func (ar *FTNsManager) RandomInterval() interf.Interval {
	interval := interf.Interval{}

	return interval
}

func (ar *FTNsManager) GroupZero(arr FTNArray) {
	GroupCount := 200

	groupMapping := ar.GroupIndexMapping(GroupCount)

	result := map[int]FTN{}
	for _, v := range ar.RevList {
		gidx := groupMapping[v.Key()]
		result[gidx] = v
	}
}

func (ar *FTNsManager) FinalPick(filenames []string) {
	if len(filenames) == 0 {
		fmt.Errorf("no file names !!!\n")
		return
	}
	filterPick := FTNArray{}
	ar.ReadJSON(filenames)

	group := NewFTNGroup(200, ar.Combinations, ar.RevList)

	for _, bt := range ar.BackTests {
		for _, ftn := range bt.PickNumbers.Balls {
			if _, gcount := group.FindGroupIndex(ftn); gcount == 0 {
				filterPick = append(filterPick, ftn)
			}
		}
	}
}
