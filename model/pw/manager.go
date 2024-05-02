package pw

import (
	"fmt"
	"lottery/config"
	"lottery/csv"
	"lottery/model/common"
	"lottery/model/df"
	"lottery/model/interf"
	"sort"
	"strconv"
	"strings"
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

func (mgr *PowerManager) GenerateTopPriceNumber(th interf.Threshold) {

	for r := 0; r < th.Round; r++ {
		filestr := ""

		result := map[string]int{}
		for _, v := range th.Combinations {
			balls := NewPowerWithInts(v)
			result[balls.Key()] = 0
		}

		featureMatchs := PowerList{}
		features := mgr.List.FeatureRange(th)
		total := len(th.Combinations) * int(th.SampleTime)
		for i := 0; i < total; i++ {
			index := uint64(uint64(common.RandomNuber() % uint64(len(result))))
			balls := NewPowerWithInts(th.Combinations[index])
			if v, ok := result[balls.Key()]; ok {
				result[balls.Key()] = v + 1
			}
		}

		count := 0
		for k, v := range result {
			arr := strings.Split(k, "_")
			if v > th.Value {
				filestr = filestr + fmt.Sprintf("%v:%v\n", k, v)
				findarr := mgr.List.findNumbers(arr, df.None)
				if len(findarr) > 0 {
					filestr = filestr + findarr.Presentation()
				}
				count++
				pw := NewPowerWithString(arr)
				for _, f := range features {
					if f.MatchFeature(pw) {
						filestr = filestr + "F:M 一樣\n"
						filestr = filestr + "F:" + f.formRow() + "\n"
						filestr = filestr + "M:" + pw.formRow() + "\n"
						filestr = filestr + f.Feature.String() + "\n"
						filestr = filestr + pw.Feature.String() + "\n"
						featureMatchs = append(featureMatchs, *pw)
						break
					}
				}
			}

		}
		filestr = filestr + "Feature Close\n"
		if len(featureMatchs) > 0 {
			for _, fpw := range featureMatchs {
				filestr = filestr + fpw.formRow() + "\n"
			}
		}

		filestr = filestr + fmt.Sprintf("%d coco, %d \n", count*100, count)
		filestr = filestr + fmt.Sprintf("done : %03d\n", r+1)
		filestr = filestr + th.Presentation()
		common.Save(filestr, fmt.Sprintf("./gendata/powercontent%s.txt", time.Now().Format(time.RFC3339)), r+1)
	}
}
