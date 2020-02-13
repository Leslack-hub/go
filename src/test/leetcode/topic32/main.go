package main

import "fmt"

func longestValidParentheses(s string) int {
	var left int
	record := make([]int, len(s))
	// 如果left == 1 表示上一个为(
	for k, v := range s {
		if v == '(' {
			left++
		} else if left > 0 {
			left--
			record[k] = 2
		}
	}
	for i := 0; i < len(s); i++ {
		if record[i] == 2 {
			j := i - 1
			for record[j] != 0 {
				j--
			}
			record[i], record[j] = 1, 1
		}
	}

	var temp, max int
	// 统计1的连续次数
	for _, v := range record {
		if v == 0 {
			temp = 0
			continue
		}
		temp++
		if temp > max {
			max = temp
		}
	}

	return temp
}
func main() {
	fmt.Println(longestValidParentheses(")()())"))
}
