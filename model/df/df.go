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

// 特徵種類
const (
	ContinuePickupNumber1 = iota // 跟前一期號碼相比, 有1個號碼連續出現
	ContinuePickupNumber2        // 跟前一期號碼相比, 有2個號碼連續出現
	ContinuePickupNumber3        // 跟前一期號碼相比, 有3個號碼連續出現
	ContinuePickupNumber4        // 跟前一期號碼相比, 有4個號碼連續出現
	ContinueNumber2              // 同一期出現相連號碼(2個) ex: 01 05 06 23 33
	ContinueNumber3              // 同一期出現相連號碼(3個) ex: 01 05 06 07 33
	ContinueNumber4              // 同一期出現相連號碼(3個) ex: 01 05 06 07 08
	ContinueNumber5              // 同一期出現相連號碼(3個) ex: 04 05 06 07 08
)

const (
	Next = iota
	Before
	Both
	None
)
