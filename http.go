package main

import (
	"log"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

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
