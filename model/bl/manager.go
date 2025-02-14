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

type BigLotterysManager struct {
	List          BLList
	RevList       BLList
	ballsCount    map[uint]NormalizeInfo
	numberToIndex map[string]int
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
	var bll BLList
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
			bl := NewBL(yd)
			if bl == nil || bl.B1.Number == "" || bl.B1.Number == "00" {
				continue
			}
			bll = append(bll, *bl)
		}
	}
	ar.RevList = make(BLList, len(bll))
	copy(ar.RevList, bll)
	ar.List = bll
	sort.Sort(ar.List)
}
