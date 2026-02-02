package main

import (
	"fmt"
	"sort"
)

var res [][]int

func subsetsWithDup(nums []int) [][]int {
	res = make([][]int, 0)
	sort.Ints(nums)
	dfs([]int{}, nums, 0)
	return res
}

func dfs(temp []int, nums []int, start int) {
	tem := make([]int, len(temp))
	copy(tem, temp)
	res = append(res, tem)
	for i := start; i < len(nums); i++ {
		if i > start && nums[i] == nums[i-1] {
			continue
		}
		temp = append(temp, nums[i])
		dfs(temp, nums, i+1)
		temp = temp[:len(temp)-1]
	}
}

func main() {
	fmt.Println(subsetsWithDup([]int{1, 2, 2}))
}
