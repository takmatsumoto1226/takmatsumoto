package pw

import (
	"fmt"
	"strconv"
	"strings"
)

// Balls ...
type Balls []Ball

func (bsi Balls) Presentation(dshift bool) string {
	rowmsg := "\n  "
	if dshift {
		rowmsg = "          "
	}

	for _, bi := range bsi {
		rowmsg = rowmsg + fmt.Sprintf("%3d|", bi.Digit)
	}
	rowmsg = rowmsg + "\n  "

	for _, bi := range bsi {
		rowmsg = rowmsg + fmt.Sprintf("%3d|", bi.Count)
	}
	return rowmsg
}

// Ball
type Ball struct {
	Number   string `json:"number"`
	Position int    `json:"position"`
	Digit    int    `json:"digit"`
	Period   int    `json:"period"`
	Continue int    `json:"continue"`
	Count    uint   `json:"count"`
}

func (b *Ball) Illegal() bool {
	return b.Number == "" || b.Number == "00"
}

func (b *Ball) Same(cb Ball) bool {
	return b.Digit == cb.Digit
}

func NewBallS(n string, pos int) *Ball {
	iB, _ := strconv.Atoi(n)
	if strings.Contains(n, "==") {
		return &Ball{Number: n, Position: 0, Digit: 0}
	} else {
		return &Ball{Number: n, Position: pos, Digit: iB}
	}
}

func NewBallI(n int, pos int) *Ball {
	number := fmt.Sprintf("%02d", n)
	if n == 0 {
		number = "=="
	}
	return &Ball{Number: number, Position: pos, Digit: n}
}
