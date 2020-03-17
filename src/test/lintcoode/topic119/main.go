package main

import (
	"fmt"
	"leslack/src/helper"
)

func EditDistance(A string, B string) int {
	m, n := len(A), len(B)
	f := make([][]int, m+1)
	for i := 0; i <= m; i++ {
		f[i] = make([]int, n+1)
	}

	for i := 0; i <= m; i++ {
		for j := 0; j <= n; j++ {
			if i == 0 {
				f[i][j] = j
				continue
			}
			if j == 0 {
				f[i][j] = i
				continue
			}
			f[i][j] = helper.Min(f[i][j-1], f[i-1][j], f[i-1][j-1]) + 1
			if A[i-1] == B[j-1] {
				f[i][j] = helper.Min(f[i][j], f[i-1][j-1])
			}
		}
	}

	fmt.Println(f)
	return f[m][n]
}

func main() {
	fmt.Println(EditDistance("abcd", "abcge"))
}
