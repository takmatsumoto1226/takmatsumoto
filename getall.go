package main

import (
	"fmt"
	"io"
	"lottery/config"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
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

			fpath, _ := csv.GetPath(&info, year)
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

func getAllFromURL() error {
	for _, dp := range config.Config.AokURLs {
		fullURLFile := dp.URL

		// Build fileName from fullPath
		fileURL, err := url.Parse(fullURLFile)
		if err != nil {
			return err
		}
		path := fileURL.Path
		segments := strings.Split(path, "/")
		fileName := segments[len(segments)-1]

		if err := os.MkdirAll(dp.LocaPath, 0755); err != nil {
			return err
		}

		fullPath := filepath.Join(dp.LocaPath, fileName)

		// Create blank file
		file, err := os.Create(fullPath)
		if err != nil {
			return err
		}
		client := http.Client{
			CheckRedirect: func(r *http.Request, via []*http.Request) error {
				r.URL.Opaque = r.URL.Path
				return nil
			},
		}
		// Put content on file
		resp, err := client.Get(fullURLFile)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		size, err := io.Copy(file, resp.Body)

		defer file.Close()

		fmt.Printf("Downloaded a file %s with size %d\n", fileName, size)
		time.Sleep(time.Second * 1)
	}
	return nil
}
