package main

import "fmt"

var b [][]byte
var w string

func exist(board [][]byte, word string) bool {
	m, n := len(board), len(board[0])
	b = board
	w = word
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if dfs(i, j, 0) {
				return true
			}
		}
	}
	return false
}

func dfs(i, j, k int) bool {
	if k == len(w) {
		return true
	}

	if i < 0 ||
		j < 0 ||
		i == len(b) ||
		j == len(b[i]) {
		return false
	}
	if b[i][j] != w[k] {
		return false
	}
	fmt.Println(i, j, k)
	tmp := b[i][j]
	b[i][j] = byte(' ')

	if dfs(i-1, j, k+1) {
		return true
	}
	if dfs(i+1, j, k+1) {
		return true
	}
	if dfs(i, j-1, k+1) {
		return true
	}
	if dfs(i, j+1, k+1) {
		return true
	}
	b[i][j] = tmp
	return false
}
func main() {
	fmt.Println(exist([][]byte{
		{'A', 'B', 'C', 'E'},
		{'S', 'F', 'C', 'S'},
		{'A', 'D', 'E', 'E'},
	}, "ACE"))
}
