package main

import (
	"fmt"
)

func multiply(num1 string, num2 string) string {
	if num1 == "0" || num2 == "0" {
		return "0"
	}
	byteList1 := []byte(num1)
	byteList2 := []byte(num2)
	temp := make([]int, len(byteList1)+len(byteList2))
	for i := 0; i < len(byteList1); i++ {
		for j := 0; j < len(byteList2); j++ {
			temp[i+j+1] += int(byteList1[i]-'0') * int(byteList2[j]-'0')
		}
	}
	for i := len(temp) - 1; i > 0; i-- {
		temp[i-1] += temp[i] / 10
		temp[i] = temp[i] % 10
	}
	if temp[0] == 0 {
		temp = temp[1:]
	}
	res := make([]byte, len(temp))
	for i := 0; i < len(temp); i++ {
		res[i] = '0' + byte(temp[i])
	}

	return string(res)
}

func numToString(num1 string) int {
	multi, num := 1, 0
	for i := len(num1) - 1; i >= 0; i-- {
		num += multi * int(num1[i]-'0')
		multi *= 10
	}
	return num
}
func main() {
	fmt.Println(multiply("123",
		"3"))
}
