package main

import (
	"fmt"
	"leslack/src/helper"
	"math"
)

func palindromePartitioningII(s string) int {
	length := len(s)
	f := make([]int, length+1)
	f[0] = 0
	for i := 1; i <= length; i++ {
		f[i] = math.MaxInt32
		for j := 0; j < i; j++ {
			r := isRoll(s[j:i])
			if r {
				f[i] = helper.Min(f[i], f[j]+1)
			}
		}
	}
	fmt.Println(f)
	return 0
}

func isRoll(s string) bool {
	r := true
	x, y := 0, len(s)-1
	for x < y {
		if s[x] == s[y] {
			x++
			y--
		} else {
			r = false
			break
		}
	}
	return r
}

func main() {
	fmt.Println(palindromePartitioningII("acabccbaac"))
}
