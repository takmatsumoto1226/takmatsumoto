package bl

import (
	"lottery/config"
	"testing"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func Test_listLikeExecl(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = BigLotterysManager{numberToIndex: map[string]int{}}
	as.Prepare()
	as.RevList.Presentation()
}

func Test_findnumbers(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = BigLotterysManager{numberToIndex: map[string]int{}}
	as.Prepare()
	as.RevList.findNumbers([]string{"03", "13", "22"}, true).Presentation()
}
