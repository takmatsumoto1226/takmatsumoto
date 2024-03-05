package ftn

/*
 特徵一 : 號碼依大小排序後. 固定選出的號碼, 然後將找出未選出的號碼的位置
 例如 : 選出號碼 為 14 16 38 則. 剩下的2個號碼會有4個位置可以放
       , a1<14, 14<a2<16, 16<a3<38, a4>38
	   依照此方法, 找出特徵
*/
/*
	一組2連號, 3連號
*/
func (fa FTNArray) Continue2(p PickParam) FTNArray {
	result := FTNArray{}
	l := uint(0)

	for _, v := range fa {
		if l > p.Interval {
			break
		}

		if v.IsContinue2() {
			result = append(result, v)
			l++
		}
	}
	return result
}

func (fa FTNArray) Continue3(p PickParam) FTNArray {
	result := FTNArray{}
	l := uint(0)

	for _, v := range fa {
		if l > p.Interval {
			break
		}

		if v.IsContinue3() {
			result = append(result, v)
			l++
		}
	}
	return result
}

func (fa FTNArray) Continue4(p PickParam) FTNArray {
	result := FTNArray{}
	l := uint(0)

	for _, v := range fa {
		if l > p.Interval {
			break
		}

		if v.IsContinue4() {
			result = append(result, v)
			l++
		}
	}
	return result
}

func (fa FTNArray) Continue5(p PickParam) FTNArray {
	result := FTNArray{}
	l := uint(0)

	for _, v := range fa {
		if l > p.Interval {
			break
		}

		if v.IsContinue5() {
			result = append(result, v)
			l++
		}
	}
	return result
}

/*
二組2連號
*/
func (fa FTNArray) Continue22(p PickParam) FTNArray {
	result := FTNArray{}
	l := uint(0)

	for _, v := range fa {
		if l > p.Interval {
			break
		}

		if v.IsContinue22() {
			result = append(result, v)
			l++
		}
	}
	return result
}

/*
	2個號碼的組合號出現次數
*/

/*
	3個號碼的組合號出現次數
*/
/*
	4個號碼的組合號出現次數
*/

/**!SECTION
model
*/

func (fa FTNArray) DTree(p PickParam) FTNArray {
	result := FTNArray{}
	l := uint(0)

	for i, v := range fa {
		if l > p.Interval || i == len(fa)-1 {
			break
		}

		if fa[i+1].IsDTree(&v) {
			result = append(result, fa[i+1])
			result = append(result, v)
			result = append(result, *Empty())
			l++
		}
	}
	return result
}

func (fa FTNArray) UTree(p PickParam) FTNArray {
	result := FTNArray{}
	l := uint(0)

	for i, v := range fa {
		if l > p.Interval || i == len(fa)-1 {
			break
		}

		if v.IsUTree(&fa[i+1]) {
			result = append(result, fa[i+1])
			result = append(result, v)
			result = append(result, *Empty())
			l++
		}
	}
	return result
}
