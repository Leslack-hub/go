package main

import "fmt"

func totalNQueens(n int) int {
	var res int
	balance := make([]int, n)
	queen(balance, 0, &res)
	return res
}

func queen(a []int, cur int, res *int) {
	if cur == len(a) {
		*res++
		return
	}
	for i := 0; i < len(a); i++ {
		a[cur] = i
		flag := true
		for j := 0; j < cur; j++ {
			ab := i - a[j]
			temp := 0
			if ab > 0 {
				temp = ab
			} else {
				temp = -ab
			}
			if a[j] == i || temp == cur-j {
				flag = false
				break
			}
		}
		if flag {
			queen(a, cur+1, res)
		}
	}
}

func main() {
	fmt.Println(totalNQueens(4))
}
