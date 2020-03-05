package main

import (
	"fmt"
	"leslack/src/helper"
	"math"
)

func maxProfit(k int, prices []int) int {
	length := len(prices)
	// 表示可以 无限制次数
	if k > length/2 {
		var max int
		for i := 1; i < length; i++ {
			num := prices[i] - prices[i-1]
			if num > 0 {
				max += num
			}
		}
		return max
	}

	nums := 2 * k
	f := make([][]int, length+1)
	for i := 0; i <= length; i++ {
		f[i] = make([]int, nums+2)
	}
	f[0][1] = 0
	for i := 2; i <= nums+1; i++ {
		f[0][i] = math.MinInt32
	}
	for i := 1; i <= length; i++ {
		for j := 1; j <= nums+1; j += 2 {
			f[i][j] = f[i-1][j]
			if j > 1 && i >= 2 && f[i-1][j-1] != math.MinInt32 {
				f[i][j] = helper.Max(f[i][j], f[i-1][j-1]+prices[i-1]-prices[i-2])
			}
		}

		for j := 2; j <= nums; j += 2 {
			f[i][j] = f[i-1][j-1]
			if i >= 2 && f[i-1][j] != math.MinInt32 {
				f[i][j] = helper.Max(f[i][j], f[i-1][j]+prices[i-1]-prices[i-2])
			}
		}
	}
	var max int
	for i := 1; i <= nums+1; i += 2 {
		max = helper.Max(max, f[length][i])
	}
	return max
}

func main() {
	fmt.Println(maxProfit(3, []int{4, 4, 6, 1, 1, 4, 2, 5}))
}
