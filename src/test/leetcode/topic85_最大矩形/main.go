package main

import (
	"fmt"
	"os"
)

func maximalRectangle(matrix [][]byte) int {
	length := len(matrix)
	dp := make([][]int, length)
	for i := 0; i < length; i++ {
		dp[i] = make([]int, len(matrix[i]))
	}
	var maxarea int
	for i := 0; i < length; i++ {
		for j := 0; j < len(matrix[i]); j++ {
			if matrix[i][j] == '1' {
				if j == 0 {
					dp[i][j] = 1
				} else {
					dp[i][j] = dp[i][j-1] + 1
				}

				width := dp[i][j]
				for k := i; k >= 0; k-- {
					width = min(width, dp[k][j])
					maxarea = max(maxarea, width*(i-k+1))
				}
				if maxarea == 6 {
					fmt.Println(dp)
					fmt.Println(i, j, width)
					os.Exit(1)
				}
			}
		}
	}

	return maxarea
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func min(a, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

func main() {
	fmt.Println(maximalRectangle([][]byte{
		{'1', '0', '1', '0', '0'},
		{'1', '0', '1', '1', '1'},
		{'1', '1', '1', '1', '1'},
		{'1', '0', '0', '1', '0'},
	}))
}
