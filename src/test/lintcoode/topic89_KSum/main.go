package main

import (
	"fmt"
)

func KSum(A []int, K int, Target int) int {
	n := len(A)
	f := make([][][]int, n+1)
	for i := 0; i <= n; i++ {
		f[i] = make([][]int, K+1)
		for j := 0; j <= K; j++ {
			f[i][j] = make([]int, Target+1)
		}
	}

	f[0][0][0] = 1
	for i := 1; i <= n; i++ {
		for j := 0; j <= K; j++ {
			for k := 0; k <= Target; k++ {
				f[i][j][k] = f[i-1][j][k]
				if j > 0 && A[i-1] <= k {
					f[i][j][k] += f[i-1][j-1][k-A[i-1]]
				}
			}
		}
	}
	return f[n][K][Target]
}

func main() {
	fmt.Println(KSum([]int{1, 2, 3, 4}, 2, 5))
}
