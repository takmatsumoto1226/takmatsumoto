package common

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/goark/mt/v2"
	"github.com/goark/mt/v2/mt19937"
)

var defaultRand *rand.Rand
var rand19937 *mt.PRNG
var randtype = 0

func TimeTaken(t time.Time, name string) {
	elapsed := time.Since(t)
	log.Printf("TIME: %s took %s\n", name, elapsed)
}

func Save(content, filename string, index int) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err := f.WriteString(content)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println(fmt.Sprintf("%03d : ", index), l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func SaveJSON(content interface{}, filename string) {
	jsonString, _ := json.Marshal(content)
	os.WriteFile(filename, jsonString, os.ModePerm)
}

func SetRandomGenerator(t int) {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	randtype = t
	switch t {
	case 1:
		rand19937 = mt.New(mt19937.New(int64(binary.LittleEndian.Uint64(b[:]))))
	default:
		defaultRand = rand.New(rand.NewSource(int64(binary.LittleEndian.Uint64(b[:]))))
	}
}

func RandomNuber() uint64 {
	if rand19937 == nil && defaultRand == nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	switch randtype {
	case 1:
		return rand19937.Uint64()
	default:
		return defaultRand.Uint64()
	}
}

func MAX(x, y int) int {
	if x <= y {
		return y
	}
	return x
}

func MIN(x, y int) int {
	if x >= y {
		return y
	}
	return x
}

type LIMap map[int]bool

func (lm LIMap) Presentation() string {
	msg := ""
	for k, _ := range lm {
		msg = msg + fmt.Sprintf("%02d|", k)
	}
	return msg
}

func Test_montecarlo(t *testing.T) {
	samplesExponent := 10

	var r1 float64
	var r2 float64
	var heads float64
	samples := math.Pow(10, float64(samplesExponent))
	heads = 0
	for range make([]struct{}, int(samples)) {
		r1 = rand.Float64()
		r2 = rand.Float64()
		toss := math.Pow(r1-0.5, 2) + math.Pow(r2-0.5, 2)
		if toss < 0.25 {
			heads++
		}
	}

	area := samples * 0.25

	pi := heads / area

	fmt.Printf("pi estimation - %f\n", pi)
}

func MultiPI(samples int, threads int) float64 {
	threadSamples := samples / threads
	results := make(chan float64, threads)

	for j := 0; j < threads; j++ {
		go func() {
			var inside int
			for i := 0; i < threadSamples; i++ {
				x, y := rand.Float64(), rand.Float64()

				if x*x+y*y <= 1 {
					inside++
				}
			}
			results <- float64(inside) / float64(threadSamples) * 4
		}()
	}

	var total float64
	for i := 0; i < threads; i++ {
		total += <-results
	}

	return total / float64(threads)
}
