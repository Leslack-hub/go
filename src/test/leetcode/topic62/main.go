package main

import (
	"fmt"
	"math"
)

func uniquePaths(m int, n int) int {
	path := [10][10]int{}
	for i := 0; i < m; i++ {
		path[i][0] = 1
	}
	for j := 0; j < n; j++ {
		path[0][j] = 1
	}
	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			path[i][j] = path[i-1][j] + path[i][j-1]
		}
	}
	fmt.Println(path)
	return path[m-1][n-1]
}
func fcc(n int) int {
	list := make([]int, n)
	list[0] = 1
	list[1] = 1
	for i := 2; i < n; i++ {
		list[i] = list[i-1] + list[i-2]
	}
	fmt.Println(list)
	return list[n-1]
}
func fcc2(n int) int {
	a, b := 0, 1
	for i := 0; i < n; i++ {
		a, b = b+a, a
	}
	return a
}
func gold(pattern []int, num int) int {
	array := make([]int, num+1)
	array[0] = 0
	for i := 1; i <= num; i++ {
		array[i] = math.MaxInt32
		for j := 0; j < len(pattern); j++ {
			if i >= pattern[j] && array[i-pattern[j]] != math.MaxInt32 {
				if array[i-pattern[j]] < array[i] {
					array[i] = array[i-pattern[j]] + 1
				}
			}
		}
	}
	if array[num] == math.MaxInt32 {
		return -1
	}
	return array[num]
}
func main() {
	//fmt.Println(uniquePaths(1, 1))
	//fmt.Println(fcc2(5))
	fmt.Println(gold([]int{2, 5, 7}, 27))
}
