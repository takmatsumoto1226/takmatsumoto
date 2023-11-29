package config

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func Test_loadconfig(t *testing.T) {
	LoadConfig("../config.yaml")
	logrus.Info(Config)
}
