package main

import (
	"fmt"
	"leslack/src/helper"
)

func houseRobber(A []int) int {
	length := len(A)
	if length == 1 {
		return A[0]
	}
	//old := 0
	//new := A[0]
	//for i := 2; i <= length; i++ {
	//	t := helper.Max(new, old+A[i-1])
	//	old, new = new, t
	//}
	f := make([]int, length+1)
	f[0] = 0
	f[1] = A[0]
	for i := 2; i <= length; i++ {
		f[i] = helper.Max(f[i-1], f[i-2]+A[i-1])
	}
	return f[length]
}

func main() {
	fmt.Println(houseRobber([]int{3, 8, 4}))
}
