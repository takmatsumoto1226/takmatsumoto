package common

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
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
