package main

import "fmt"

func removeElement(nums []int, val int) int {
	p := 0
	for i := 0; i < len(nums); i++ {
		if nums[i] != val {
			nums[p] = nums[i]
			p++
		}
	}
	fmt.Println(nums)

	return p
}

func main() {
	fmt.Println(removeElement([]int{0, 1, 2, 2, 3, 0, 4, 2}, 2))
}
