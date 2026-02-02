package topic64_最小路径之和

import "leslack/src/helper"

func minPathSum(grid [][]int) int {
	length := len(grid)
	f := make([][]int, length)
	for i := 0; i < length; i++ {
		f[i] = make([]int, len(grid[i]))
	}

	for i := 1; i < length; i++ {
		for j := 1; j < len(grid[i]); j++ {
			f[i][j] = helper.Min(f[i-1][j], f[i][j-1]) + grid[i][j]
		}
	}

	return f[length-1][len(grid[0])-1]
}
