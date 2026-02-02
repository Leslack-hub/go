package main

import "fmt"

func grayCode(n int) []int {
	res := make([]int, 1<<n)
	index := 1
	for i := 0; i < n; i += 1 {
		r := 1<<i - 1
		for r >= 0 {
			v := res[r]
			v = v | (1 << i)
			res[index] = v
			r -= 1
			index += 1
		}
	}
	return res
}

func main() {
	fmt.Println(grayCode(22))
}
