package main

import "fmt"

func backpackV(A []int, m int) int {
	n := len(A)
	f := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		f[i] = make([]int, m+1)
	}
	f[0][0] = 1
	for i := 1; i <= n; i++ {
		f[i][0] = 1
		for j := 1; j <= m; j++ {
			f[i][j] = f[i-1][j]
			if j >= A[i-1] {
				f[i][j] += f[i-1][j-A[i-1]]
			}
		}
	}
	fmt.Println(f)
	return f[n][m]
}

func backpackV2(A []int, m int) int {
	n := len(A)
	f := make([]int, m+1)
	f[0] = 1
	for i := 1; i <= n; i++ {
		for j := m; j >= 0; j-- {
			if j >= A[i-1] {
				f[j] += f[j-A[i-1]]
			}
		}
	}
	fmt.Println(f)
	return 0
}
func main() {
	fmt.Println(backpackV([]int{1, 2, 3, 3, 7}, 5))
}
