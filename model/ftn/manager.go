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
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/stat/combin"
)

type IntervalCount struct {
	Appear    []int `json:"appear"`
	Disappear []int `json:"disappear"`
}

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

func (ar *FTNsManager) IntervalBallsCountStatic(params PickParams) {

	if len(params) == 0 {
		logrus.Error(errors.New("不可沒設定interval"))
		return
	}
	result := map[uint]NormalizeInfo{}
	for _, p := range params {
		if p.SortType == df.Descending {
			result[p.Interval] = ar.RevList[:p.Interval].IntervalBallsCountStatic(p)
		} else if p.SortType == df.Ascending {
			result[p.Interval] = ar.List[:p.Interval].IntervalBallsCountStatic(p)
		} else {
			logrus.Error("必須指定型態")
			break
		}
	}

	ar.ballsCount = result
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

func (ar *FTNsManager) JSONGenerateTopPriceNumber(th interf.Threshold) {
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
		bt.Features.Title = "features row"
		bt.Features.Balls = features

		count := 0
		thNumbTops := FTNArray{}
		threadholdNumbers := FTNArray{}
		featuresFTNs := FTNArray{}
		for k, v := range result {
			if v > th.Value {
				thNumb := NewFTNWithStrings(strings.Split(k, "_"))
				threadholdNumbers = append(threadholdNumbers, *thNumb)
				ftnarr := ar.List.findNumbers(thNumb.toStringArray(), df.None)
				if len(ftnarr) > 0 {
					thNumbTops = append(thNumbTops, ftnarr...)
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

		bt.HistoryTopCount = len(thNumbTops)

		bt.PickNumbers.Title = "Feature Close"
		bt.PickNumbers.Balls = featuresFTNs
		bt.PickupCount = len(featuresFTNs)
		bt.NumbersHistoryTopsPercent = float32(len(thNumbTops)) / float32(count) * 100.0

		// exclude tops
		pures := threadholdNumbers.FilterFeatureExcludes(ar.List)

		bt.ExcludeTops.Title = "Pures"
		bt.ExcludeTops.Balls = pures

		bts = append(bts, *bt)
	}

	ar.BackTests = bts
}

func B2i(b bool) int8 {
	if b {
		return 1
	}
	return 0
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

func (ar *FTNsManager) DoBackTesting(filenames []string, d string) {
	top := ar.RevList.GetFTNWithDate(d)
	ar.ReadJSON(filenames)
	for _, bt := range ar.BackTests {
		bt.DoBacktesting(top)
	}
}

func (ar *FTNsManager) Predictions(filenames []string) {
	ar.ReadJSON(filenames)

	interval := interf.Interval{Index: 0, Length: 20}
	tops := ar.List.WithRange(interval.Index, interval.Length)
	for _, bt := range ar.BackTests {
		bt.ThresholdNumbers.DoPrediction(tops)
		bt.PickNumbers.DoPrediction(tops)
		bt.Save()
	}
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

func (ar *FTNsManager) CompareLatestAndHistoryFeature(i interf.Interval) {
	latest := ar.RevList[0]
	histories := ar.List.WithRange(i.Index, i.Length)
	for _, his := range histories {
		if his.MatchFeature(&latest) {
			fmt.Println(his.formRow())
		}
	}
}

func (ar *FTNsManager) SaveBTs() {
	err := os.MkdirAll(filepath.Join(RootDir, time.Now().Format("20060102")), 0755)
	if err != nil {
		logrus.Error(err)
		return
	}
	for _, bt := range ar.BackTests {
		fmt.Println(bt.Save())
	}
}

func (ar *FTNsManager) FilterByGroupIndex(group *FTNGroup, c int) FTNArray {
	arr := FTNArray{}
	for _, bt := range ar.BackTests {
		arr = append(arr, bt.PickNumbers.Balls.FilterByGroupIndex(group, c)...)
	}
	return arr.Distinct()
}

func (ar *FTNsManager) GodPick(arr FTNArray, c int) {
	if len(arr) == 0 {
		return
	}
	common.SetRandomGenerator(1)
	picks := FTNArray{}
	for i := 0; i < c; i++ {
		a := arr[common.RandomNuber()%uint64(len(arr))]
		picks = append(picks, a)
	}
	fmt.Println(picks.Presentation())
}
