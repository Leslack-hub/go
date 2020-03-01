package main

import (
	"fmt"
	"math"
)

func uniquePathsWithObstacles(obstacleGrid [][]int) int {
	row := len(obstacleGrid)
	if row == 0 {
		return 0
	}
	col := len(obstacleGrid[0])
	if col == 0 {
		return 0
	}
	f := make([][]int, row)
	for i := 0; i < row; i++ {
		f[i] = make([]int, col)
	}
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			if obstacleGrid[i][j] == 1 {
				f[i][j] = 0
				continue
			}
			if i == 0 && j == 0 {
				f[i][j] = 1
				continue
			}
			if i > 0 {
				f[i][j] += f[i-1][j]
			}
			if j > 0 {
				f[i][j] += f[i][j-1]
			}
		}
	}
	return f[row-1][col-1]
}

func paintHouse(num int, cost [][]int) int {
	array := make([][]int, num+1)
	for i := 0; i < num; i++ {
		array[i] = make([]int, 3)
	}
	if num < 2 {
		return 0
	}
	for i := 1; i <= num; i++ {
		for j := 0; j < 3; j++ {
			array[i][j] = math.MaxInt32
			for k := 0; k < 3; k++ {
				if j != k {
					lastVal := array[i-1][k] + cost[i-1][j]
					if array[i][j] > lastVal {
						array[i][j] = lastVal
					}
				}
			}
		}
	}
	min := math.MaxInt32
	for _, v := range array[num] {
		if min > v {
			min = v
		}
	}
	return min
}

func main() {
	fmt.Println(uniquePathsWithObstacles([][]int{
		{0, 0, 0},
		{0, 1, 0},
		{0, 0, 0},
	}))
}
