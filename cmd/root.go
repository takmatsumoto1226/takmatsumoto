package cmd

import (
	"lottery/config"
	"lottery/model/common"
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
	}
}
