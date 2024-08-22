package main

import "fmt"

func spiralOrder(matrix [][]int) []int {
	if len(matrix) == 0 {
		return []int{}
	}
	u, d, l, r := 0, len(matrix)-1, 0, len(matrix[0])-1
	var res []int
	for {
		for i := l; i <= r; i++ {
			res = append(res, matrix[u][i])
		}

		if u < d {
			u++
		} else {
			break
		}

		for i := u; i <= d; i++ {
			res = append(res, matrix[i][r])
		}

		if r > l {
			r--
		} else {
			break
		}

		for i := r; i >= l; i-- {
			res = append(res, matrix[d][i])
		}

		if d > u {
			d--
		} else {
			break
		}

		for i := d; i >= u; i-- {
			res = append(res, matrix[i][l])
		}

		if l < r {
			l++
		} else {
			break
		}
	}
	return res
}

func spiralOrder2(matrix [][]int) []int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return []int{}
	}

	rows, cols := len(matrix), len(matrix[0])
	total := rows * cols
	result := make([]int, 0, total)

	top, bottom, left, right := 0, rows-1, 0, cols-1

	for len(result) < total {
		// 从左到右
		for i := left; i <= right; i++ {
			result = append(result, matrix[top][i])
		}
		top++

		// 从上到下
		for i := top; i <= bottom; i++ {
			result = append(result, matrix[i][right])
		}
		right--

		// 从右到左
		for i := right; i >= left; i-- {
			result = append(result, matrix[bottom][i])
		}
		bottom--

		// 从下到上
		for i := bottom; i >= top; i-- {
			result = append(result, matrix[i][left])
		}
		left++
	}

	return result
}

func main() {
	fmt.Println(spiralOrder2([][]int{
		{1, 2, 3, 10},
		{4, 5, 6, 20},
		{7, 8, 9, 30},
		{17, 18, 19, 130},
		{27, 28, 29, 230},
		{37, 38, 39, 330},
		{47, 48, 49, 430},
		{57, 58, 59, 530},
	}))
}
