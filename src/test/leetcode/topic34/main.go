package main

import "fmt"

func searchRange(nums []int, target int) []int {
	// 寻找左边界
	left := searchLeft(nums, target)
	right := searchRight(nums, target)
	return []int{left, right}
}

func searchLeft(nums []int, target int) int {
	size := len(nums)
	if size == 0 {
		return -1
	}
	left, right := 0, size-1
	// 找到左边边界
	for left <= right {
		mid := (left + right) / 2
		if nums[mid] == target {
			right = mid - 1
		} else if nums[mid] < target {
			left = mid + 1
		} else if nums[mid] > target {
			right = mid - 1
		}
	}
	if left <= size-1 &&
		nums[left] == target {
		return left
	} else {
		return -1
	}
}

func searchRight(nums []int, target int) int {
	size := len(nums)
	if size == 0 {
		return -1
	}
	left, right := 0, size-1
	for left <= right {
		mid := (left + right) / 2
		if nums[mid] == target {
			left = mid + 1
		} else if nums[mid] < target {
			left = mid + 1
		} else if nums[mid] > target {
			right = mid - 1
		}
	}
	if right >= 0 &&
		nums[right] == target {
		return right
	} else {
		return -1
	}
}

func main() {
	//fmt.Println(searchLeft([]int{5, 7, 7, 8, 8, 10}, 8))
	//fmt.Println(searchRight([]int{5, 7, 7, 8, 8, 8, 8, 9, 10}, 8))
	fmt.Println(searchRange([]int{5, 7, 7, 8, 8, 8, 8, 8, 10}, 8))
}
