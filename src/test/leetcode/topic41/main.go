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
func main() {
	fmt.Println(firstMissingPositive([]int{7,8,9,11,12}))
}
