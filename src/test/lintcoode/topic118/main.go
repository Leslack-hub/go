package main

import "fmt"

/**
 * 两个字符相同最长长度
 */
func DistinctSubsequences(A string, B string) int {
	m, n := len(A), len(B)
	f := make([][]int, m+1)
	for i := 0; i <= m; i++ {
		f[i] = make([]int, n+1)
	}

	for i := 0; i <= m; i++ {
		for j := 0; j <= n; j++ {
			if j == 0 {
				f[i][j] = 1
				continue
			}
			if i == 0 {
				f[i][j] = 0
				continue
			}
			f[i][j] = f[i-1][j]
			if A[i-1] == B[j-1] {
				f[i][j] += f[i-1][j-1]
			}
		}
	}
	fmt.Println(f)

	return f[m][n]
}

func main() {
	fmt.Println(DistinctSubsequences("rabbbit", "rabbit"))
}
