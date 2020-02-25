package main

import (
	"fmt"
)

func canJump(nums []int) bool {
	k := 0
	for i := 0; i < len(nums); i++ {
		if i > k {
			return false
		}
		if i+nums[i] > k {
			k = i + nums[i]
		}
	}
	return true
}
func main() {
	fmt.Println(canJump([]int{3, 2, 1, 1, 4}))
}
