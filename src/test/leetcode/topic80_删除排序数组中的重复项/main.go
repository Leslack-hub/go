package main

import "fmt"

func removeDuplicates(nums []int) int {
	i, j := 1, 2
	for ; j < len(nums); j++ {
		if nums[j] != nums[i-1] {
			i++
			nums[i] = nums[j]
		}
	}
	return i + 1
}

func removeDuplicates1(nums []int) int {
	j, count := 1, 1
	for i := 1; i < len(nums); i++ {
		if nums[i] == nums[i-1] {
			count++
		} else {
			count = 1
		}
		if count <= 2 {
			nums[j] = nums[i]
			j++
		}
	}
	return j
}
func main() {
	fmt.Println(removeDuplicates1([]int{1, 1, 1, 1, 1, 2, 3, 3}))
}
