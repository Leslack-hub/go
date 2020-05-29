package main

import "fmt"

func climbStairs(n int) int {
	if n == 1 {
		return 1
	}
	f := make([]int, n+1)
	f[0], f[1] = 1, 1
	for i := 3; i <= n; i++ {
		f[i] = f[i-2] + f[i-1]
	}
	return f[n]
}

func main() {
	fmt.Println(climbStairs(10))
}
