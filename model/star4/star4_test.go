package star4

import (
	"fmt"
	"lottery/config"
	"testing"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func Test_listLikeExecl(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = Star4Manager{numberToIndex: map[string]int{}}
	as.Prepare()
	as.list()
}

func Test_start4Permutations(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = Star4Manager{numberToIndex: map[string]int{}}
	as.Prepare()
	arr := as.least().Permutation()
	for _, ar := range arr {
		fmt.Println(ar)
	}
}
func Test_permutations(t *testing.T) {
	arr := permutations([]int{5, 5, 3, 8})
	fmt.Println(len(arr))
	for _, ar := range arr {
		fmt.Println(ar)
	}
}
