package main

import (
	"fmt"
	"leslack/src/helper"
	"math"
)

/**
题意：
给定数组A,每个元素是不超过100的正整数
将A种每个元素修改形成数组B
要求B种任意两个相邻的元素的差不能超过target
求最小修改代价 即A[0]-B[0] + .. A[n-1]-B[n-1]
f[i][j] 为将A前i个元素改成B的最小代价，确保前i个改好的元素种任意两个相邻的元素的差不超过Target 并且A[i-1]改成j
*/
func MinimumAdjustment(A []int, target int) int {
	n := len(A)
	f := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		f[i] = make([]int, 101)
	}
	for i := 1; i <= 100; i++ {
		// 第一个数字 修改为 1 - 100 的代价
		f[1][i] = int(math.Abs(float64(A[0] - i)))
	}

	for i := 2; i <= n; i++ {
		for j := 1; j <= 100; j++ {
			f[i][j] = math.MaxInt32
			// j 改成k 的最小代价
			for k := j - target; k <= j+target; k++ {
				if k < 1 || k > 100 {
					continue
				}
				f[i][j] = helper.Min(f[i][j], f[i-1][k])
			}
			f[i][j] += int(math.Abs(float64(j - A[i-1])))
		}
	}
	res := math.MaxInt32
	for i := 1; i <= 100; i ++ {
		res = helper.Min(res, f[n][i])
	}
	fmt.Println(f)
	return res
}

func main() {
	fmt.Println(MinimumAdjustment([]int{1, 4, 2, 3}, 1))
}
