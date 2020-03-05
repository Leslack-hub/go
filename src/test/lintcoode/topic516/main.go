package main

import (
	"fmt"
	"math"
)

func paintHouse(n int, k int, cost [][]int) int {
	f := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		f[i] = make([]int, k)
	}
	min1, min2, id1, id2 := 0, 0, 0, 0
	for i := 1; i <= n; i++ {
		min1, min2 = math.MaxInt32, math.MaxInt32
		for j := 0; j < k; j++ {
			if f[i-1][j] < min1 {
				min2 = min1
				id2 = id1
				min1 = f[i-1][j]
				id1 = j
			} else {
				if f[i-1][j] < min2 {
					min2 = f[i-1][j]
					id2 = j
				}
			}
		}
		for z := 0; z < k; z++ {
			f[i][z] = cost[i-1][z]
			if z != id1 {
				f[i][z] += min1
			} else {
				f[i][z] += min2
			}
		}
		fmt.Println(id2)
	}
	res := math.MaxInt32
	for i := 0; i < k; i++ {
		if f[n][i] < res {
			res = f[n][i]
		}
	}
	return res

}

func main() {
	fmt.Println(paintHouse(3, 3, [][]int{
		{14, 2, 11},
		{11, 14, 5},
		{14, 3, 10},
	}))
}
