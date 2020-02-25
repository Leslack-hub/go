package main

import (
	"fmt"
	"math"
)

func maxSubArray(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	currentSum, maxSum := nums[0], nums[0]
	for i := 1; i < len(nums); i++ {
		currentSum = int(math.Max(float64(nums[i]), float64(currentSum+nums[i])))
		maxSum = int(math.Max(float64(currentSum), float64(maxSum)))
	}
	return maxSum
}
func main() {
	fmt.Println(maxSubArray([]int{-1}))
}
