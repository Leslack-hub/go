package main

import "fmt"

func coinsInLine(n int) bool {
	f := make([]bool, n+1)
	f[0] = false
	f[1], f[2] = true, true
	for i := 3; i <= n; i++ {
		if f[i-1] == false ||
			f[i-2] == false {
			f[i] = true
		} else {
			f[i] = false
		}
	}
	return f[n]
}

func main() {
	fmt.Println(coinsInLine(3))
}
