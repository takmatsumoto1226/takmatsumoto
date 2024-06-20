package pw

import (
	"encoding/json"
	"fmt"
	"lottery/model/common"
	"lottery/model/interf"
	"os"
	"path/filepath"
	"time"
)

const RootDir = "./gendata"
const SubDir = "20240620"
const FileNameTemplate = "content_%02d_%02.1f_%s.json"

type SessionData struct {
	Title          string    `json:"title"`
	Balls          PowerList `json:"balls"`
	TopMatch       PowerList `json:"topmatch"`
	PredictionTops PowerList `json:"predictiontops"`
	Price          int       `json:"price"`
}

func (sd *SessionData) Presentation() string {
	msg := ""
	msg = msg + "Title : " + sd.Title + "\n"
	msg = msg + sd.Balls.Presentation()
	return msg
}

func (sd *SessionData) append(ftn Power) {
	sd.Balls = append(sd.Balls, ftn)
}

func (sd *SessionData) appendTop(top Power) {
	sd.TopMatch = append(sd.TopMatch, top)
}

func (sd *SessionData) appendPTop(top Power) {
	sd.PredictionTops = append(sd.PredictionTops, top)
}

func (sd *SessionData) DoBT(top Power) {
	sd.TopMatch = PowerList{}
	for _, pn := range sd.Balls {
		price := top.AdariPrice(&pn)
		sd.Price = sd.Price + price
		if price >= PriceTop {
			sd.appendTop(pn)
		}
	}
}

func (sd *SessionData) DoPrediction(tops PowerList) int {
	sd.PredictionTops = PowerList{}
	total := 0

	for _, pn := range sd.Balls {
		for _, t := range tops {
			price := pn.AdariPrice(&t)
			total = total + price
			if price >= PriceTop {
				sd.appendPTop(t)
			}
		}
	}
	return total
}

type BackTest struct {
	ID                        string           `json:"id"`
	Date                      time.Time        `json:"date"`
	FileName                  string           `json:"file_name"`
	FullPath                  string           `json:"full_path"`
	Threshold                 interf.Threshold `json:"Threshold"`
	Features                  SessionData      `json:"features"`
	ThresholdNumbers          SessionData      `json:"thresholdnumbers"`
	HistoryTopsMatch          SessionData      `json:"historytopsmatch"`
	PickNumbers               SessionData      `json:"picknumbers"`
	ExcludeTops               SessionData      `json:"excludetops"`
	ThreadHoldCount           int              `json:"threadholdcount"`
	PickupCount               int              `json:"pickupcount"`
	HistoryTopCount           int              `json:"historytopcount"`
	NumbersHistoryTopsPercent float32          `json:"numbershistorytopspercent"`
}

func NewBackTest(date time.Time, th interf.Threshold) *BackTest {
	id := date.Format("20060102150405")
	filename := fmt.Sprintf("content_%02d_%02.1f_%s.json", th.Value, th.SampleTime, id)
	return &BackTest{
		ID:        id,
		Date:      date,
		FileName:  filename,
		FullPath:  filepath.Join(RootDir, date.Format("20060102"), filename),
		Threshold: th,
	}
}

func (bt *BackTest) Presentation() string {
	msg := "ID : " + bt.ID + "\n"
	msg = msg + bt.Features.Presentation()
	msg = msg + "\n\n"

	msg = msg + bt.ThresholdNumbers.Presentation()
	msg = msg + "\n\n"

	msg = msg + bt.HistoryTopsMatch.Presentation()
	msg = msg + "\n\n"

	msg = msg + bt.PickNumbers.Presentation()
	msg = msg + "\n\n"

	msg = msg + fmt.Sprintf("Tops:%d, EnumCount:%d\n", bt.HistoryTopCount, bt.ThreadHoldCount)
	msg = msg + fmt.Sprintf("%f\n", bt.NumbersHistoryTopsPercent)

	msg = msg + bt.Threshold.Presentation()

	return msg
}

func (bt *BackTest) Save() string {
	if bt.FileName == "" {
		bt.FileName = fmt.Sprintf(FileNameTemplate, bt.Threshold.Value, bt.Threshold.SampleTime, bt.ID)
		bt.FullPath = filepath.Join(RootDir, bt.Date.Format("0102"), bt.FileName)
	}
	jsonString, _ := json.Marshal(bt)
	os.WriteFile(bt.FullPath, jsonString, os.ModePerm)
	return bt.FileName
}

func (bt *BackTest) Report() {
	common.Save(bt.Presentation(), filepath.Join(RootDir, SubDir, fmt.Sprintf("powercontent%s_report.json", bt.ID)), 0)
}

func (bt *BackTest) Summery() string {
	msg := "ID : " + bt.ID + "\n"
	msg = msg + fmt.Sprintf("Tops:%d, EnumCount:%d, Pickup:%d\n", bt.HistoryTopCount, bt.ThreadHoldCount, bt.PickupCount)
	msg = msg + fmt.Sprintf("%f\n", bt.NumbersHistoryTopsPercent)
	msg = msg + bt.Threshold.Presentation()
	return msg
}
