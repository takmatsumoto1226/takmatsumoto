package df

const (
	InfoFTN = iota
	Info49
	InfoPOWER
	Info4STAR
)

const (
	HOT = iota
	COOL
)

const (
	REVERSE = iota
)

const (
	Descending = iota // raw data 年份大到小
	Ascending         // raw data 年份小到大
)

const (
	Biggerfront = iota // 球數出現次數統計後, 出現次數多得在前面
	Smallfront         // 球數出現次數統計後, 出現次數少的在前面
	Normal             // ball的數字由小到大排序
)
