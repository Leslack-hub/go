package main

import "fmt"

func main() {
	fmt.Println(findMedianSortedArrays([]int{4,2,3}, []int{2}))
}

func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
	nums := combine(nums1, nums2)
	fmt.Println(nums)
	return medianOf(nums)
}
// mis [1,2] njs [3,4]
func combine(mis, njs []int) []int {
	// 左数组长度
	lenMis, i := len(mis), 0
	// 右边数组的长度
	lenNjs, j := len(njs), 0
	res := make([]int, lenMis+lenNjs)

	for k := 0; k < lenMis+lenNjs; k++ {
		// 第一次循环 false
		if i == lenMis ||
			(i < lenMis && j < lenNjs && mis[i] > njs[j]) {
			res[k] = njs[j]
			j++
			continue
		}

		// 第一次循环
		if j == lenNjs ||
			(i < lenMis && j < lenNjs && mis[i] <= njs[j]) {
			res[k] = mis[i]
			i++
		}
	}

	return res
}

func medianOf(nums []int) float64 {
	l := len(nums)

	if l == 0 {
		panic("切片的长度为0，无法求解中位数。")
	}

	if l%2 == 0 {
		return float64(nums[l/2]+nums[l/2-1]) / 2.0
	}

	return float64(nums[l/2])
}