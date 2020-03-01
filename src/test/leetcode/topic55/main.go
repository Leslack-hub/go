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

func canJump2(nums []int) bool {
	length := len(nums)
	f := make([]bool, length)
	f[0] = true
	for j := 1; j < length; j++ {
		f[j] = false
		for i := 0; i < j; i++ {
			if f[i] && i+nums[i] >= j {
				f[j] = true
				break
			}
		}
	}
	return f[length-1]
}
func main() {
	fmt.Println(canJump2([]int{3, 2, 1, 0, 4}))
}
