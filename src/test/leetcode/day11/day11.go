package main

import "fmt"

func maxArea(height []int) int {
	if len(height) < 2 {
		return height[0]
	}
	var area, a int
	for k, v := range height {
		for j := k + 1; j < len(height); j++ {
			if height[j] > v {
				a = v * (j - k)
			} else {
				a = height[j] * (j - k)
			}
			if a > area {
				area = a
			}
		}
	}
	return area
}
func main() {
	fmt.Println(maxArea([]int{1, 8, 6, 2, 5, 4, 8, 3, 7}))
}
