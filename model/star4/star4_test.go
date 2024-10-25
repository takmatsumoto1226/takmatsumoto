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
	var as = NewStar4Manager()
	as.Prepare()
	as.List.Presentation()
}

func Test_start4Permutations(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = NewStar4Manager()
	as.Prepare()
	arr := as.List.Least().Permutation()
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

func Test_findnumber(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = NewStar4Manager()
	as.Prepare()
	// arr := permutations([]int{9, 5, 2, 7})
	arr := [][]int{{9, 1, 9, 1}}
	for _, v := range arr {
		str := fmt.Sprintf("%d%d%d%d", v[0], v[1], v[2], v[3])
		as.findNumbers(str).Presentation()
	}
}

func Test_quickSort(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = NewStar4Manager()
	as.Prepare()

	as.List.quickSort().Presentation()
}

func Test_statics(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = NewStar4Manager()
	as.Prepare()
	as.info.formRow()
	fmt.Println(as.info.Len())
}

func Test_bitsTest(t *testing.T) {
	// bits.Len8(00100100)
}
