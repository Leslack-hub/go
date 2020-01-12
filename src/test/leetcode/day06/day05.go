package main

import (
	"bytes"
	"fmt"
)
func convert(s string, numRows int) string {
	if len(s) < 2 || numRows == 1 {
		return s
	}
	list := make([][]byte, numRows)
	i := 0
	re := false
	for _, v := range s {
		list[i] = append(list[i], byte(v))
		if re {
			i--
			if i == 0 {
				re = false
			}
		} else {
			i++
			if i == numRows-1 {
				re = true
			}
		}
	}
	res := bytes.Buffer{}
	for _, i := range list {
		for _, j := range i {
			res.WriteByte(j)
		}
	}
	return res.String()
}

func main() {
	// 输出 LEETCODEISHIRING
	// 输出 LCIRETOESIIGEDHN
	fmt.Println(convert("AB", 1))
}
