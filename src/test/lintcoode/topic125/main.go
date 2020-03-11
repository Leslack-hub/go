package main

import (
	"fmt"
	"leslack/src/helper"
)

func backpackII(A []int, P []int, m int) int {
	// n个物品
	n := len(A)
	f := make([][]int, n+1)
	// f[i][w] 前i个物品拼出重量w的最大重甲值
	for i := 0; i <= n; i++ {
		f[i] = make([]int, m+1)
	}

	f[0][0] = 0
	for i := 1; i <= m; i++ {
		f[0][i] = -1
	}

	for i := 1; i <= n; i++ {
		for w := 1; w <= m; w++ {
			// A[w]表示物品重量 P[w]表示中价值
			f[i][w] = f[i-1][w]
			if w >= A[i-1] && f[i-1][w-A[i-1]] != -1 {
				f[i][w] = helper.Max(f[i][w], f[i-1][w-A[i-1]]+P[i-1])
			}
		}
	}
	var res int
	for i := 0; i <= m; i++ {
		res = helper.Max(res, f[n][i])
	}
	return res
}
func main() {
	fmt.Println(backpackII([]int{2, 3, 5, 7}, []int{1, 5, 2, 4}, 11))
}
