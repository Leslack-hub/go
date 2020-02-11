package main

import (
	"fmt"
	"math"
)

func trap(height []int) int {
	res := 0
	for i := 1; i < len(height)-1; i++ {
		leftMax, rightMax := 0, 0
		for j := i; j > 0; j-- {
			leftMax = int(math.Max(float64(leftMax), float64(height[j])))
		}

		for j := i; j < len(height)-1; j++ {
			rightMax = int(math.Max(float64(rightMax), float64(height[j])))
		}
		res += int(math.Min(float64(leftMax), float64(rightMax))) - height[i]
	}
	return res
}
func trap2(height []int) int {
	left, right := 0, len(height)-1
	res, leftMax, rightMax := 0, 0, 0
	for left < right {
		if height[left] < height[right] {
			if height[left] >= leftMax {
				leftMax = height[left]
			} else {
				res += leftMax - height[left]
			}
			left++
		} else {
			if height[right] >= rightMax {
				rightMax = height[right]
			} else {
				res += rightMax - height[right]
			}
			right--
		}
	}

	return res
}
func trap3(height []int) int {
	var res int
	maxNum := 0
	for i := range height {
		if maxNum < height[i] {
			maxNum = height[i]
		}
	}
	for i := 0; i < maxNum; i++ {
		isStart, sum := false, 0
		for j := 0; j < len(height); j++ {
			if isStart && height[j] < i {
				sum++
			}
			if height[j] >= i {
				res += sum
				sum = 0
				isStart = true
			}
		}
	}
	return res
}
func main() {
	fmt.Println(trap3([]int{6, 4, 2, 0, 3, 2, 0, 3, 1, 4, 5, 3, 2, 7, 5, 3, 0, 1, 2, 1, 3, 4, 6, 8, 1, 3}))
}
