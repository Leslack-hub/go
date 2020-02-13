package main

import (
	"fmt"
	"math"
)

func jump(nums []int) int {
	end, steps, maxPosition := 0, 0, 0
	for i := 0; i < len(nums)-1; i++ {
		maxPosition = int(math.Max(float64(maxPosition), float64(nums[i]+i)))
		if i == end {
			end = maxPosition
			steps++
		}
	}
	return steps
}

func main() {
	fmt.Println(jump([]int{2, 1, 1, 1, 4}))
}
