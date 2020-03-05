package helper

/**
 * 求最大值
 */
func Max(nums ...int) int {
	max := 0
	for i := 0; i < len(nums); i++ {
		if nums[i] > max {
			max = nums[i]
		}
	}

	return max
}

/**
 * 最小值
 */
func Min(nums ...int) int {
	length := len(nums)
	min := 0
	if length < 1 {
		return min
	}
	min = nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i] < min {
			min = nums[i]
		}
	}

	return min
}
