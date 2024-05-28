package common

import (
	"fmt"
	"testing"
)

func Test_byteTest(t *testing.T) {
	var b byte = 0x0001
	var c byte = 0x0010
	fmt.Printf("%d\n", b)
	fmt.Printf("%d\n", c)
	d := c & b
	fmt.Printf("%d\n", d)
}
