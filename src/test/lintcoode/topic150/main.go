package main

import "fmt"

/**
 * @param prices: Given an integer array
 * @return: Maximum profit
 */
func maxProfit(prices []int) int {
	max := 0
	for i := 1; i < len(prices); i++ {
		num := prices[i] - prices[i-1]
		if num > 0 {
			max += num
		}
	}
	return max
}

func main() {
	fmt.Println(maxProfit([]int{2, 1, 2, 0, 1}))
}
