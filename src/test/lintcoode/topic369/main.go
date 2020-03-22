package main

import (
	"fmt"
	"leslack/src/helper"
)

func CoinsInAlineII(A []int) bool {
	// f[i][j] 为i一方在另一方面对A[i..j]这些数字时的最大差值
	n := len(A)
	f := make([][]int, n)
	for i := 0; i < n; i++ {
		f[i] = make([]int, n)
	}

	for i := 0; i < n; i++ {
		f[i][i] = A[i]
	}

	for length := 2; length <= n; length++ {
		for i := 0; i <= n-length; i++ {
			j := i + length - 1
			f[i][j] = helper.Max(A[i]-f[i+1][j], A[j]-f[i][j-1])
		}
	}
	fmt.Println(f)
	return f[0][n-1] >= 0
}

func main() {
	fmt.Println(CoinsInAlineII([]int{1, 5, 233, 7}))
}
