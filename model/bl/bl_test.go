package bl

import (
	"encoding/binary"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"lottery/algorithm"
	"lottery/config"
	"lottery/model/common"
	"lottery/model/df"
	"lottery/model/interf"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/stat/combin"

	crypto_rand "crypto/rand"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func Test_listLikeExecl(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = BigLotterysManager{numberToIndex: map[string]int{}}
	as.Prepare()
	fmt.Println(as.RevList.Presentation())
}

func Test_findnumbers(t *testing.T) {
	config.LoadConfig("../../config.yaml")
	var as = BigLotterysManager{numberToIndex: map[string]int{}}
	as.Prepare()
	fmt.Println(as.RevList.findNumbers([]string{"03", "13", "22"}, df.Both).Presentation())
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

func Test_random2(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Top Price Taken Time")
	// Create and seed the generator.
	// Typically a non-fixed seed should be used, such as time.Now().UnixNano().
	// Using a fixed seed will produce the same output on every run.
	config.LoadConfig("../../config.yaml") // 17591400
	var as = BigLotterysManager{numberToIndex: map[string]int{}}
	as.Prepare()

	balls := 6
	combarr := combin.Combinations(49, balls)
	// lens := len(combarr)
	th := interf.Threshold{Round: 1, Value: 11, SampleTime: 4, Sample: len(combarr)}

	for i := 0; i < th.Round; i++ {
		result := map[string]int{}

		var b [8]byte
		_, err := crypto_rand.Read(b[:])
		if err != nil {
			panic("cannot seed math/rand package with cryptographically secure random number generator")
		}

		rnumber := rand.New(rand.NewSource(int64(binary.LittleEndian.Uint64(b[:]))))

		for _, v := range combarr {
			balls := NewPowerWithInts(v)
			result[balls.Key()] = 0
		}

		total := int(float64(len(combarr)) * th.SampleTime)
		fmt.Println(total)

		for i := 0; i < total; i++ {
			index := uint32(rnumber.Uint32() % uint32(len(result)))
			balls := NewPowerWithInts(combarr[index])
			if v, ok := result[balls.Key()]; ok {
				result[balls.Key()] = v + 1
			}

			// fmt.Println(combarr[index])
		}

		lottos := as.List.WithRange(20)
		featuresList := BLList{}
		count := 0
		for k, v := range result {
			if v > th.Value {
				fmt.Printf("%v:%v\n", k, v)
				arr := strings.Split(k, "_")
				fmt.Println(as.List.findNumbers(arr, df.Next).Presentation())

				bl := NewPowerWithString(arr)
				for _, l := range lottos {
					if l.CompareFeature(bl) {
						featuresList = append(featuresList, *bl)
					}
				}
				count++
			}
		}
		fmt.Printf("%d , %d \n", count*50, count)
		fmt.Println("")
		fmt.Println("features")
		fmt.Println(featuresList.Presentation())
		fmt.Printf("done Round : %02d\n", th.Round)
	}

}

func Test_ExportAllNumber(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_TenGroupManager")
	config.LoadConfig("../../config.yaml")
	var mgr = BigLotterysManager{numberToIndex: map[string]int{}}
	mgr.Prepare()
	mgr.List.Reverse().CSVExport1()

}

func Test_ExportAllCombination(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_ExportAllCombination")
	combarr := combin.Combinations(49, 6)
	// 建立 CSV 檔案
	file, err := os.Create("data.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 建立 CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 將資料寫入 CSV
	for _, f := range combarr {
		if err := writer.Write([]string{
			fmt.Sprintf("%02d", f[0]+1),
			fmt.Sprintf("%02d", f[1]+1),
			fmt.Sprintf("%02d", f[2]+1),
			fmt.Sprintf("%02d", f[3]+1),
			fmt.Sprintf("%02d", f[4]+1),
			fmt.Sprintf("%02d", f[5]+1),
		}); err != nil {
			panic(err)
		}
	}

	// 檢查是否有錯誤
	if err := writer.Error(); err != nil {
		panic(err)
	}

	println("CSV 檔案輸出完成！")

}

func Test_ExportbinaryAllNumber(t *testing.T) {
	defer common.TimeTaken(time.Now(), "Test_TenGroupManager")
	config.LoadConfig("../../config.yaml")
	var mgr = BigLotterysManager{numberToIndex: map[string]int{}}
	mgr.Prepare()
	mgr.List.Reverse().CSVExport("/Users/tak 1/Documents/gitlab_project/LotteryAi/resultbl.csv")

}

func Test_NewWithStrings(t *testing.T) {
	fileName := "/Users/tak 1/Desktop/examplebl.csv"
	config.LoadConfig("../../config.yaml")
	var mgr = BigLotterysManager{numberToIndex: map[string]int{}}
	mgr.Prepare()

	// 打開 CSV 檔案
	file, err := os.Open(fileName)
	if err != nil {

		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// 讀取 CSV 檔案內容
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Error reading CSV file: %v\n", err)
		return
	}

	// 檢查是否有資料
	if len(records) < 1 {
		fmt.Println("No data in CSV file")
		return
	}
	fmt.Println(len(records))

	arr := BLList{}
	for _, record := range records {
		ftn := NewPowerWithString(record)
		arr = append(arr, *ftn)
	}
	fmt.Println(arr.Presentation())
	newtop := mgr.List[0]
	fmt.Printf("Cost : %d\n", len(arr)*50)
	fmt.Printf("Top:\n%s\n", newtop.simpleFormRow())
	fmt.Println(arr.AdariPrice(&newtop))
}
