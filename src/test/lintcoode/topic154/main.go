package main

import "fmt"

/**
 * 向后正则匹配
 */
func RegularExpressingMatching(A string, B string) bool {
	m, n := len(A), len(B)
	f := make([][]bool, m+1)
	for i := 0; i <= m; i++ {
		f[i] = make([]bool, n+1)
	}
	for i := 0; i <= m; i++ {
		for j := 0; j <= n; j++ {
			if i == 0 && j == 0 {
				f[i][j] = true
				continue
			}
			if i > 0 && j > 0 && B[j-1] != '*' {
				if B[j-1] == '.' || A[i-1] == B[j-1] {
					f[i][j] = f[i-1][j-1]
				}
			} else {
				if j > 1 {
					f[i][j] = f[i][j] || f[i][j-2]
				}

				if i > 0 && j > 1 {
					if B[j-2] == '.' || B[j-2] == A[i-1] {
						fmt.Println(f)
						f[i][j] = f[i][j] || f[i-1][j]
					}
				}
			}
		}
	}
	return f[m][n]
}

func main() {
	fmt.Println(RegularExpressingMatching("aa", "a*"))
}
