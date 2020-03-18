package main

import (
	"fmt"
	"leslack/src/helper"
)

/**
题意：
给定T个01串S0, S1, S2,... S(T-1)
现有m个0，n个1
问最多能组成多少个给定01串
每个串最多组成一次
子问题：
	0和1的个数在变化，如何记录？
	直接放入状态
	状态：设f[i][j][k]为前i个01串最多能有多少个被j个0和k个1组成
*/
func OneAndZeroes(A []string, m int, n int) int {
	length := len(A)
	f := make([][][]int, length+1)
	for i := 0; i <= length; i++ {
		f[i] = make([][]int, m+1)
		for j := 0; j <= m; j++ {
			f[i][j] = make([]int, n+1)
		}
	}
	for i := 1; i <= length; i++ {
		a0, a1 := 0, 0
		for _, v := range A[i-1] {
			if v == '0' {
				a0++
			} else {
				a1++
			}
		}
		for j := 1; j <= m; j++ {
			for k := 1; k <= n; k++ {
				f[i][j][k] = f[i-1][j][k]
				if j >= a0 && k >= a1 {
					f[i][j][k] = helper.Max(f[i][j][k], f[i-1][j-a0][k-a1]+1)
				}
			}
		}
	}
	return f[length][m][n]
}

func main() {
	fmt.Println(OneAndZeroes([]string{"10", "0001", "111001", "1", "0"}, 5, 3))
}
