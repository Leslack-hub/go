package main

import (
	"fmt"
	"sort"
)

func combinationSum2(candidates []int, target int) [][]int {
	var res [][]int
	sort.Ints(candidates)
	backtrack(candidates, target, &res, 0, []int{})
	return res
}

func backtrack(candidates []int, target int, res *[][]int, sum int, path []int) {
	if target == sum {
		*res = append(*res, path)
	} else {
		tempNum := -1
		for i := 0; i < len(candidates); i++ {
			if tempNum == candidates[i] {
				continue
			} else {
				tempNum = candidates[i]
			}
			tempSum := sum + candidates[i]
			if tempSum > target {
				break
			}
			path = append(path, candidates[i])
			record := make([]int, len(path))
			copy(record, path)
			backtrack(candidates[i+1:], target, res, tempSum,record)
			path = path[:len(path)-1]
		}
	}
}

func inArray(array []int, num int) bool {
	inArray := false
	for i := range array {
		if array[i] == num {
			inArray = true
			break
		}
	}
	return inArray
}

func main() {
	fmt.Println(combinationSum2([]int{3,1,3,5,1,1}, 8))
}
