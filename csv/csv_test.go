package csv

import (
	"lottery/config"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func Test_csv(t *testing.T) {
	a := assert.New(t)
	config.LoadConfig("../config.yaml")
	info := &config.Info{Label: "39-%d.htm", Path: "39-year", BaseYear: "2003"}
	fullpath, err := GetPath(info, 2004)
	a.NoError(err)
	logrus.Info(fullpath)
}
