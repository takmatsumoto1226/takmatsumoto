package main

import (
	"lottery/config"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// var qApp *widgets.QApplication

// TextEdit ...
// type TextEdit struct {
// 	widgets.QMainWindow

// 	textEdit *widgets.QTextEdit
// }

var RootCmd = &cobra.Command{
	Use:   "amumugw",
	Short: "Mitake Amumu Gateway Service",
	RunE:  start,
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		logrus.Errorf("無法啟用 : %s", err.Error())
		os.Exit(1)
	}
	// Observe signal notification
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case <-term:
		logrus.Infof("User term")
	}
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
		getAll()
	}
}
