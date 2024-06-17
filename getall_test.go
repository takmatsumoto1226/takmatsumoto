package main

import (
	"lottery/config"
	"lottery/model/common"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func Test_getAll(t *testing.T) {
	a := assert.New(t)

	config.LoadConfig("./config.yaml")
	common.GetAll()
	if err := common.GetAllFromURL(); err != nil {
		a.NoError(err)
	}
}
