package main

import (
	"fmt"
)

func LongestIncreasingSubsequence(A []int) int {
	n := len(A)
	f := make([]int, n)
	var max int
	for i := 0; i < n; i++ {
		f[i] = 1
		for j := 0; j < i; j++ {
			if A[j] < A[i] {
				if f[i] < f[j]+1 {
					f[i] = f[j] + 1
				}
			}
		}
		if max < f[i] {
			max = f[i]
		}
	}
	return max
}

func main() {
	fmt.Println(LongestIncreasingSubsequence([]int{4, 2, 4, 5, 3, 7}))
}
