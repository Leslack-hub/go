package service

func GetSum(inputs ...int64) int64 {
	var res int64
	for _, v := range inputs {
		res += v
	}
	return res
}
