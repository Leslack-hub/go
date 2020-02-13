package main

import (
	"fmt"
	"sort"
)

func combinationSum(candidates []int, target int) [][]int {
	var res [][]int
	sort.Ints(candidates)
	calculate(candidates, target, &res, 0, []int{})
	return res
}

func calculate(list []int, target int, res *[][]int, sum int, record []int) {
	if target == sum {
		*res = append(*res, record)
		return
	} else {
		for i := range list {
			temp := sum + list[i]
			if temp > target {
				break
			}
			tempRecord := record
			tempRecord = append(tempRecord, list[i])
			calculate(list[i:], target, res, temp, tempRecord)
		}
	}
}

func combinationSum2(candidates []int, target int) [][]int {
	sort.Ints(candidates)
	size := len(candidates)
	var path []int
	var res [][]int
	if size == 0 {
		return res
	}
	backtrack(candidates, 0, size, path, &res, target)
	return res
}

type quest struct {
	candidates [] int
	target     int
	res        *[][]int
}

func backtrack(candidates []int, begin int, size int, path []int, res *[][]int, target int) {
	if target == 0 {
		*res = append(*res, path)
	}
	for index := begin; index < size; index++ {
		residue := target - candidates[index]
		if residue < 0 {
			break
		}
		path = append(path, candidates[index])
		record := make([]int, len(path))
		copy(record, path)
		backtrack(candidates, index, size, record, res, residue)
		path = path[:len(path)-1]
	}
}

func combinationSum3(candidates []int, target int) [][]int {
	sort.Ints(candidates)

	res := [][]int{}
	solution := []int{}
	cs(candidates, solution, target, &res)

	return res
}

func cs(candidates, solution []int, target int, result *[][]int) {
	if target == 0 {
		*result = append(*result, solution)
	}

	if len(candidates) == 0 || target < candidates[0] {
		// target < candidates[0] 因为candidates是排序好的
		return
	}

	// 这样处理一下的用意是，让切片的容量等于长度，以后append的时候，会分配新的底层数组
	// 避免多处同时对底层数组进行修改，产生错误的答案。
	// 可以注释掉以下语句，运行单元测试，查看错误发生。
	solution = solution[:len(solution):len(solution)]

	cs(candidates, append(solution, candidates[0]), target-candidates[0], result)

}

func combinationSum4(candidates []int, target int) [][]int {
	record := make(map[int][][]int)
	for i := 1; i < target+1; i++ {
		record[i] = make([][]int, 1)
	}

	for i := 1; i < target+1; i++ {
		for j := range candidates {
			if i == candidates[j] {
				temp := []int{candidates[j]}
				record[i] = append(record[i], temp)
			} else if i > candidates[j] {
				for _, k := range record[i-candidates[j]] {
					k = append(k, candidates[j])
					sort.Ints(k)
					inArray := false
					for _, v := range record[i] {
						for a := range v {
							if a == candidates[j] {
								inArray = true
							}
						}
					}

					if !inArray {
						record[i] = append(record[i], k)
					}
				}
			}
		}
	}
	return record[target]
}

func main() {
	fmt.Println(combinationSum2([]int{2, 3, 7}, 9))
}
