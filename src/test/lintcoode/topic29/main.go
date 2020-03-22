package main

import "fmt"

func InterleavingString(A string, B string, X string) bool {
	m, n := len(A), len(B)
	f := make([][]bool, m+1)
	for i := 0; i <= m; i++ {
		f[i] = make([]bool, n+1)
	}

	for i := 0; i <= m; i++ {
		for j := 0; j <= n; j++ {
			if i == 0 || j == 0 {
				f[i][j] = true
				continue
			}
			if i > 0 && X[i+j-1] == A[i-1] {
				f[i][j] = f[i][j] || f[i-1][j]
			}
			if j > 0 && X[i+j-1] == B[j-1] {
				f[i][j] = f[i][j] || f[i][j-1]
			}
		}
	}
	return f[m][n]
}
func main() {
	fmt.Println(InterleavingString("aabc", "bbac", "aabcbbac"))
}
