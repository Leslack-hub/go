package main

import (
	"fmt"
	"leslack/src/helper"
)

func LongestPalindromicSubsequence(s string) int {
	n := len(s)
	f := make([][]int, n)
	for i := 0; i < n; i++ {
		f[i] = make([]int, n)
	}

	for i := 0; i < n; i++ {
		f[i][i] = 1
	}

	for i := 0; i < n-1; i++ {
		if s[i] == s[i+1] {
			f[i][i+1] = 2
		} else {
			f[i][i+1] = 1
		}
	}
	for len := 3; len <= n; len++ {
		for i := 0; i <= n-len; i++ {
			j := i + len - 1
			f[i][j] = helper.Max(f[i][j-1], f[i+1][j])
			if s[i] == s[j] {
				f[i][j] = helper.Max(f[i][j], f[i+1][j-1]+2)
			}
		}
	}

	return f[0][n-1]
}

func main() {
	fmt.Println(LongestPalindromicSubsequence("bbbab"))
}
