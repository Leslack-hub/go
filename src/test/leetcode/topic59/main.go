package main

import "fmt"

func generateMatrix(n int) [][]int {
	matrix := make([][]int, n)
	u, d, l, r := 0, n-1, 0, n-1
	for i := 0; i < len(matrix); i++ {
		matrix[i] = make([]int, n)
	}
	num := 1
	for num <= n*n {
		for i := l; i <= r; i++ {
			matrix[u][i] = num
			num++
		}

		u++
		for i := u; i <= d; i++ {
			matrix[i][r] = num
			num++
		}

		r--

		for i := r; i >= l; i-- {
			matrix[d][i] = num
			num++
		}

		d--

		for i := d; i >= u; i-- {
			matrix[i][l] = num
			num++
		}

		l++
	}

	return matrix
}
func main() {
	fmt.Println(generateMatrix(3))
}
