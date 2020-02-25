package main

import (
	"fmt"
	"strings"
)

func solveNQueens(n int) [][]string {
	board := make([][]string, n)
	for i := 0; i < n; i++ {
		board[i] = make([]string, n)
		for j := 0; j < len(board[i]); j++ {
			board[i][j] = "."
		}
	}
	var res [][]string
	backtrack(board, 0, &res)

	return res
}

func backtrack(board [][]string, col int, res *[][]string) [][]string {
	if col >= len(board) {
		var temp []string
		for i := range board {
			temp = append(temp, strings.Join(board[i], ""))
		}
		*res = append(*res, temp)
		return *res
	}

	for i := 0; i < len(board[col]); i++ {
		if !judgeReasonable(board, col, i) {
			continue
		}
		board[col][i] = "Q"
		record := make([][]string, len(board))
		copy(record, board)
		backtrack(record, col+1, res)
		board[col][i] = "."
	}
	return *res
}

func judgeReasonable(board [][]string, col int, row int) bool {
	for i := 0; i < len(board); i++ {
		if board[col][i] == "Q" {
			return false
		}
		if board[i][row] == "Q" {
			return false
		}
	}
	var tempCol, tempRow int
	if col >= row {
		tempRow = 0
		tempCol = col - row
	} else {
		tempCol = 0
		tempRow = row - col
	}

	for tempCol < len(board) && tempRow < len(board) {
		if board[tempCol][tempRow] == "Q" {
			return false
		} else {
			tempCol++
			tempRow++
		}
	}

	if col+row < len(board) {
		tempCol = col + row
		tempRow = 0
	} else {
		tempCol = len(board) - 1
		tempRow = col + row - tempCol
	}

	for tempCol >= 0 && tempRow < len(board) {
		if board[tempCol][tempRow] == "Q" {
			return false
		} else {
			tempCol--
			tempRow++
		}
	}

	return true
}

func main() {
	//fmt.Println(judgeReasonable([][]string{
	//	{".", ".", "Q", "."},
	//	{".", ".", ".", "."},
	//	{".", ".", ".", "Q"},
	//	{".", ".", ".", "."},
	//}, 1, 1))
	fmt.Println(len(solveNQueens(8)))
}
