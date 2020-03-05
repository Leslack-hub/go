package main

import (
	"fmt"
	"leslack/src/helper"
	"math"
)

/**
 * @param prices: Given an integer array
 * @return: Maximum profit
 */
func maxProfit(prices []int) int {
	length := len(prices)
	f := make([][]int, length+1)
	for i := 0; i <= length; i++ {
		f[i] = make([]int, 6)
	}
	f[0][1] = 0
	for i := 2; i <= 5; i++ {
		f[0][i] = math.MinInt32
	}
	for i := 1; i <= length; i++ {
		for j := 1; j <= 5; j += 2 {
			f[i][j] = f[i-1][j]
			if j > 1 && i >= 2 && f[i-1][j-1] != math.MinInt32 {
				f[i][j] = helper.Max(f[i][j], f[i-1][j-1]+prices[i-1]-prices[i-2])
			}

		}

		for j := 2; j < 5; j += 2 {
			f[i][j] = f[i-1][j-1]
			if i >= 2 && f[i-1][j] != math.MinInt32 {
				f[i][j] = helper.Max(f[i][j], f[i-1][j]+prices[i-1]-prices[i-2])
			}
		}
	}

	return helper.Max(f[length][1], f[length][3], f[length][5])
}

func main() {
	fmt.Println(maxProfit([]int{4, 4, 6, 1, 1, 4, 2, 5}))
}
