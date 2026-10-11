package main

import (
	"lottery/api"
	lotteryapi "lottery/api/lottery"
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

	stores := []*lotteryapi.Store{}
	for _, g := range lotteryapi.Games() {
		stores = append(stores, lotteryapi.NewStore(g))
	}
	reloadAll := func() {
		for _, s := range stores {
			if err := s.Reload(); err != nil {
				logrus.WithError(err).Errorf("載入%s資料失敗", s.Game.Name)
			}
		}
	}
	reloadAll()
	cmd.OnRawDataUpdated(reloadAll)

	addr := config.Config.Server.Addr
	if addr == "" {
		addr = ":8080"
	}
	e := initGINServer(stores)
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

func initGINServer(stores []*lotteryapi.Store) *gin.Engine {
	e := gin.New()
	e.Use(gin.Logger(), gin.Recovery())

	v1 := e.Group("/api/v1")
	v1.GET("/health", lotteryapi.HealthHandler(stores))

	authed := v1.Group("", api.APIKeyAuth(config.Config.Server.APIKey))
	authed.GET("/games", lotteryapi.GamesHandler(stores))
	for _, s := range stores {
		lotteryapi.Register(authed.Group("/"+s.Game.Code), s)
	}
	return e
}
