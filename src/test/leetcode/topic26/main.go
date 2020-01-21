package main

import "fmt"

func removeDuplicates(nums []int) int {
	if len(nums) < 1 {
		return 0
	}
	temp := nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i] == temp {
			return removeDuplicates(append(nums[:i], nums[i+1:]...))
		} else {
			temp = nums[i]
		}
	}
	return len(nums)
}

// 双指针
func removeDuplicates2(nums []int) int {
	if nums == nil || len(nums) < 1 {
		return 0
	}
	p := 0
	q := 1
	for q < len(nums) {
		if nums[q] != nums[p] {
			nums[p+1] = nums[q]
			p++
		}
		q++
	}

	return p + 1
}

func main() {
	fmt.Println(removeDuplicates2([]int{1, 1, 2, 2, 3, 4, 5}))
}
