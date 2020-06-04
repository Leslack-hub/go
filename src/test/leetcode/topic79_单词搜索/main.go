package main

var directions = [][]int{
	{0, -1}, // 左
	{0, 1},  // 右
	{1, 0},  // 下
	{-1, 0}, // 上
}

func exist(board [][]byte, word string) bool {

}

func main() {
	exist([][]byte{
		{'A', 'B', 'C', 'E'},
		{'S', 'F', 'C', 'S'},
		{'A', 'D', 'E', 'E'},
	}, "SEE")
}
