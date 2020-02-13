package main

import (
	"fmt"
)

func isValidSudoku(board [][]byte) bool {
	size := len(board)
	var res bool
	for i := 0; i < size; i++ {
		res = check(board, i)
		if res == false {
			return false
		}
	}
	return res
}

func check(bytes [][]byte, i int) bool {
	var m = map[byte]int{'1': 0, '2': 0, '3': 0, '4': 0, '5': 0, '6': 0, '7': 0, '8': 0, '9': 0,}
	var m1 = map[byte]int{'1': 0, '2': 0, '3': 0, '4': 0, '5': 0, '6': 0, '7': 0, '8': 0, '9': 0,}
	var m2 = map[byte]int{'1': 0, '2': 0, '3': 0, '4': 0, '5': 0, '6': 0, '7': 0, '8': 0, '9': 0,}
	for k, v := range bytes[i] {
		item, ok := m[v]
		if ok {
			if item >= 1 {
				return false
			} else {
				m[v]++
			}
		}
		// 检查竖排
		if i == 0 {
			m1 = resetList(m1)
			for j := 0; j < 9; j++ {
				item, ok := m1[bytes[j][k]]
				if ok {
					if item >= 1 {
						fmt.Println(j, k)
						return false
					} else {
						m1[bytes[j][k]]++
					}
				}
			}
		}

		// 周围检查
		if i%3 == 0 {
			for x := 0; x < 9; x += 3 {
				m2 = resetList(m2)
				for y := 0; y < 3; y++ {
					for z := 0; z < 3; z++ {
						item, ok := m2[bytes[i+y][x+z]]
						if ok {
							if item >= 1 {
								fmt.Println(i+y, x+z)
								return false
							} else {
								m2[bytes[i+y][x+z]]++
							}
						}
					}
				}
			}
		}
	}

	return true
}

func resetList(m map[byte]int) map[byte]int {
	for k := range m {
		m[k] = 0
	}
	return m
}

func isValidSudoku2(board [][]byte) bool {
	rows := make(map[int]map[byte]int, 9)
	columns := make(map[int]map[byte]int, 9)
	boxes := make(map[int]map[byte]int, 9)
	for i := 0; i < 9; i++ {
		rows[i] = make(map[byte]int, 9)
		columns[i] = make(map[byte]int, 9)
		boxes[i] = make(map[byte]int, 9)
	}

	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			num := board[i][j]
			if num != '.' {
				n := num
				box_index := (i/3)*3 + j/3
				rows[i][n] = rows[i][n] + 1
				columns[j][n] = columns[j][n] + 1
				boxes[box_index][n] = boxes[box_index][n] + 1
				if rows[i][n] > 1 ||
					columns[j][n] > 1 ||
					boxes[box_index][n] > 1 {
					return false
				}
			}
		}
	}
	return true
}

func main() {
	fmt.Println(isValidSudoku2([][]byte{
		{'.', '2', '3', '.', '7', '.', '.', '.', '.'},
		{'4', '5', '6', '1', '9', '8', '.', '.', '.'},
		{'7', '8', '9', '.', '.', '.', '.', '6', '.'},
		{'8', '.', '.', '.', '6', '.', '.', '.', '3'},
		{'5', '.', '.', '8', '.', '3', '.', '.', '1'},
		{'9', '.', '.', '.', '2', '.', '.', '.', '6'},
		{'.', '6', '.', '.', '.', '.', '2', '8', '.'},
		{'.', '.', '.', '4', '1', '9', '.', '.', '5'},
		{'.', '.', '.', '.', '8', '.', '.', '7', '9'},
	}))
}
