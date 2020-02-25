package main

func solveSudoku(board [][]byte) {
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
			d := board[i][j]

		}
	}
	box_index := func(row int, col int) int {
		return (row/3)*3 + col/3
	}

}
func place_number(d, row, col int) {
}
func main() {
	solveSudoku([][]byte{
		{'.', '2', '3', '.', '7', '.', '.', '.', '.'},
		{'4', '5', '6', '1', '9', '8', '.', '.', '.'},
		{'7', '8', '9', '.', '.', '.', '.', '6', '.'},
		{'8', '.', '.', '.', '6', '.', '.', '.', '3'},
		{'5', '.', '.', '8', '.', '3', '.', '.', '1'},
		{'9', '.', '.', '.', '2', '.', '.', '.', '6'},
		{'.', '6', '.', '.', '.', '.', '2', '8', '.'},
		{'.', '.', '.', '4', '1', '9', '.', '.', '5'},
		{'.', '.', '.', '.', '8', '.', '.', '7', '9'},
	})
}
