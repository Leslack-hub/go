package main

import "fmt"

var res [][]int
var num int
var target int

func combine(n int, k int) [][]int {
	num, target = n, k
	backtrack(1, []int{})
	return res
}

func backtrack(first int, temp []int) {
	if len(temp) == target {
		res = append(res, temp)
		return
	}
	for i := first; i <= num; i++ {
		temp = append(temp, i)
		record := make([]int, len(temp))
		copy(record, temp)
		backtrack(i+1, record)
		temp = temp[:len(temp)-1]
	}
}

func main() {
	fmt.Println(combine(1, 1))
}
