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
func main() {
	fmt.Println(spiralOrder([][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}))
}
