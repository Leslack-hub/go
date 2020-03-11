package main

import (
	"fmt"
	"leslack/src/helper"
	"math"
)

func copyBooks(K int, books []int) int {
	n := len(books)
	f := make([][]int, K+1)
	for i := 0; i <= K; i++ {
		f[i] = make([]int, n+1)
	}
	f[0][0] = 0
	for i := 1; i <= n; i++ {
		f[0][i] = math.MaxInt32
	}

	// f[1][1]
	for k := 1; k <= K; k++ {
		f[k][0] = 0
		for i := 1; i <= n; i++ {
			f[k][i] = math.MaxInt32
			max := 0
			for j := 0; j < i; j++ {
				f[k][i] = helper.Min(f[k][i], helper.Max(f[k-1][j], max))
				max += books[j]
				fmt.Println("max1", max)
			}

			//for j := i; j >= 0; j-- {
			//	fmt.Println("max2", max)
			//	f[k][i] = helper.Min(f[k][i], helper.Max(f[k-1][j], max))
			//	if j > 0 {
			//		max += books[j-1]
			//	}
			//}
		}
	}
	fmt.Println(f)
	return f[K][n]
}

func main() {
	fmt.Println(copyBooks(2, []int{3, 2, 4}))
}
