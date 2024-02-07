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
	as.List.Presentation()
}

func Test_start4Permutations(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = Star4Manager{numberToIndex: map[string]int{}}
	as.Prepare()
	arr := as.List.Least().Permutation()
	for _, ar := range arr {
		fmt.Println(ar)
	}
}
func Test_permutations(t *testing.T) {
	arr := permutations([]int{1, 2, 2, 6})
	fmt.Println(len(arr))
	for _, ar := range arr {
		fmt.Println(ar)
	}
}

func Test_findnumber(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = Star4Manager{numberToIndex: map[string]int{}}
	as.Prepare()
	arr := permutations([]int{9, 5, 2, 7})
	for _, v := range arr {
		str := fmt.Sprintf("%d%d%d%d", v[0], v[1], v[2], v[3])
		fmt.Println(str)
		fmt.Println(as.findNumbers(str))
	}
}
