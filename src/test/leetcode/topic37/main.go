package main

import "fmt"

func solveSudoku(board [][]byte) {
	rows := make([][]bool, 10)
	columns := make([][]bool, 10)
	boxes := make([][][]bool, 3)

	for i := 0; i < 10; i++ {
		rows[i] = make([]bool, 10)
		columns[i] = make([]bool, 10)
	}

	for i := 0; i < 3; i++ {
		boxes[i] = make([][]bool, 3)
		for j := 0; j < 3; j++ {
			boxes[i][j] = make([]bool, 10)
		}
	}
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			val := board[i][j] - '0'
			if 1 <= val && val <= 9 {
				rows[i][val] = true
				columns[j][val] = true
				boxes[i/3][j/3][val] = true
			}
		}
	}

	fmt.Println(recusiveSolveSudoku(board, rows, columns, boxes, 0, 0))
	for _, i := range board {
		for _, j := range i {
			fmt.Printf("%s",string(j))
		}
		fmt.Println()
	}
}

func recusiveSolveSudoku(board [][]byte, rows [][]bool, columns [][]bool, boxes [][][]bool, row int, col int) bool {
	if col == len(board[0]) {
		col = 0
		row++
		if row == len(board) {
			return true
		}
	}

	if board[row][col] == '.' {
		for num := 1; num <= 9; num++ {
			canUse := rows[row][num] || columns[col][num] || boxes[row/3][col/3][num]
			if !canUse {
				rows[row][num] = true
				columns[col][num] = true
				boxes[row/3][col/3][num] = true

				board[row][col] = byte(num + '0')
				if recusiveSolveSudoku(board, rows, columns, boxes, row, col+1) {
					return true
				}
				board[row][col] = '.'
				rows[row][num] = false
				columns[col][num] = false
				boxes[row/3][col/3][num] = false
			}
		}
	} else {
		return recusiveSolveSudoku(board, rows, columns, boxes, row, col+1)
	}
	return false
}

func main() {
	solveSudoku([][]byte{
		{'.', '5', '.', '.', '.', '.', '.', '2', '.'},
		{'4', '.', '.', '2', '.', '6', '.', '.', '7'},
		{'.', '.', '8', '.', '3', '.', '1', '.', '.'},
		{'.', '1', '.', '.', '.', '.', '.', '6', '.'},
		{'.', '.', '9', '.', '.', '.', '5', '.', '.'},
		{'.', '7', '.', '.', '.', '.', '.', '9', '.'},
		{'.', '.', '5', '.', '8', '.', '3', '.', '.'},
		{'7', '.', '.', '9', '.', '1', '.', '.', '4'},
		{'.', '2', '.', '.', '.', '.', '.', '7', '.'},
	})
}
