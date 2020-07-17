package main

import (
	"fmt"
	"math"
)

/**
 * 动态规划
 */
func LongestIncreasingSubsequence(A []int) int {
	n := len(A)
	f := make([]int, n)
	var max int
	for i := 0; i < n; i++ {
		f[i] = 1
		for j := 0; j < i; j++ {
			if A[j] < A[i] {
				if f[i] < f[j]+1 {
					f[i] = f[j] + 1
				}
			}
		}
		if max < f[i] {
			max = f[i]
		}
	}
	return max
}

/**
 *优化方法
 */
func LongestIncreasingSubsequence2(A []int) int {
	length := len(A)
	b := make([]int, length+1)
	start, end := 0, 0
	b[0] = math.MinInt32
	top := 0
	for i := 0; i < length; i++ {
		j := 0
		end = top
		// 二分法
		for start <= end {
			mid := (start + end) / 2
			if b[mid] < A[i] {
				j = mid
				start = mid + 1
			} else {
				end = mid - 1
			}
		}
		b[j+1] = A[i]
		if j+1 > top {
			top = j + 1
		}
	}
	return top
}

func main() {
	fmt.Println(LongestIncreasingSubsequence2([]int{4, 2, 4, 5, 3, 7}))
}
