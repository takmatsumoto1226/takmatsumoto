package main

import (
	"fmt"
	"lottery/config"
	"net/url"
	"path"
	"strconv"
	"time"

	"lottery/csv"

	"github.com/sirupsen/logrus"
)

func getAll() error {
	infos := config.Config.HTTP.Infos
	now := time.Now()
	for _, info := range infos {
		iyear, err := strconv.Atoi(info.BaseYear)
		if err != nil {
			logrus.Error(err)
			continue
		}
		for year := iyear; year <= now.Year(); year++ {
			url, err := url.Parse(config.Config.HTTP.Base)
			if err != nil {
				logrus.Error(err)
				continue
			}
			url.Path = path.Join(url.Path, info.Path, fmt.Sprintf(info.Label, year))
			logrus.Info(url.String())

			fpath, err := csv.GetPath(&info, year)
			if year < time.Now().Year() {
				if csv.FileExists(fpath) {
					logrus.Info("非今年的檔案，已下載過")
					continue
				}
			}

			grid, err := parseHTML(url.String())
			if err != nil {
				logrus.Error(err)
				continue
			}

			if err := csv.WriteToCSV(grid, &info, year); err != nil {
				logrus.Error(err)
				continue
			}
		}
	}
	return nil
}
