package common

import (
	"fmt"
	"io"
	"log"
	"lottery/config"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"lottery/csv"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

func GetAll() error {
	infos := config.Config.HTTP.Infos
	now := time.Now()
	for _, info := range infos {
		iyear, err := strconv.Atoi(info.BaseYear)
		if err != nil {
			// logrus.Error(err)
			fmt.Errorf("%v", err)
			continue
		}
		for year := iyear; year <= now.Year(); year++ {
			url, err := url.Parse(config.Config.HTTP.Base)
			if err != nil {
				fmt.Errorf("%v", err)
				continue
			}
			url.Path = path.Join(url.Path, info.Path, fmt.Sprintf(info.Label, year))
			// logrus.Info(url.String())
			fmt.Println(url.String())

			fpath, _ := csv.GetPath(&info, year)
			if year < time.Now().Year() {
				if csv.FileExists(fpath) {
					// logrus.Info("非今年的檔案，已下載過")
					fmt.Println("非今年的檔案，已下載過")
					continue
				}
			}

			grid, err := parseHTML(url.String())
			if err != nil {
				// logrus.Error(err)
				fmt.Errorf("%v", err)
				continue
			}

			if err := csv.WriteToCSV(grid, &info, year); err != nil {
				// logrus.Error(err)
				fmt.Errorf("%v", err)
				continue
			}
		}
	}
	return nil
}

func GetAllFromURL() error {
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

// res, err := http.Get("http://www.nfd.com.tw/lottery/49-year/49-2004.htm")
// res, err := http.Get("http://www.nfd.com.tw/lottery/power-38/2008.htm")
// res, err := http.Get("http://www.nfd.com.tw/lottery/39-year/39-2007.htm")
// res, err := http.Get("http://www.nfd.com.tw/lottery/4-star/2003.htm")

func parseHTML(URL string) ([][]string, error) {
	// Request the HTML page.
	res, err := http.Get(URL)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		return nil, err
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	grid := [][]string{}
	// Find the review items
	doc.Find("table").Each(func(i int, s *goquery.Selection) {
		s.Find("tr").Each(func(j int, d *goquery.Selection) {
			if j > 0 {
				row := []string{}
				d.Find("td").Each(func(k int, g *goquery.Selection) {
					title := g.Find("b").Text()
					// fmt.Println(title)
					reg, err := regexp.Compile("[^0-9]+")
					if err != nil {
						log.Fatal(err)
					}
					processedString := reg.ReplaceAllString(title, "")
					row = append(row, processedString)
				})
				logrus.Info(row)
				grid = append(grid, row)
			}
		})
	})
	// logrus.Info(grid)
	return grid, nil
}
