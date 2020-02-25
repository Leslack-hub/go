package main

import (
	"fmt"
	"sort"
)

// Interval Definition for an interval.
type Interval struct {
	Start int
	End   int
}

func merge(its [][]int) [][]int {
	if len(its) <= 1 {
		return its
	}

	var array []Interval
	for k := range its {
		array = append(array, Interval{its[k][0], its[k][1]})
	}

	// 按照start 排序
	sort.Slice(array, func(i int, j int) bool {
		return array[i].Start < array[j].Start
	})

	res := make([][]int, 0, len(array))

	temp := array[0]
	for i := 1; i < len(array); i++ {
		if array[i].Start <= temp.End {
			temp.End = max(temp.End, array[i].End)
		} else {
			res = append(res, []int{temp.Start, temp.End})
			temp = array[i]
		}
	}
	res = append(res, []int{temp.Start, temp.End})
	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func main() {
	fmt.Println(merge([][]int{
		{1, 3}, {15, 18}, {8, 10}, {4, 6},
	}))
}
