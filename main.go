package main

import (
	"lottery/api"
	ftnapi "lottery/api/ftn"
	"lottery/cmd"
	"lottery/config"
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

	store := ftnapi.NewStore()
	if err := store.Reload(); err != nil {
		logrus.WithError(err).Error("載入今彩539資料失敗")
	}
	cmd.OnRawDataUpdated(func() {
		if err := store.Reload(); err != nil {
			logrus.WithError(err).Error("重新載入今彩539資料失敗")
		}
	})

	addr := config.Config.Server.Addr
	if addr == "" {
		addr = ":8080"
	}
	e := initGINServer(store)
	go func() {
		if err := e.Run(addr); err != nil {
			logrus.WithError(err).Fatal("api server 啟動失敗")
		}
	}()

	// Observe signal notification
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case <-term:
		logrus.Infof("User term")
	}
}

func initGINServer(store *ftnapi.Store) *gin.Engine {
	e := gin.New()
	e.Use(gin.Logger(), gin.Recovery())

	v1 := e.Group("/api/v1")
	v1.GET("/health", ftnapi.HealthHandler(store))

	authed := v1.Group("", api.APIKeyAuth(config.Config.Server.APIKey))
	ftnapi.Register(authed.Group("/ftn"), store)
	return e
}
