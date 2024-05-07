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
	Pickups      FTNArray
	BackTest     []BackTest
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
	ar.BackTest = []BackTest{}
	return nil
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
			sort.Sort(sort.Reverse(norball.NorBalls))
		} else if p.Whichfront == df.Smallfront {
			sort.Sort(norball.NorBalls)
		} else {

		}

		if len(norball.NorBalls) > 5 {
			pool := BallsInfo{}
			// logrus.Infof("%s %s %s %s %s", balls[blIdx].Ball.Number, balls[blIdx-1].Ball.Number, balls[blIdx-2].Ball.Number, balls[blIdx-3].Ball.Number, balls[blIdx-4].Ball.Number)
			for _, ball := range norball.NorBalls {

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

func (ar *FTNsManager) GenerateTopPriceNumber(th interf.Threshold) {

	resultss := map[string]int{} // n round result
	common.SetRandomGenerator(th.Randomer)

	for r := 0; r < th.Round; r++ {
		filestr := ""
		result := map[string]int{}

		for _, v := range ar.Combinations {
			balls := NewFTNWithInts(v)
			result[balls.Key()] = 0
		}
		total := int(float32(th.Sample) * th.SampleTime)

		for i := 0; i < total; i++ {
			index := common.RandomNuber() % uint64(th.Sample)
			balls := NewFTNWithInts(ar.Combinations[index])
			bK := balls.Key()
			if v, ok := result[bK]; ok {
				result[bK] = v + 1
			}
		}

		features := ar.List.FeatureRange(th)
		filestr = filestr + "features row\n\n\n"
		filestr = filestr + features.Presentation()
		count := 0
		tops := FTNArray{}
		featuresFTNs := FTNArray{}
		for k, v := range result {
			if v > th.Value {
				filestr = filestr + fmt.Sprintf("%v:%v\n", k, v)
				numbersArr := strings.Split(k, "_")
				ftnarr := ar.List.findNumbers(numbersArr, df.None)
				if len(ftnarr) > 0 {
					filestr = filestr + ftnarr.Presentation()
					tops = append(tops, ftnarr...)
				}

				ftn := NewFTNWithStrings(numbersArr)
				for _, l := range features {
					if l.MatchFeature(ftn) {
						filestr = filestr + "F:M 一樣\n"
						filestr = filestr + "F:" + l.formRow() + "\n"
						filestr = filestr + "M:" + ftn.formRow() + "\n"
						filestr = filestr + l.Feature.String() + "\n"
						filestr = filestr + ftn.Feature.String() + "\n"
						featuresFTNs = append(featuresFTNs, *ftn)
						break
					}
				}

				if v2, ok := resultss[k]; ok {
					resultss[k] = v2 + v
				} else {
					resultss[k] = v
				}
				count++
			}
		}

		filestr = filestr + fmt.Sprintf("%d TWD, %d\n", count*45, count)
		filestr = filestr + fmt.Sprintf("群 %02d, get %d Top\n", r+1, len(tops))
		filestr = filestr + fmt.Sprintf("%.9f tops\n", float32(len(tops))/float32(count))
		filestr = filestr + fmt.Sprintf("done %02d\n", r+1)
		filestr = filestr + "\n"
		filestr = filestr + "\n"
		filestr = filestr + "\n"

		filestr = filestr + fmt.Sprintf("Value:%d\nRound:%d\n\n", th.Value, th.Round)
		filestr = filestr + "Feature Close : \n"
		filestr = filestr + featuresFTNs.Presentation()
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
		filestr = filestr + "Pures : \n"
		filestr = filestr + pures.Presentation()
		filestr = filestr + th.Presentation()

		ar.Pickups = append(ar.Pickups, featuresFTNs...)

		common.Save(filestr, fmt.Sprintf("./gendata/content%s.txt", time.Now().Format(time.RFC3339)), r+1)
	}

	if len(ar.Pickups) > 0 {
		pickupsFile := "Pickups:\n"
		pickupsFile = pickupsFile + ar.Pickups.Distinct().Presentation()
		// pickupsFile = pickupsFile + ar.Pickups.intervalBallsCountStatic()
		common.Save(pickupsFile, fmt.Sprintf("./gendata/pickers%s.txt", time.Now().Format(time.RFC3339)), 0)
	}
}

func (ar *FTNsManager) JSONGenerateTopPriceNumber(th interf.Threshold) []string {
	common.SetRandomGenerator(th.Randomer)
	filenames := []string{}
	for r := 0; r < th.Round; r++ {
		bt := BackTest{}
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
		pickupcount := 0
		tops := FTNArray{}
		threadholdNumbers := FTNArray{}
		featuresFTNs := FTNArray{}
		for k, v := range result {
			if v > th.Value {
				threadholdNumber := NewFTNWithStrings(strings.Split(k, "_"))
				threadholdNumbers = append(threadholdNumbers, *threadholdNumber)
				ftnarr := ar.List.findNumbers(threadholdNumber.toStringArray(), df.None)
				if len(ftnarr) > 0 {
					tops = append(tops, ftnarr...)
				}

				for _, l := range features {
					if l.MatchFeature(threadholdNumber) {
						featuresFTNs = append(featuresFTNs, *threadholdNumber)
						pickupcount++
						break
					}
				}
				count++
			}
		}
		bt.ThresholdNumbers.Title = "Thread Hold Numbers"
		bt.ThresholdNumbers.Balls = threadholdNumbers

		bt.HistoryTopsMatch.Title = "History Match Tops"
		bt.HistoryTopsMatch.Balls = tops
		bt.PickNumbers.Title = "Feature Close"
		bt.PickNumbers.Balls = featuresFTNs
		bt.ThreadHoldCount = count
		bt.PickupCount = pickupcount
		bt.ID = time.Now().Format("20060102150405")
		bt.HistoryTopCount = len(tops)
		bt.NumbersHistoryTopsPercent = float32(len(tops)) / float32(count)
		bt.Threshold = th

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

		// common.Save(filestr, fmt.Sprintf("./gendata/content%s.txt", time.Now().Format(time.RFC3339)), r+1)
		filename := fmt.Sprintf("./gendata/content%s.json", bt.ID)
		common.SaveJSON(bt, filename, r+1)
		filenames = append(filenames, filename)
	}

	return filenames
}

func (ar *FTNsManager) ReadJSON(filenames []string) {
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
		ar.BackTest = append(ar.BackTest, bt)
	}
}

func (ar *FTNsManager) RandomInterval() interf.Interval {
	interval := interf.Interval{}

	return interval
}
