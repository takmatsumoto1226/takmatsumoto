package main

import (
	"lottery/config"
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
	getAll()
	if err := getAllFromURL(); err != nil {
		a.NoError(err)
	}
}
