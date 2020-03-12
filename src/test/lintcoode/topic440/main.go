package main

import (
	"fmt"
	"leslack/src/helper"
)

func backpackIII(A []int, V []int, m int) int {
	// 转义方式f[i][w] = 用前i种物品拼出w重量时的最大价值（-1）表示拼不出
	n := len(A)
	f := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		f[i] = make([]int, m+1)
	}

	f[0][0] = 0
	for i := 1; i <= m; i++ {
		f[0][i] = -1
	}

	for i := 1; i <= n; i++ {
		for w := 1; w <= m; w++ {
			f[i][w] = f[i-1][w]
			if w >= A[i-1] && f[i][w-A[i-1]] != -1 {
				f[i][w] = helper.Max(f[i][w], f[i][w-A[i-1]]+V[i-1])
			}
		}
	}
	// 同理
	for i := 1; i <= n; i++ {
		for w := 1; w <= m; w++ {
			f[i][w] = f[i-1][w]
			k := 1
			for w >= k*A[i-1] && f[i-1][w-k*A[i-1]] != -1 {
				f[i][w] = helper.Max(f[i][w], f[i-1][w-k*A[i-1]]+k*V[i-1])
				f[i-1][w] = f[i][w]
				k++
			}
		}
	}
	fmt.Println(f)
	var res int
	for i := 0; i <= m; i++ {
		res = helper.Max(res, f[n][i])
	}
	return res
}

func main() {
	fmt.Println(backpackIII([]int{2, 3, 5, 7}, []int{1, 5, 2, 4}, 10))
}
