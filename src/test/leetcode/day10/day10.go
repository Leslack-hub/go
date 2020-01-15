package main

import (
	"fmt"
)

func isMatch(s, p string) bool {
	return match(s, p, len(s), len(p))
}

func match(s, p string, i, j int) bool {
	if i == 0 && j == 0 {
		return true
	}
	if i != 0 && j == 0 {
		return false
	}
	if i == 0 && j != 0 {
		if p[j-1] == '*' {
			return match(s, p, i, j-2)
		}
		return false
	}
	// aabbb a*b 5 3
	if s[i-1] == p[j-1] || p[j-1] == '.' {
		return match(s, p, i-1, j-1)
	} else if p[j-1] == '*' {
		if match(s, p, i, j-2) {
			return true
		}
		if s[i-1] == p[j-2] || p[j-2] == '.' {
			return match(s, p, i-1, j)
		}
		return false
	}
	return false
}

func isMatch2(s, p string) bool {
	sSize := len(s)
	pSize := len(p)

	dp := make([][]bool, sSize+1)
	for i := range dp {
		dp[i] = make([]bool, pSize+1)
	}

	/* dp[i][j] 代表了 s[:i] 能否与 p[:j] 匹配 */

	dp[0][0] = true
	/**
	 * 根据题目的设定， "" 可以与 "a*b*c*" 相匹配
	 * 所以，需要把相应的 dp 设置成 true
	 */
	for j := 1; j < pSize && dp[0][j-1]; j += 2 {
		if p[j] == '*' {
			dp[0][j+1] = true
		}
	}

	for i := 0; i < sSize; i++ {
		for j := 0; j < pSize; j++ {
			if p[j] == '.' || p[j] == s[i] {
				/* p[j] 与 s[i] 可以匹配上，所以，只要前面匹配，这里就能匹配上 */
				dp[i+1][j+1] = dp[i][j]
			} else if p[j] == '*' {
				/* 此时，p[j] 的匹配情况与 p[j-1] 的内容相关。 */
				if p[j-1] != s[i] && p[j-1] != '.' {
					/**
					 * p[j] 无法与 s[i] 匹配上
					 * p[j-1:j+1] 只能被当做 ""*
					 */
					dp[i+1][j+1] = dp[i+1][j-1]
				} else {
					/**
					 * p[j] 与 s[i] 匹配上
					 * p[j-1;j+1] 作为 "x*", 可以有三种解释
					 */
					dp[i+1][j+1] = dp[i+1][j-1] || /* "x*" 解释为 "" */
						dp[i+1][j] || /* "x*" 解释为 "x" */
						dp[i][j+1] /* "x*" 解释为 "xx..." */
				}
			}
		}
	}

	return dp[sSize][pSize]
}

func isMatch3(s string, p string) bool {
	i := len(s)
	j := len(p)
	if i == 0 && j == 0 {
		return true
	}
	if i != 0 && j == 0 {
		return false
	}
	firstMatch := len(s) != 0 && (s[0] == p[0] || p[0] == '.')
	if j >= 2 && p[1] == '*' {
		return isMatch3(s, p[2:]) || (firstMatch && isMatch3(s[1:], p))
	} else {
		return firstMatch && isMatch3(s[1:], p[1:])
	}
}

func main() {
	fmt.Println(isMatch3("a", "a*b"))
}
