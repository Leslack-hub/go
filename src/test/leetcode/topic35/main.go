package main

import "fmt"

func searchInsert(nums []int, target int) int {
	res := 0
	for i := 0; i < len(nums); i++ {
		if nums[i] == target {
			return i
		} else if nums[i] < target {
			res = i + 1
		}
	}
	return res
}
func main() {
	fmt.Println(searchInsert([]int{1, 3, 5, 6}, 2))
}
