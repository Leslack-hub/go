package main

import (
	"fmt"
	"strconv"
)

func countAndSay(n int) string {
	res := ""
	if n == 1 {
		return "1"
	}
	res += countAndSay(n - 1)
	var str string
	for i := 0; i < len(res); i++ {
		q := i + 1
		if q < len(res) {
			p := 1
			val := res[i]
			for q < len(res) {
				if res[q] == val {
					p++
					q++
					i++
				} else {
					q++
					break
				}
			}
			str += strconv.Itoa(p) + string(val)
		} else {
			str += "1" + string(res[i])
		}
	}
	return str
}

func main() {
	fmt.Println(countAndSay(5))
}
