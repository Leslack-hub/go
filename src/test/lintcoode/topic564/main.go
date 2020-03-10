package main

import "fmt"

/**
 * 背包问题
 */
func backpackVI(A []int, m int) int {
	n := len(A)
	f := make([]int, m+1)
	f[0] = 1
	// i= 3; j = 0; f[3] += f[3-1];
	for i := 1; i <= m; i++ {
		for j := 0; j < n; j++ {
			if i >= A[j] {
				f[i] += f[i-A[j]]
			}
		}
	}
	fmt.Println(f)
	return f[m]
}

func main() {
	fmt.Println(backpackVI([]int{1, 2, 4}, 4))
}
