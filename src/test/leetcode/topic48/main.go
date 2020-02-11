package main

func rotate(matrix [][]int) {
	m := len(matrix)
	res := make([][]int, m)
	for i := 0; i < m; i++ {
		n := len(matrix[i])
		res[i] = make([]int, n)
		for j := 0; j < n; j++ {
			res[i][j] = matrix[m-j-1][i]
		}
	}
	for i := 0; i < m; i++ {
		n := len(matrix[i])
		for j := 0; j < n; j++ {
			matrix[i][j] = res[i][j]
		}
	}
}

func main() {
	rotate([][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	})
}
