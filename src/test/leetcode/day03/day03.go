package main

import (
	"fmt"
)

func lengthOfLongestSubstring(s string) int {
	list := make(map[uint8]int)
	// 相当于 每次有新的值进入 计算出现在的长度，如果有相同的值之前出现过，
	// 将之前的这个值之前出现的丢掉，如果值大于已经存入的长度，则替换，否者截取数组的长度
	max, left := 0, 0

	for i := 0; i < len(s); i++ {
		if list[s[i]] >= left {
			left = list[s[i]] + 1
		} else if i+1-left > max {
			max = i + 1 - left
		}
		list[s[i]] = i
	}

	return max
}

func main() {
	fmt.Println(lengthOfLongestSubstring("abcabc"))
}
