package common

import (
	"encoding/csv"
	"os"
)

// ReadCSV ...
func ReadCSV(fullpath string) ([][]string, error) {

	f, err := os.Open(fullpath)
	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
}
