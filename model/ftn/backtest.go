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
const SubDir = "20240516"

type SessionData struct {
	Title          string   `json:"title"`
	Balls          FTNArray `json:"balls"`
	TopMatch       FTNArray `json:"topmatch"`
	PredictionTops FTNArray `json:"predictiontops"`
	Price          int      `json:"price"`
}

func (sd *SessionData) append(ftn FTN) {
	sd.Balls = append(sd.Balls, ftn)
}

func (sd *SessionData) appendTop(top FTN) {
	sd.TopMatch = append(sd.TopMatch, top)
}

func (sd *SessionData) appendPTop(top FTN) {
	sd.PredictionTops = append(sd.PredictionTops, top)
}

func (sd *SessionData) DoBT(top FTN) {
	sd.TopMatch = FTNArray{}
	for _, pn := range sd.Balls {
		price := top.AdariPrice(&pn)
		sd.Price = sd.Price + price
		if price >= PriceTop {
			sd.appendTop(pn)
		}
	}
}

func (sd *SessionData) DoPrediction(ftns FTNArray) int {
	sd.PredictionTops = FTNArray{}
	total := 0
	for _, ftn := range ftns {
		for _, pn := range sd.Balls {
			price := ftn.AdariPrice(&pn)
			total = total + price
			if price >= PriceTop {
				sd.appendPTop(pn)
			}
		}
	}
	return total
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

func (bt *BackTest) Summery() string {
	msg := "ID : " + bt.ID + "\n"
	msg = msg + fmt.Sprintf("Tops:%d, EnumCount:%d, Pickup:%d\n", bt.HistoryTopCount, bt.ThreadHoldCount, bt.PickupCount)
	msg = msg + fmt.Sprintf("%f\n", bt.NumbersHistoryTopsPercent)
	msg = msg + bt.Threshold.Presentation()
	return msg
}

func (bt *BackTest) DoBacktesting(top FTN) {
	bt.ThresholdNumbers.DoBT(top)
	bt.PickNumbers.DoBT(top)
	bt.Save()
}

// func (bt *BackTest) DoPrediction(ftns FTNArray) {
// 	for _, ftn := range ftns {
// 		bt.ThresholdNumbers.DoPrediction(ftn)
// 		bt.PickNumbers.DoPrediction(ftn)
// 	}
// }

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
