package main

import (
	"fmt"
	"math"
	"sort"
)

func nextPermutation(nums []int) {
	p := len(nums)
	isModify := true
	for i := p - 1; i >= 0; i-- {
		// 找比nums[i] 大于的数
		matchKey := findNum(nums[i+1:], nums[i])
		fmt.Println(matchKey)
		if matchKey != -1 && isModify {
			nums[i], nums[i+1+matchKey] = nums[i+1+matchKey], nums[i]
			sort.Ints(nums[i+1:])
			isModify = false
		}
	}
	fmt.Println(nums)
	if isModify {
		i := 0
		j := p - 1
		for i < p/2 {
			nums[i], nums[j] = nums[j], nums[i]
			i++
			j--
		}
	}

}

func findNum(nums []int, i int) int {
	temp := math.MaxInt32
	res := -1
	for k, v := range nums {
		if v > i &&
			v-i < temp &&
			v-i != 0 {
			temp = v - i
			res = k
		}
	}

	return res
}

func main() {
	nextPermutation([]int{1, 3, 2})
}
