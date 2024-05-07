package ftn

import (
	"fmt"
	"lottery/model/interf"
)

type SessionData struct {
	Title string   `json:"title"`
	Balls FTNArray `json:"balls"`
}

func (sd *SessionData) Presentation() string {
	msg := ""
	msg = msg + "Title : " + sd.Title + "\n"
	msg = msg + sd.Balls.Presentation()
	return msg
}

type BackTest struct {
	ID                        string           `json:"id"`
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

	return msg
}
