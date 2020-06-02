package main

import (
	"fmt"
	"leslack/src/helper"
)

func minDistance(word1 string, word2 string) int {
	m, n := len(word1), len(word2)
	f := make([][]int, m+1)
	for i := 0; i <= m; i++ {
		f[i] = make([]int, n+1)
	}
	f[0][0] = 0
	for i := 1; i <= m; i++ {
		f[i][0] = i
	}
	for i := 1; i <= n; i++ {
		f[0][i] = i
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if word1[i-1] == word2[j-1] {
				f[i][j] = f[i-1][j-1]
			} else {
				f[i][j] = f[i-1][j-1] + 1
			}
			f[i][j] = helper.Min(f[i][j], f[i][j-1]+1, f[i-1][j]+1)
		}
	}
	return f[m][n]
}

func main() {
	fmt.Println(minDistance("horse", "ros"))
}
