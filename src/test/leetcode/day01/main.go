package main

import "fmt"

func twoSum(nums []int, target int) []int {
	num := len(nums)
	for k, v := range nums {
		for i := k + 1; i < num; i++ {
			if v+nums[i] == target {
				return []int{k, i}
			}
		}
	}
	return []int{0, 0}
}

func main() {
	fmt.Println(twoSum([]int{3, 4, 5, 1, 2}, 9))
}
