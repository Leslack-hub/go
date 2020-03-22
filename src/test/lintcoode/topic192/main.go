package main

import "fmt"


func WildcardMatching(A string, B string) bool {
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
			if j > 0 &&
				B[j-1] != '*' {
				if i > 0 && (A[i-1] == B[j-1] || B[j-1] == '?') {
					f[i][j] = f[i][j] || f[i-1][j-1]
				}
			} else {
				if j > 0 {
					f[i][j] = f[i][j] || f[i][j-1]
				}
				if i > 0 {
					f[i][j] = f[i][j] || f[i-1][j]
				}
			}
		}
	}
	return f[m][n]
}

func main() {
	fmt.Println(WildcardMatching("aaab", "a*"))
}
