package pw

import (
	"lottery/config"
	"lottery/csv"
	"lottery/model/common"
	"lottery/model/df"
	"sort"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type PowerManager struct {
	List    PowerList
	RevList PowerList
	// ballsCount    map[uint]NormalizeInfo
	numberToIndex map[string]int
}

func (ar *PowerManager) Prepare() error {

	initNumberToIndex()

	// LoadAllData
	ar.loadAllData()
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
			if ftn.B1 == "" || ftn.B1 == "00" {
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
