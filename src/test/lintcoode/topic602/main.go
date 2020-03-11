package main

import (
	"fmt"
	"leslack/src/helper"
	"sort"
)

func dollEnvelopes(p [][]int) int {
	sort.Slice(p, func(i, j int) bool {
		return p[i][0] < p[j][0]
	})
	length := len(p)
	f := make([]int, length)
	var max int
	for i := 0; i < length; i++ {
		f[i] = 1
		for j := 0; j < i; j++ {
			if p[j][0] < p[i][0] &&
				p[j][1] < p[i][1] {
				f[i] = helper.Max(f[i], f[j]+1)
			}
		}
		max = helper.Max(max, f[i])
	}

	return max
}

func main() {
	fmt.Println(dollEnvelopes([][]int{
		{5, 4},
		{6, 4},
		{6, 7},
		{2, 3},
	}))
}
