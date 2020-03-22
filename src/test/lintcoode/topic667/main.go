package main

import "leslack/src/helper"

func LongestPalindromicSubsequence(s string) int {
	n := len(s)
	if n <= 1 {
		return 1
	}

	f := make([][]int, n)
	for i := 0; i < n; i++ {
		f[i] = make([]int, n)
		f[i][i] = 1
	}

	for i := 0; i < n-1; i++ {
		if s[i] == s[i+1] {
			f[i][i+1] = 2
		} else {
			f[i][i+1] = 1
		}
	}

	for length := 3; length <= n; length++ {
		for i := 0; i <= n-i; i++ {
			j := i + length - 1
			f[i][j] = helper.Max(f[i+1][j], f[i][j-1])
			if s[i] == s[j] {
				f[i][j] = helper.Max(f[i][j], f[i+1][j-1]+2)
			}
		}
	}

}

func main() {
	LongestPalindromicSubsequence("bbbab")
}
