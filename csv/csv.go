package csv

import (
	"encoding/csv"
	"fmt"
	"lottery/config"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

// WriteToCSV ...
func WriteToCSV(grid [][]string, info *config.Info, year int) error {

	fpath, err := GetPath(info, year)
	logrus.Info(fpath)

	if grid == nil {
		logrus.Error("沒有資料")
		return err
	}
	if err != nil {
		logrus.Error(err)
		return err
	}
	if year < time.Now().Year() {
		if FileExists(fpath) {
			logrus.Info("非今年的檔案，已下載過")
			return nil
		}
	}

	file, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.UseCRLF = true
	defer writer.Flush()

	for _, row := range grid {
		if err := writer.Write(row); err != nil {
			return err // let's return errors if necessary, rather than having a one-size-fits-all error handler
		}
	}

	return nil
}

// FileExists ...
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// GetPath 取得csv檔案的路徑
func GetPath(info *config.Info, year int) (string, error) {
	ext := filepath.Ext(info.Label)
	filename := fmt.Sprintf(info.Label, year)
	filename = filename[:len(filename)-len(ext)] + ".csv"
	dirpath := filepath.Join(config.Config.Local.Path, info.Path)
	err := os.MkdirAll(dirpath, 0755)
	if err != nil {
		return "", nil
	}
	fullpath := filepath.Join(dirpath, filename)
	return fullpath, nil
}
