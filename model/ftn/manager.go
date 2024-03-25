package ftn

import (
	"errors"
	"fmt"
	"lottery/config"
	"lottery/csv"
	"lottery/model/common"
	"lottery/model/df"
	"sort"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

// FTNsManager ...
type FTNsManager struct {
	List       FTNArray
	RevList    FTNArray
	ballsCount map[uint]NormalizeInfo
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
			if ftn.B1 == "" || ftn.B1 == "00" {
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
	return nil
}

func (ar *FTNsManager) intervalBallsCountStatic(params PickParams) {
	if len(params) == 0 {
		logrus.Error(errors.New("不可沒設定interval"))
		return
	}

	ar.ballsCount = map[uint]NormalizeInfo{}
	for _, p := range params {
		if p.Interval == 0 {
			logrus.Error(errors.New("不可指定0"))
			return
		}
		var FTNIntervalCount = [ballsCountFTN]uint{}
		var intervalFTNs = FTNArray{}
		if p.SortType == df.Descending {
			intervalFTNs = ar.RevList[:p.Interval]
		} else if p.SortType == df.Ascending {
			intervalFTNs = ar.List[:p.Interval]
		} else {
			logrus.Error("必須指定型態")
			break
		}
		for _, t := range intervalFTNs {
			FTNIntervalCount[numberToIndex[t.B1]]++
			FTNIntervalCount[numberToIndex[t.B2]]++
			FTNIntervalCount[numberToIndex[t.B3]]++
			FTNIntervalCount[numberToIndex[t.B4]]++
			FTNIntervalCount[numberToIndex[t.B5]]++
		}
		arr := BallsCount{}
		for i, count := range FTNIntervalCount {
			b := BallInfo{Count: count, Ball: Ball{fmt.Sprintf("%02d", i+1), i, i + 1}}
			arr = append(arr, b)
		}
		ar.ballsCount[p.Interval] = NormalizeInfo{NorBalls: arr, Param: p}
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
