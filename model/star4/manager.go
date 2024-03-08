package star4

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

type Star4Manager struct {
	List    Star4s
	RevList Star4s
	info    *StaticInfo
}

func NewStar4Manager() *Star4Manager {
	return &Star4Manager{info: NewStaticInfo()}
}

func (ar *Star4Manager) Prepare() error {

	ar.loadAllData()

	ar.info.prepare()
	return nil
}

func (ar *Star4Manager) loadAllData() {
	info := config.Config.HTTP.Infos[df.Info4STAR]
	now := time.Now()

	iyear, err := strconv.Atoi(info.BaseYear)
	if err != nil {
		logrus.Error(err)
		return
	}
	var star4s Star4s
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
			bl := NewStar4(strconv.Itoa(year), yd)
			for i, _ := range bl {
				if err := bl[i].Normalize(); err != nil {
					logrus.Error(err)
				}
				if v, ok := ar.info.numberStatic[bl[i].Balls]; ok {
					ar.info.numberStatic[bl[i].Balls] = v + 1
				} else {
					ar.info.numberStatic[bl[i].Balls] = 1
				}
			}
			star4s = append(star4s, bl...)
		}
	}
	sort.Sort(star4s)
	ar.RevList = make(Star4s, len(star4s))
	copy(ar.RevList, star4s)
	ar.List = star4s

}

func (ar *Star4Manager) StaticPresentation() {
	ar.info.formRow()
}
