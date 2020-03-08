package main

import (
	"fmt"
	"leslack/src/helper"
	"math"
)

func perfectSquares(n int) int {
	f := make([]int, n+1)
	f[0] = 0
	for i := 1; i <= n; i++ {
		// i = 2; j=1; f[2] = f[2-1]+1 = f[1]+1=2; j =2
		f[i] = math.MaxInt32
		for j := 1; j*j <= i; j++ {
			f[i] = helper.Min(f[i], f[i-j*j]+1)
		}
	}
	return f[n]
}

func main() {
	fmt.Println(perfectSquares(13))
}
