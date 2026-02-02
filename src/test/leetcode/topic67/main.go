package main

import (
	"fmt"
	"strconv"
)

func addBinary(a string, b string) string {
	i, j := len(a)-1, len(b)-1
	res := ""
	var flag, carray int
	for i >= 0 || j >= 0 {
		inta, intb := 0, 0
		if i >= 0 {
			inta = int(a[i] - '0')
		}
		if j >= 0 {
			intb = int(b[j] - '0')
		}
		carray = inta + intb + flag
		flag = 0
		if carray >= 2 {
			flag = 1
			carray = carray - 2
		}
		cur := strconv.Itoa(carray)
		res = cur + res
		i--
		j--
	}
	if flag == 1 {
		res = "1" + res
	}

	return res
}

func main() {
	fmt.Println(addBinary("10001", "1001"))
}
