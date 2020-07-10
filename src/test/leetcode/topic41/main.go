package main

import (
	"fmt"
	"sort"
)

func firstMissingPositive(nums []int) int {
	if len(nums) == 0 {
		return 1
	}
	sort.Ints(nums)
	q := 1
	for i := 0; i < len(nums); i++ {
		if q == nums[i] {
			q++
		}
	}
	return q

}

func firstMissingPositive2(nums []int) int {
	length := len(nums)
	for i := 0; i < length; i++ {
		for nums[i] != i+1 && nums[i] < length {
			nums[nums[i]-1], nums[i] = nums[i], nums[nums[i]-1]
		}
	}
	for i := 0; i < length; i++ {
		if nums[i] != i+1 {
			return i + 1
		}
	}
	return 0
}

func main() {
	fmt.Println(firstMissingPositive2([]int{7, 8, 9, 11, 12}))
}
