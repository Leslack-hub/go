package main

import "fmt"

func permute(nums []int) [][]int {
	var res [][]int
	if len(nums) <= 1 {
		return append(res, nums)
	}
	calculate(nums, []int{}, &res)
	return res
}
func calculate(nums []int, path []int, res *[][]int) {
	if len(nums) == 2 {
		*res = append(*res, append(path, nums...))
		nums[0], nums[1] = nums[1], nums[0]
		*res = append(*res, append(path, nums...))
		return
	}
	for i := 0; i < len(nums); i++ {
		recordNums := make([]int, len(nums))
		copy(recordNums, nums)
		path = append(path, nums[i])
		record := make([]int, len(path))
		copy(record, path)
		calculate(append(recordNums[:i], recordNums[i+1:]...), record, res)
		path = path[:len(path)-1]
	}
}
func main() {
	fmt.Println(permute([]int{1, 2, 3, 4}))
}
