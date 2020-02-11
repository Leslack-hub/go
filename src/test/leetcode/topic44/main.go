package main

import "fmt"

func isMatch(s string, p string) bool {
	m, n := len(s), len(p)
	dp := make([][]bool, m+1)
	// 原字符串
	for i := 0; i <= m; i++ {
		dp[i] = make([]bool, n+1)
	}
	dp[0][0] = true
	// 匹配的值
	for i := 1; i < n+1; i++ {
		dp[0][i] = dp[0][i-1] && p[i-1] == '*'
	}
	fmt.Println(dp)
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if s[i-1] == p[j-1] || p[j-1] == '?' {
				dp[i][j] = dp[i-1][j-1]
			} else if p[j-1] == '*' {
				dp[i][j] = dp[i][j-1] || dp[i-1][j]
			}
			fmt.Printf("dp[%d][%d] = %v\n", i, j, dp[i][j])
		}
	}
	return dp[m][n]
}

func main() {
	fmt.Println(isMatch("acdcb",
		"*cbb"))
}
