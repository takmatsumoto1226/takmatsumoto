package star4

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
)

type RawStar4 struct {
	Date1      string `json:"date1"`
	YearIndex1 string `json:"yearindex1"`
	Balls1     string `json:"balls1"`
	Space      string `json:"space"`
	Date2      string `json:"date2"`
	YearIndex2 string `json:"yearindex2"`
	Balls2     string `json:"balls2"`
}

const (
	arrIdxDate1 = iota
	arrIdxYearIndex1
	arrIdxBalls1
	arrIdxSpace
	arrIdxDate2
	arrIdxYearIndex2
	arrIdxBalls2
	RawStar4Total
)

type Star4 struct {
	Year      string `json:"year"`
	Date      string `json:"date"`
	YearIndex string `json:"yearindex"`
	Balls     string `json:"numbers"`
	P1        string `json:"p1"`
	P2        string `json:"p2"`
	P3        string `json:"p3"`
	P4        string `json:"p4"`
}

func (fa *Star4) toStringArray() []string {
	return []string{fa.P1, fa.P2, fa.P3, fa.P4}
}

func (fa *Star4) Normalize() error {
	a := []rune(fa.Balls)
	if len(a) == 4 {
		fa.P1 = string(a[0])
		fa.P2 = string(a[1])
		fa.P3 = string(a[2])
		fa.P4 = string(a[3])
		return nil
	}
	return errors.New("數字不符合4星彩格式")
}

func (fa Star4) formRow() {
	rowmsg := fmt.Sprintf("%s|", fa.Year)
	rowmsg = rowmsg + fmt.Sprintf("%s|%s|%s|%s|%s|%s|", fa.Date, fa.Balls, fa.P1, fa.P2, fa.P3, fa.P4)
	fmt.Println(rowmsg)
}

func NewStar4(year string, arr []string) []Star4 {
	if len(arr) == RawStar4Total {
		result := []Star4{}
		if len(arr[arrIdxDate1]) > 0 && len(arr[arrIdxBalls1]) > 0 {
			result = append(result, Star4{
				Year:      year,
				Date:      arr[arrIdxDate1],
				YearIndex: arr[arrIdxYearIndex1],
				Balls:     arr[arrIdxBalls1],
				P1:        "",
				P2:        "",
				P3:        "",
				P4:        "",
			})
		}

		if len(arr[arrIdxDate2]) > 0 && len(arr[arrIdxBalls2]) > 0 {
			result = append(result, Star4{
				Year:      year,
				Date:      arr[arrIdxDate2],
				YearIndex: arr[arrIdxYearIndex2],
				Balls:     arr[arrIdxBalls2],
				P1:        "",
				P2:        "",
				P3:        "",
				P4:        "",
			})
		}
		return result
	}

	logrus.Error("NewStar4 資料格式錯誤")
	return nil
}

func Empty() *Star4 {
	return &Star4{"====", "====", "==", "==", "==", "==", "==", "=="}
}

func (fa *Star4) toInts() []int {
	p1, _ := strconv.Atoi(fa.P1)
	p2, _ := strconv.Atoi(fa.P2)
	p3, _ := strconv.Atoi(fa.P3)
	p4, _ := strconv.Atoi(fa.P4)
	return []int{p1, p2, p3, p4}
}

func (fa Star4) Permutation() [][]int {
	var helper func([]int, int)
	res := [][]int{}

	helper = func(arr []int, n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			// if _, ok := trims[tmp]; !ok {
			res = append(res, tmp)
			// }
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(fa.toInts(), 4)
	return res
}
