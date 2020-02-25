package main

import (
	"fmt"
	"sort"
)

func insert(intervals [][]int, newInterval []int) [][]int {
	intervals = append(intervals, newInterval)
	var array []Interval
	for k := range intervals {
		array = append(array, Interval{intervals[k][0], intervals[k][1]})
	}
	sort.Slice(array, func(i, j int) bool {
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

type Interval struct {
	Start int
	End   int
}

func main() {
	fmt.Println(insert([][]int{
		{1, 3}, {6, 9},
	}, []int{2, 5}))
}
