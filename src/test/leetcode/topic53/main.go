package main

import (
	"fmt"
)

func maxSubArray(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	currentSum, maxSum := nums[0], nums[0]
	for i := 1; i < len(nums); i++ {
		currentSum = max(nums[i], currentSum+nums[i])
		maxSum = max(currentSum, maxSum)
	}
	return maxSum
}
func main() {
	fmt.Println(maxSubArray([]int{-1}))
}
