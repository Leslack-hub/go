package main

import "fmt"

func search(nums []int, target int) int {
	rotated := indexOfMin(nums)
	size := len(nums)
	left, right := 0, size-1
	for left <= right {
		mid := (left + right) / 2
		rotatedMid := (mid + rotated) % size
		switch {
		case nums[rotatedMid] < target:
			left = mid + 1
		case nums[rotatedMid] > target:
			right = mid - 1
		default:
			return rotatedMid
		}

	}
	return -1
}

/* nums 是被旋转了的递增数组 */
func indexOfMin(nums []int) int {
	size := len(nums)
	left, right := 0, size-1
	for left < right {
		mid := (left + right) / 2
		if nums[right] < nums[mid] {
			left = mid + 1
		} else {
			right = mid
		}
	}
	return left
}

func main() {
	fmt.Println(search([]int{4,5,6,7,0,1,2}, 0))
}
