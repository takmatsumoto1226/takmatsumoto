package bl

import (
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

const ballsCountBigLottery = 49

// BigLottery ...
type BigLottery struct {
	Year     string
	MonthDay string
	LIdx     string
	B1       string
	B2       string
	B3       string
	B4       string
	B5       string
	B6       string
	B7       string
	TIdx     string
}

func Empty() *BigLottery {
	return &BigLottery{"====", "====", "==", "==", "==", "==", "==", "==", "==", "==", "=="}
}

func (fa *BigLottery) toStringArray() []string {
	return []string{fa.B1, fa.B2, fa.B3, fa.B4, fa.B5, fa.B5, fa.B6, fa.B7}
}

func (fa BigLottery) formRow() {
	rowmsg := fmt.Sprintf("%s|", fa.Year)
	rowmsg = rowmsg + fmt.Sprintf("%s|", fa.MonthDay)
	iB1, _ := strconv.Atoi(fa.B1)
	iB2, _ := strconv.Atoi(fa.B2)
	iB3, _ := strconv.Atoi(fa.B3)
	iB4, _ := strconv.Atoi(fa.B4)
	iB5, _ := strconv.Atoi(fa.B5)
	iB6, _ := strconv.Atoi(fa.B6)
	iB7, _ := strconv.Atoi(fa.B7)
	ballarr := []int{iB1, iB2, iB3, iB4, iB5, iB6}
	bi := 0
	for i := 1; i <= ballsCountBigLottery; i++ {
		if ballarr[bi] == i {
			rowmsg = rowmsg + fmt.Sprintf("%02d|", ballarr[bi])
			if bi < 5 {
				bi++
			}
		} else {
			rowmsg = rowmsg + "  |"
		}
	}

	rowmsg = rowmsg + fmt.Sprintf("      |%02d|", iB7)
	fmt.Println(rowmsg)
}

type BigLotteryList []BigLottery

func (fa BigLotteryList) Len() int {
	return len(fa)
}

func (fa BigLotteryList) List() {
	for _, bl := range fa {
		bl.formRow()
	}
}

// Less ...
func (fa BigLotteryList) Less(i, j int) bool {
	ii, _ := strconv.Atoi(fa[i].TIdx)
	jj, _ := strconv.Atoi(fa[j].TIdx)
	return ii > jj
}

// Swap swaps the elements with indexes i and j.
func (fa BigLotteryList) Swap(i, j int) {
	fa[i], fa[j] = fa[j], fa[i]
}

// Ball 球
type Ball struct {
	Number string
}

// PickParam ...
type PickParam struct {
	Key        string
	SortType   uint
	Interval   uint
	Whichfront uint
	Spliter    uint
	Hot        uint
}

type BallsCount []BallInfo

// BallCount ...
type BallInfo struct {
	Count uint
	Ball  Ball
}

type NormalizeInfo struct {
	NorBalls BallsCount
	Param    PickParam
}

/*
BigLottery
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
	arrIdxB6
	arrIdxB7
	arrIdxTIdx
	arrBigLCount
)

func NewBigLottery(arr []string) *BigLottery {
	if len(arr) == arrBigLCount {
		return &BigLottery{
			arr[arrIdxYear],
			arr[arrIdxMonthDay],
			arr[arrIdxLIdx],
			arr[arrIdxB1],
			arr[arrIdxB2],
			arr[arrIdxB3],
			arr[arrIdxB4],
			arr[arrIdxB5],
			arr[arrIdxB6],
			arr[arrIdxB7],
			arr[arrIdxTIdx],
		}
	}
	logrus.Error("NewBigLottery 資料格式錯誤")
	return nil
}

type BigLotterysManager struct {
	List          BigLotteryList
	RevList       BigLotteryList
	ballsCount    map[uint]NormalizeInfo
	numberToIndex map[string]int
}

func (ar *BigLotterysManager) list() {
	ar.RevList.List()
}

func (ar *BigLotterysManager) initNumberToIndex() {
	for i := 0; i < ballsCountBigLottery; i++ {
		key := fmt.Sprintf("%02d", i+1)
		ar.numberToIndex[key] = i
	}
}

func (ar *BigLotterysManager) Prepare() error {

	ar.initNumberToIndex()

	// LoadAllData
	ar.loadAllData()
	return nil
}

func (ar *BigLotterysManager) loadAllData() {
	info := config.Config.HTTP.Infos[df.Info49]
	now := time.Now()

	iyear, err := strconv.Atoi(info.BaseYear)
	if err != nil {
		logrus.Error(err)
		return
	}
	var bll BigLotteryList
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
			bl := NewBigLottery(yd)
			bll = append(bll, *bl)
		}
	}
	ar.RevList = make(BigLotteryList, len(bll))
	copy(ar.RevList, bll)
	ar.List = bll
	sort.Sort(ar.List)
}
