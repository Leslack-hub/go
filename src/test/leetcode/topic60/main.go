package main

import (
	"fmt"
)

func getPermutation(n int, k int) string {
	if n == 0 {
		return ""
	}
	res := make([]byte, n)
	var rec []byte
	for i := 0; i < n; i++ {
		rec = append(rec, byte(i)+'1')
	}
	k--
	base := 1
	for i := 2; i < n; i++ {
		base *= i
	}
	for i := 0; i < n-1; i++ {
		idx := k / base
		res[i] = rec[idx]
		rec = append(rec[:idx], rec[idx+1:]...)
		k %= base
		base /= n - i - 1
	}
	res[n-1] = rec[0]

	return string(res)
}

func main() {
	fmt.Println(getPermutation(4, 23))
}
