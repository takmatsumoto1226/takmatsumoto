package googleplayparser

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (r *AppVersion) history(packagenames []string) error {

	logrus.Info("Init history data")
	path, err := filepath.Abs("./comments")
	if err != nil {
		return errors.Wrap(err, "Google Play Comment Load History error")
	}
	for _, pkn := range packagenames {
		fullpath := filepath.Join(path, pkn+".json")
		var vs AppVersion
		decbs, err := ioutil.ReadFile(fullpath)
		if err != nil {
			logrus.Error(err.Error())
		}
		err = json.Unmarshal(decbs, &vs)
		if err != nil {
			logrus.Error(err.Error())
		}
	}

	return nil
}

type AndroidVersion struct {
	ProductVersionID int    `json:"product_version_id"`
	UpdateDate       int    `json:"update_date"`
	Name             string `json:"name"`
	ReleaseNotes     string `json:"release_notes"`
	Date             int    `json:"date"`
	Facet            string `json:"facet"`
}

type AndroidVersions struct {
	Versions []AndroidVersion `json:"facets"`
}

type AppVersion struct {
	Data AndroidVersions `json:"data"`
	PID  string          `json:"pid"`
}
