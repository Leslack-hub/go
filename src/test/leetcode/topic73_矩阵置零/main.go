package main

import "fmt"

func setZeroes(matrix [][]int) {
	c := len(matrix)
	r := len(matrix[0])
	queue := make([][]int, 0)
	for i := 0; i < c; i++ {
		for j := 0; j < r; j++ {
			if matrix[i][j] == 0 {
				queue = append(queue, []int{i, j})
			}
		}
	}
	for i := 0; i < len(queue); i++ {
		cloumn := queue[i][0]
		row := queue[i][1]
		for j := 0; j < r; j++ {
			matrix[cloumn][j] = 0
		}
		for j := 0; j < c; j++ {
			matrix[j][row] = 0
		}
	}
	fmt.Println(matrix)
}

func main() {
	(setZeroes([][]int{
		{0, 1, 2, 0},
		{3, 4, 5, 2},
		{1, 3, 1, 5},
	}))
}
