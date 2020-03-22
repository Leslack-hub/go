package main

import (
	"fmt"
	"leslack/src/helper"
)

/**
 * 双序列动态规划
 */
func LongestCommonSubsequence(A string, B string) int {
	m := len(A)
	n := len(B)
	f := make([][]int, m+1)
	for i := 0; i <= m; i++ {
		f[i] = make([]int, n+1)
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			f[i][j] = helper.Max(f[i-1][j], f[i][j-1])
			if A[i-1] == B[j-1] {
				f[i][j] = helper.Max(f[i][j], f[i-1][j-1]+1)
			}
		}
		fmt.Println(f)
	}
	return f[m][n]
}

func main() {
	fmt.Println(LongestCommonSubsequence("jiuzhang", "lijiang"))
}
