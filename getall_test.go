package main

import (
	"fmt"
	"lottery/config"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func Test_getAll(t *testing.T) {
	config.LoadConfig("./config.yaml")
	getAll()
}

func Test_getMondhLastDay(t *testing.T) {
	now := time.Now()
	currentYear, _, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, 2, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	fmt.Println(firstOfMonth)
	fmt.Println(lastOfMonth)
}
