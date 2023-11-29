package googleplayparser

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func Test_parser(t *testing.T) {
	var vs AppVersion
	decbs, err := os.ReadFile("./test.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	err = json.Unmarshal(decbs, &vs)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(vs)

}
