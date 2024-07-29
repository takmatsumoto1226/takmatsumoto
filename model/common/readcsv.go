package common

import (
	"encoding/csv"
	"log"
	"os"
)

// ReadCSV ...
func ReadCSV(fullpath string) ([][]string, error) {

	f, err := os.Open(fullpath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
}
