package main

import "fmt"

func largestRectangleArea(heights []int) int {
	n := len(heights)
	left, right := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		right[i] = n
	}
	monoStack := []int{}
	for i := 0; i < n; i++ {
		for len(monoStack) > 0 && heights[monoStack[len(monoStack)-1]] >= heights[i] {
			right[monoStack[len(monoStack)-1]] = i
			monoStack = monoStack[:len(monoStack)-1]
		}
		if len(monoStack) == 0 {
			left[i] = -1
		} else {
			left[i] = monoStack[len(monoStack)-1]
		}
		monoStack = append(monoStack, i)
		fmt.Println(monoStack)
	}
	ans := 0
	for i := 0; i < n; i++ {
		ans = max(ans, (right[i]-left[i]-1)*heights[i])
	}
	return ans
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

/**
 * 暴力解法
 */
func largestRectangleArea2(heights []int) int {
	length := len(heights)
	var maxVar int
	for i := 0; i < length; i++ {
		left, right := i-1, i+1
		for right < length && heights[right] >= heights[i] {
			right++
		}
		for left >= 0 && heights[left] >= heights[i] {
			left--
		}
		maxVar = max(maxVar, heights[i]*(right-left-1))
	}
	return maxVar
}

func main() {
	fmt.Println(largestRectangleArea([]int{2, 1, 5, 6, 2, 3}))
}
