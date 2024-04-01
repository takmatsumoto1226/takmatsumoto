package bl

import (
	"encoding/json"
	"fmt"
	"lottery/algorithm"
	"lottery/config"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/stat/combin"
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

func Test_combination(t *testing.T) {
	// fmt.Println(algorithm.All([]string{"09", "14", "30"}))
	// fmt.Println(Ball39())
	balls := 6
	combarr := algorithm.Combinations(Ball49(), balls)

	bytes, err := json.Marshal(combarr)
	if err != nil {
		logrus.Error(err)
		return
	}

	err = os.WriteFile(fmt.Sprintf("./blcombination%d.json", balls), bytes, 0777)
	if err != nil {
		logrus.Error(err)
		return
	}

}

func Test_combination2(t *testing.T) {
	// fmt.Println(algorithm.All([]string{"09", "14", "30"}))
	// fmt.Println(Ball39())
	balls := 6
	combarr := combin.Combinations(49, balls)
	fmt.Printf("combination : %d\n", len(combarr))

	// bytes, err := json.Marshal(combarr)
	// if err != nil {
	// 	logrus.Error(err)
	// 	return
	// }

	// err = os.WriteFile(fmt.Sprintf("./ftncombination%d.json", balls), bytes, 0777)
	// if err != nil {
	// 	logrus.Error(err)
	// 	return
	// }

	// bl : 13983816
	// ftn : 575757
	// pw : 2760681

}

func Test_random(t *testing.T) {
	// Create and seed the generator.
	// Typically a non-fixed seed should be used, such as time.Now().UnixNano().
	// Using a fixed seed will produce the same output on every run.

	rseed := rand.New(rand.NewSource(time.Now().Unix()))

	// Uint32 410958445 3256529108 1564474733

	fmt.Println("")
	// for i := 0; i < 24; i++ {
	fmt.Println(int(time.Now().Unix()))
	fmt.Printf("%v\n", rseed.Intn(int(time.Now().Unix())))
	// }
	rnumber := rand.New(rand.NewSource(int64(rseed.Intn(int(time.Now().Unix())))))

	balls := 6
	combarr := combin.Combinations(49, balls)

	for i := 0; i < 10; i++ {
		index := int(rnumber.Intn(len(combarr)))
		// fmt.Println(index)
		time.Sleep(time.Second)
		fmt.Println(combarr[index])
	}

}
