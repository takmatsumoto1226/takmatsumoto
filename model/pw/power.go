package pw

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

const ballsCountPower = 38
const ballsCountPowerS2 = 8

var numberToIndex = map[string]int{}

func NewPower(arr []string) *Power {
	if len(arr) == arrPowerCount {
		return &Power{
			arr[arrIdxYear],
			arr[arrIdxMonthDay],
			arr[arrIdxLIdx],
			arr[arrIdxB1],
			arr[arrIdxB2],
			arr[arrIdxB3],
			arr[arrIdxB4],
			arr[arrIdxB5],
			arr[arrIdxB6],
			arr[arrIdxS1],
			arr[arrIdxTIdx],
		}
	}
	logrus.Error("POWER 資料格式錯誤")
	return nil
}

// Power ...
type Power struct {
	Year     string
	MonthDay string
	LIdx     string
	B1       string
	B2       string
	B3       string
	B4       string
	B5       string
	B6       string
	S1       string
	TIdx     string
}

func (fa Power) formRow() {
	rowmsg := fmt.Sprintf("%s|", fa.Year)
	rowmsg = rowmsg + fmt.Sprintf("%s|", fa.MonthDay)
	iB1, _ := strconv.Atoi(fa.B1)
	iB2, _ := strconv.Atoi(fa.B2)
	iB3, _ := strconv.Atoi(fa.B3)
	iB4, _ := strconv.Atoi(fa.B4)
	iB5, _ := strconv.Atoi(fa.B5)
	iB6, _ := strconv.Atoi(fa.B6)
	iS1, _ := strconv.Atoi(fa.S1)
	ballarr := []int{iB1, iB2, iB3, iB4, iB5, iB6}
	bi := 0
	for i := 1; i <= ballsCountPower; i++ {
		if ballarr[bi] == i {
			rowmsg = rowmsg + fmt.Sprintf("%02d|", ballarr[bi])
			if bi < 5 {
				bi++
			}
		} else {
			rowmsg = rowmsg + "  |"
		}
	}

	rowmsg = rowmsg + "      |"

	for i := 1; i <= ballsCountPowerS2; i++ {
		if iS1 == i {
			rowmsg = rowmsg + fmt.Sprintf("%02d|", iS1)
		} else {
			rowmsg = rowmsg + "  |"
		}
	}
	fmt.Println(rowmsg)
}

type PowerList []Power

func (fa PowerList) Len() int {
	return len(fa)
}

// Less ...
func (fa PowerList) Less(i, j int) bool {
	ii, _ := strconv.Atoi(fa[i].TIdx)
	jj, _ := strconv.Atoi(fa[j].TIdx)
	return ii > jj
}

// Swap swaps the elements with indexes i and j.
func (fa PowerList) Swap(i, j int) {
	fa[i], fa[j] = fa[j], fa[i]
}

func (fa PowerList) Presentation() {
	for _, f := range fa {
		f.formRow()
	}
}

/*
Power38 索引
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
	arrIdxS1
	arrIdxTIdx
	arrPowerCount
)

type PowerManager struct {
	List    PowerList
	RevList PowerList
	// ballsCount    map[uint]NormalizeInfo
	numberToIndex map[string]int
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
			ftns = append(ftns, *ftn)
		}
	}
	ar.RevList = make(PowerList, len(ftns))
	copy(ar.RevList, ftns)
	ar.List = ftns
	sort.Sort(ar.RevList)
}

func (ar *PowerManager) Prepare() error {

	initNumberToIndex()

	// LoadAllData
	ar.loadAllData()
	return nil
}

func initNumberToIndex() {
	for i := 0; i < ballsCountPower; i++ {
		key := fmt.Sprintf("%02d", i+1)
		numberToIndex[key] = i
	}
}
