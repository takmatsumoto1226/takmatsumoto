package cmd

import (
	"lottery/config"
	"lottery/model/common"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "lot",
	Short: "Mitake Amumu Gateway Service",
	RunE:  start,
}

func start(cmd *cobra.Command, args []string) error {
	config.LoadConfig("./config.yaml")

	go checkRawData()

	return nil
}

var (
	rawDataMu      sync.Mutex
	rawDataUpdated func()
)

// OnRawDataUpdated 每日抓完資料後呼叫
func OnRawDataUpdated(fn func()) {
	rawDataMu.Lock()
	defer rawDataMu.Unlock()
	rawDataUpdated = fn
}

func checkRawData() {

	for {
		now := time.Now()
		// 早上 10:00
		nt := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, now.Location())
		if now.After(nt) {
			nt = nt.Add(time.Hour * 24)
		}
		logrus.Debugf("Next cron.DumpLoginLog date (%s)", nt)
		select {
		case <-time.After(nt.Sub(time.Now())):
		}
		common.GetAll()

		rawDataMu.Lock()
		fn := rawDataUpdated
		rawDataMu.Unlock()
		if fn != nil {
			fn()
		}
	}
}
