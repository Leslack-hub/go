package main

import (
	"fmt"
	"sort"
)

func fourSum(nums []int, target int) [][]int {
	if len(nums) < 4 {
		return nil
	}
	sort.Ints(nums)
	fmt.Println(nums)
	res := [][]int{}
	max := len(nums)
	for i := 0; i < len(nums); i++ {
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		for j := i + 1; j < len(nums); j++ {
			L := j + 1
			R := max - 1
			for L < R {
				val := nums[i] + nums[j] + nums[L] + nums[R]
				if val == target {
					res = append(res, []int{nums[i], nums[j], nums[L], nums[R]})
					for j < R && nums[j] == nums[j+1] {
						j++
					}
					for L < R && nums[L] == nums[L+1] {
						L++
					}

					for L < R && nums[R] == nums[R-1] {
						R--
					}
					L++
					R--
				} else if val > target {
					R--
				} else {
					L++
				}
			}
		}
	}
	return res
}

func main() {
	fmt.Println(fourSum([]int{-3, -2, -1, 0, 0, 1, 2, 3}, 0))
}
