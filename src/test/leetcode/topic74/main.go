package main

import "fmt"

func searchMatrix(matrix [][]int, target int) bool {
	column := len(matrix)
	if column == 0 {
		return false
	}
	row := len(matrix[0])
	if row == 0 {
		return false
	}

	for i := 0; i < column; i++ {
		if matrix[i][row-1] >= target {
			for j := 0; j < row; j++ {
				if matrix[i][j] == target {
					return true
				}
			}
			return false
		}
	}
	return false
}

func main() {
	fmt.Println(searchMatrix([][]int{
		{1, 3, 5, 7},
		{10, 11, 16, 20},
		{23, 30, 34, 50},
	}, 24))
}
