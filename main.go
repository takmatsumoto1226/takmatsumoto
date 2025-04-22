package main

import (
	api "lottery/api/ftn"
	"lottery/cmd"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// var qApp *widgets.QApplication

// TextEdit ...
// type TextEdit struct {
// 	widgets.QMainWindow

// 	textEdit *widgets.QTextEdit
// }

func main() {

	if err := cmd.RootCmd.Execute(); err != nil {
		logrus.Errorf("無法啟用 : %s", err.Error())
		os.Exit(1)
	}

	e, err := initGINServer("")
	if err != nil {
		logrus.WithError(err).Logger.Exit(1)
	}
	go e.Run()

	// Observe signal notification
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case <-term:
		logrus.Infof("User term")
	}
}

func initGINServer(env string) (*gin.Engine, error) {
	e := gin.Default()
	ftnGroup := e.Group("ftn")
	ftnGroup.POST("list", api.FTNListCtx)

	gstatics := e.Group("statics")
	gstatics.POST("/numbers")
	return e, nil
}
