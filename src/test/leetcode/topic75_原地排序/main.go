package main

import "fmt"

func sortColors(nums []int) {
	p0, cur := 0, 0
	p2 := len(nums) - 1
	for cur <= p2 {
		if nums[cur] == 0 {
			nums[p0], nums[cur] = nums[cur], nums[p0]
			p0++
			cur++
		} else if nums[cur] == 2 {
			nums[cur], nums[p2] = nums[p2], nums[p0]
			p2--
		} else {
			cur++
		}
	}
	fmt.Println(nums)
}

func main() {
	sortColors([]int{2, 0, 2, 1, 1, 0})
}
