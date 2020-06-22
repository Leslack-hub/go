package main

import "fmt"

func search(nums []int, target int) bool {
	start, length := 0, len(nums)-1
	if length == 0 {
		return false
	}
	var mid int
	end := length
	for start <= end {
		mid = start + (end-start)/2
		if nums[mid] == target {
			return true
		}
		if nums[start] == nums[mid] {
			start++
			continue
		}

		if nums[start] < nums[mid] {
			// target在前面
			if nums[mid] > target && nums[start] <= target {
				end = mid - 1
			} else {
				start = mid + 1
			}
		} else {
			if nums[mid] < target && nums[end] >= target {
				start = mid + 1
			} else {
				end = mid - 1
			}
		}
	}

	return false
}

func main() {
	fmt.Println(search([]int{2, 5, 6, 0, 0, 1, 2}, 2))
}
