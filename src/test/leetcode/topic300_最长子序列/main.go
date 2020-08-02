package main

import (
	"fmt"
)

func lengthOfLIS(nums []int) int {
	length := len(nums)
	if length == 0 {
		return 0
	}
	if length == 1 {
		return 1
	}
	dp := make([]int, length)
	for i := 0; i < length; i++ {
		dp[i] = 1
	}
	for i := 1; i < length; i++ {
		for j := 0; j < i; j++ {
			if nums[i] > nums[j] {
				if dp[i] < dp[j]+1 {
					dp[i] = dp[j] + 1
				}
			}
		}
	}
	max := 0
	for i := 0; i < length; i++ {
		if max < dp[i] {
			max = dp[i]
		}
	}
	return max
}

func main() {
	fmt.Println(lengthOfLIS([]int{10, 9, 2, 5, 3, 7, 101, 18}))
}
