package config

import (
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Info ...
type Info struct {
	Label    string `yaml:"label"`
	Path     string `yaml:"path"`
	BaseYear string `yaml:"baseyear"`
}

// HTTP ...
type HTTP struct {
	Base  string `yaml:"base"`
	Infos []Info `yaml:"infos"`
}

// Local ...
type Local struct {
	Path string `yaml:"path"`
}

type Combination struct {
	Path     string `yaml:"path"`
	Template string `yaml:"template"`
}

// Param ...
type Param struct {
	HTTP             HTTP          `yaml:"http"`
	Local            Local         `yaml:"local"`
	AokURLs          []DocInfo     `yaml:"aokurls"`
	CombinationsInfo []Combination `yaml:"combinations"`
}

type DocInfo struct {
	URL      string `yaml:"url"`
	LocaPath string `yaml:"local_path"`
}

// Config ...
var Config = Param{}

// LoadConfig ...
func LoadConfig(filepath string) {
	bs, err := ioutil.ReadFile(filepath)
	if err != nil {
		epath, _ := os.Executable()
		logrus.Infof("pwd %s", epath)
		logrus.Errorf("read config %s, %v", filepath, err)
		os.Exit(1)
	}
	err = yaml.Unmarshal(bs, &Config)
	if err != nil {
		logrus.Errorf("umarshal config: %v", err)
		os.Exit(1)
	}
}
