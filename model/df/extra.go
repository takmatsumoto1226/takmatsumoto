package df

// 以下特徵皆假設 IBalls 已由小到大排序

var continueRowTypeNames = map[int]string{
	ContinueRowNone: "none",
	ContinueRow2:    "2",
	ContinueRow3:    "3",
	ContinueRow4:    "4",
	ContinueRow5:    "5",
	ContinueRow22:   "22",
	ContinueRow32:   "32",
}

// ContinueRowTypeName 連號型態的文字表示
func ContinueRowTypeName(t int) string {
	if n, ok := continueRowTypeNames[t]; ok {
		return n
	}
	return "unknown"
}

// AC 值 : 號碼兩兩差值的種類數 - (球數 - 1)
func (f *Feature) AC() int {
	diffs := map[int]bool{}
	for i := 0; i < len(f.IBalls); i++ {
		for j := i + 1; j < len(f.IBalls); j++ {
			diffs[f.IBalls[j]-f.IBalls[i]] = true
		}
	}
	return len(diffs) - (len(f.IBalls) - 1)
}

// Span 跨度 : 最大號 - 最小號
func (f *Feature) Span() int {
	if len(f.IBalls) == 0 {
		return 0
	}
	return f.IBalls[len(f.IBalls)-1] - f.IBalls[0]
}

// CountAtLeast 大於等於 n 的號碼個數 (大小比用)
func (f *Feature) CountAtLeast(n int) int {
	count := 0
	for _, b := range f.IBalls {
		if b >= n {
			count++
		}
	}
	return count
}

// Mod3Counts 012路 : 除以3餘0, 餘1, 餘2 的個數
func (f *Feature) Mod3Counts() [3]int {
	r := [3]int{}
	for _, b := range f.IBalls {
		r[b%3]++
	}
	return r
}

// TailSum 尾數和
func (f *Feature) TailSum() int {
	sum := 0
	for _, b := range f.IBalls {
		sum += b % 10
	}
	return sum
}

// SameTailGroups 同尾號組數 (同一個尾數出現2次以上算1組)
func (f *Feature) SameTailGroups() int {
	tails := map[int]int{}
	for _, b := range f.IBalls {
		tails[b%10]++
	}
	groups := 0
	for _, c := range tails {
		if c >= 2 {
			groups++
		}
	}
	return groups
}

// Gaps 相鄰號碼間隔
func (f *Feature) Gaps() []int {
	gaps := []int{}
	for i := 0; i < len(f.IBalls)-1; i++ {
		gaps = append(gaps, f.IBalls[i+1]-f.IBalls[i])
	}
	return gaps
}

// MaxGap 最大間隔
func (f *Feature) MaxGap() int {
	m := 0
	for _, g := range f.Gaps() {
		if g > m {
			m = g
		}
	}
	return m
}

// MinGap 最小間隔
func (f *Feature) MinGap() int {
	gaps := f.Gaps()
	if len(gaps) == 0 {
		return 0
	}
	m := gaps[0]
	for _, g := range gaps {
		if g < m {
			m = g
		}
	}
	return m
}
