package ftn

import (
	"errors"
	"fmt"
	"lottery/config"
	"lottery/csv"
	"sort"
	"strconv"
	"time"

	"lottery/model/common"

	"lottery/model/df"

	"github.com/sirupsen/logrus"
)

var ballPools = map[string]Balls{}

// PickParams ...
type PickParams []PickParam

// PickParam ...
type PickParam struct {
	Key        string
	SortType   uint
	Interval   uint
	Whichfront uint
	Spliter    uint
	Hot        uint
}

// GetKey 取得key
func (p *PickParam) GetKey() string {
	return fmt.Sprintf("%d_%d_%d", p.SortType, p.Interval, p.Whichfront)
}

const ballsCountFTN = 39

// Balls ...
type Balls []Ball

// Ball 球
type Ball struct {
	Number string
}

// BallsCount ...
type BallsCount []BallInfo
type NormalizeInfo struct {
	NorBalls BallsCount
	Param    PickParam
}

// Len ...
func (fa BallsCount) Len() int {
	return len(fa)
}

// Less ...
func (fa BallsCount) Less(i, j int) bool {
	return fa[i].Count < fa[j].Count
}

// Swap swaps the elements with indexes i and j.
func (fa BallsCount) Swap(i, j int) {
	fa[i], fa[j] = fa[j], fa[i]
}

// BallCount ...
type BallInfo struct {
	Count uint
	Ball  Ball
}

/*
FTN
*/
const (
	arrIdxYear = iota
	arrIdxMonthDay
	arrIdxLIdx
	arrIdxB1
	arrIdxB2
	arrIdxB3
	arrIdxB4
	arrIdxB5
	arrIdxTIdx
	arrFTNCount
)

// FTNArray ...
type FTNArray []FTN

// FTNsManager ...
type FTNsManager struct {
	List       FTNArray
	RevList    FTNArray
	ballsCount map[uint]NormalizeInfo
}

// Len ...
func (fa FTNArray) Len() int {
	return len(fa)
}

// Less ...
func (fa FTNArray) Less(i, j int) bool {
	ii, _ := strconv.Atoi(fa[i].TIdx)
	jj, _ := strconv.Atoi(fa[j].TIdx)
	return ii > jj
}

// Swap swaps the elements with indexes i and j.
func (fa FTNArray) Swap(i, j int) {
	fa[i], fa[j] = fa[j], fa[i]
}
func (fa FTNArray) List() {
	for _, ftn := range fa {
		ftn.formRow()
	}
}

var numberToIndex = map[string]int{}

func initNumberToIndex() {
	for i := 0; i < ballsCountFTN; i++ {
		key := fmt.Sprintf("%02d", i+1)
		numberToIndex[key] = i
	}
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
			b := BallInfo{Count: count, Ball: Ball{fmt.Sprintf("%02d", i+1)}}
			arr = append(arr, b)
		}
		ar.ballsCount[p.Interval] = NormalizeInfo{NorBalls: arr, Param: p}
		sort.Sort(ar.ballsCount[p.Interval].NorBalls)
	}
}

func (ar *FTNsManager) picknumber(params PickParams) map[string]Balls {
	for _, p := range params {
		norball := ar.ballsCount[p.Interval]
		if p.Whichfront == df.Biggerfront {
			sort.Sort(sort.Reverse(norball.NorBalls))
		}

		logrus.Info(norball.NorBalls)
		if len(norball.NorBalls) > 5 {
			pool := Balls{}
			// logrus.Infof("%s %s %s %s %s", balls[blIdx].Ball.Number, balls[blIdx-1].Ball.Number, balls[blIdx-2].Ball.Number, balls[blIdx-3].Ball.Number, balls[blIdx-4].Ball.Number)
			for _, ball := range norball.NorBalls {
				if ball.Count > 0 {
					pool = append(pool, Ball{Number: ball.Ball.Number})
				}
			}
			ballPools[p.GetKey()] = pool
		} else {
			logrus.Error("數字不足")
			return nil
		}
	}
	return ballPools
}

func (ar *FTNsManager) list() {
	ar.List.List()
}

// FTN ...
type FTN struct {
	Year     string
	MonthDay string
	LIdx     string
	B1       string
	B2       string
	B3       string
	B4       string
	B5       string
	TIdx     string
}

func (fa *FTN) toStringArray() []string {
	return []string{fa.B1, fa.B2, fa.B3, fa.B4, fa.B5}
}

func (fa *FTN) formRow() {
	rowmsg := fmt.Sprintf("%s|", fa.Year)
	rowmsg = rowmsg + fmt.Sprintf("%s|", fa.MonthDay)
	iB1, _ := strconv.Atoi(fa.B1)
	iB2, _ := strconv.Atoi(fa.B2)
	iB3, _ := strconv.Atoi(fa.B3)
	iB4, _ := strconv.Atoi(fa.B4)
	iB5, _ := strconv.Atoi(fa.B5)
	ballarr := []int{iB1, iB2, iB3, iB4, iB5}
	bi := 0
	for i := 1; i <= ballsCountFTN; i++ {
		if ballarr[bi] == i {
			rowmsg = rowmsg + fmt.Sprintf("%02d|", ballarr[bi])
			if bi < 4 {
				bi++
			}
		} else {
			rowmsg = rowmsg + "  |"
		}
	}
	fmt.Println(rowmsg)
}

// NewFTN ...
func NewFTN(arr []string) *FTN {
	if len(arr) == arrFTNCount {
		return &FTN{arr[arrIdxYear], arr[arrIdxMonthDay], arr[arrIdxLIdx], arr[arrIdxB1], arr[arrIdxB2], arr[arrIdxB3], arr[arrIdxB4], arr[arrIdxB5], arr[arrIdxTIdx]}
	}
	logrus.Error("FTN 資料格式錯誤")
	return Empty()
}

func Empty() *FTN {
	return &FTN{"====", "====", "==", "==", "==", "==", "==", "==", "=="}
}

func AdjacentNumberRecordCount() error {
	return nil
}

type Options struct {
	next bool
}
