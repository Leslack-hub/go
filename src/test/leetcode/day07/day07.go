package main

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
)

func reverse(x int) int {
	// 处理 res 的溢出问题
	var res int

	sign := 1
	if x < 0 {
		sign = -1
		x = x * -1
	}

	str := strconv.Itoa(x)

	var resv []byte
	for i := len(str); i > 0; i-- {
		resv = append(resv, str[i-1])
	}
	nums := bytes.Buffer{}
	for _, j := range resv {
		nums.WriteByte(j)
	}

	res, _ = strconv.Atoi(nums.String())
	res = res * sign
	if res > math.MaxInt32 || res < math.MinInt32 {
		res = 0
	}

	return res
}

func main() {
	fmt.Println(reverse(-321))
}
