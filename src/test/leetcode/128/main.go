package main

import (
	"fmt"
)

func longestConsecutive(nums []int) int {
	set := make(map[int]bool)
	for _, num := range nums {
		set[num] = true
	}
	var r [][]int
	var longset int
start:
	for _, num := range nums {
		if !set[num-1] {
			for _, v := range r {
				if num >= v[0] && num < v[1] {
					continue start
				}
			}
			cur := num
			long := 1
			for set[cur+1] {
				long++
				cur++
			}
			if long > longset {
				r = append(r, []int{num, cur})
				longset = long
			}
		}
	}
	return longset
}

func main() {
	fmt.Println(longestConsecutive([]int{4, 0, -4, -2, 2, 5, 2, 0, -8, -8, -8, -8, -1, 7, 4, 5, 5, -4, 6, 6, -3}))
}
