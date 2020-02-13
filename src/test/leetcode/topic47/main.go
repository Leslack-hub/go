package main

import "fmt"

func permuteUnique(nums []int) [][]int {
	var res [][]int
	q := &quest{
		nums: nums,
		res:  &res,
	}
	q.backtrack(nums, 0)
	return *q.res
}

type quest struct {
	nums []int
	res  *[][]int
}

func (q *quest) backtrack(nums []int, first int) {
	if first == len(nums) {
		*q.res = append(*q.res, nums)
		return
	}
	var used []int
	for i := first; i < len(nums); i++ {
		record := make([]int, len(nums))
		copy(record, nums)
		if inArray(used, record[i]) {
			continue
		}
		record[first], record[i] = record[i], record[first]
		q.backtrack(record, first+1)
		record[first], record[i] = record[i], record[first]
		used = append(used, record[i])
	}
}

func inArray(array []int, val int) bool {
	var inArray bool
	for i := range array {
		if array[i] == val {
			inArray = true
			break
		}
	}
	return inArray
}

func main() {
	fmt.Println(permuteUnique([]int{1, 1, 2}))
}
