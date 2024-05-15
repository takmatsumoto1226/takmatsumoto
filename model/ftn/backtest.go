package ftn

import (
	"encoding/json"
	"fmt"
	"lottery/model/common"
	"lottery/model/interf"
	"os"
	"path/filepath"
)

const RootDir = "./gendata"
const SubDir = "20240514"

type SessionData struct {
	Title    string   `json:"title"`
	Balls    FTNArray `json:"balls"`
	TopMatch FTNArray `json:"topmatch"`
}

func (sd *SessionData) append(ftn FTN) {
	sd.Balls = append(sd.Balls, ftn)
}

func (sd *SessionData) appendTop(top FTN) {
	sd.TopMatch = append(sd.TopMatch, top)
}

func (sd *SessionData) DoBTFrom(top FTN) {
	for _, pn := range sd.Balls {
		currentPrice := top.AdariPrice(&pn)
		if currentPrice >= PriceTop {
			sd.appendTop(pn)
		}
	}
}

func (sd *SessionData) Presentation() string {
	msg := ""
	msg = msg + "Title : " + sd.Title + "\n"
	msg = msg + sd.Balls.Presentation()
	return msg
}

type BackTest struct {
	ID                        string           `json:"id"`
	FileName                  string           `json:"file_name"`
	FullPath                  string           `json:"full_path"`
	Threshold                 interf.Threshold `json:"Threshold"`
	Features                  SessionData      `json:"features"`
	ThresholdNumbers          SessionData      `json:"thresholdnumbers"`
	HistoryTopsMatch          SessionData      `json:"historytopsmatch"`
	PickNumbers               SessionData      `json:"picknumbers"`
	PickNumbersTopMatch       SessionData      `json:"picknumberstopmatch"`
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

	// msg = msg + bt.ThresholdNumbers.Presentation()
	// msg = msg + "\n\n"

	// msg = msg + bt.HistoryTopsMatch.Presentation()
	// msg = msg + "\n\n"

	msg = msg + bt.PickNumbers.Presentation()
	msg = msg + "\n\n"

	msg = msg + fmt.Sprintf("Tops:%d, EnumCount:%d, Pickup:%d\n", bt.HistoryTopCount, bt.ThreadHoldCount, bt.PickupCount)
	msg = msg + fmt.Sprintf("%f\n", bt.NumbersHistoryTopsPercent)

	msg = msg + bt.Threshold.Presentation()

	return msg
}

func (bt *BackTest) Backtesting(top FTN) {
	bt.ThresholdNumbers.DoBTFrom(top)
	bt.PickNumbers.DoBTFrom(top)
	bt.Save()
}

func (bt *BackTest) BackFilter() FTNArray {
	bfs := FTNArray{}
	for _, pn := range bt.ThresholdNumbers.Balls {
		for _, f := range bt.Features.Balls {
			if pn.MatchFeature(&f) {
				bfs = append(bfs, pn)
			}
		}
	}
	return bfs
}

func (bt *BackTest) Save() {
	if bt.FileName == "" {
		bt.FileName = fmt.Sprintf("content_%02d_%02.1f_%s.json", bt.Threshold.Value, bt.Threshold.SampleTime, bt.ID)
		bt.FullPath = filepath.Join(RootDir, SubDir, bt.FileName)
	}
	jsonString, _ := json.Marshal(bt)
	os.WriteFile(bt.FullPath, jsonString, os.ModePerm)
}

func (bt *BackTest) Report() {
	filename := fmt.Sprintf("content_%02d_%02.1f_%s_report.txt", bt.Threshold.Value, bt.Threshold.SampleTime, bt.ID)
	common.Save(bt.Presentation(), filepath.Join(RootDir, SubDir, filename), 0)
}

type RowGroup struct {
	ID   int      `json:"id"`
	Rows FTNArray `json:"rows"`
}
