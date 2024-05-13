package pw

import "lottery/model/interf"

const RootDir = "./gendata"
const SubDir = "20240513"

type SessionData struct {
	Title string    `json:"title"`
	Balls PowerList `json:"balls"`
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
