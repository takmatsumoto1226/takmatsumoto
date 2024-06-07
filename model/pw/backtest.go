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
const SubDir = "20240513"

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
		bt.FileName = fmt.Sprintf("content_%02d_%02.1f_%s.json", bt.Threshold.Value, bt.Threshold.SampleTime, bt.ID)
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
