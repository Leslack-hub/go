package main

import "fmt"

func backpack(A []int, m int) int {
	n := len(A)
	f := make([][]bool, n+1)
	for i := 0; i <= n; i++ {
		f[i] = make([]bool, m+1)
	}
	f[0][0] = true
	for i := 1; i <= n; i++ {
		// i = 2 j = 5
		for j := 1; j <= m; j++ {
			// f[2][5] =f[1][5] false
			f[i][j] = f[i-1][j]
			//  5 >= A[i-1] 3
			if j >= A[i-1] {
				f[i][j] = f[i][j] || f[i-1][j-A[i-1]]
			}
		}
	}
	var res int
	for i := m; i >= 0; i-- {
		if f[n][i] {
			res = i
			break
		}
	}
	return res
}

func main() {
	fmt.Println(backpack([]int{2, 3, 5, 7}, 11))
}
