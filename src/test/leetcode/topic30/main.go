package main

import "fmt"

func findSubstring(s string, words []string) []int {
	pattern := sortStringNum(words)
	var res []int
	for _, v := range pattern {
		length := len(v)
		for i := 0; i < len(s)+1-length; i++ {
			temp := 0
			for j := 0; j < length; j++ {
				if s[i+j] == v[j] {
					temp++
				}
			}
			if temp == length {
				if !inArray(res, i) {
					res = append(res, i)
				}
			}
		}
	}
	return res
}

func sortStringNum(strings []string) []string {
	if len(strings) < 2 {
		return strings
	}
	if len(strings) == 2 {
		return []string{strings[0] + strings[1], strings[1] + strings[0]}
	}
	var res, pattern []string
	strList := make([]string, len(strings))
	for i := 0; i < len(strList); i++ {
		copy(strList, strings)
		pattern = sortStringNum(append(strList[:i], strList[i+1:]...))
		for _, j := range pattern {
			res = append(res, strings[i]+j)
		}
	}
	return res
}

func inArray(nums []int, v int) bool {
	if len(nums) == 0 {
		return false
	}
	for _, i := range nums {
		if v == i {
			return true
		}
	}
	return false
}

func findSubstring2(s string, words []string) []int {
	lens := len(s)
	res := make([]int, 0, lens)

	lenws := len(words)
	if lens == 0 || lenws == 0 || lens < lenws*len(words[0]) {
		return res
	}

	lenw := len(words[0])

	record := make(map[string]int, lenws)
	for _, w := range words {
		if len(w) != lenw {
			return res
		}
		record[w]++
	}

	remain := make(map[string]int, lenws)
	count := 1
	left, right := 0, 0

	/**
	 * s[left:right] 作为一个窗口存在
	 * 假设 word:= s[right:right+lenw]
	 * 如果 remain[word]>0 ，那么移动窗口 右边，同时更新移动后 s[left:right] 的统计信息
	 * 		remain[word]--
	 * 		right += lenw
	 * 		count++
	 * 否则，移动窗口 左边，同时更新移动后 s[left:right] 的统计信息
	 * 		remain[s[left:left+lenw]]++
	 * 		count--
	 * 		left += lenw
	 *
	 * 每次移动窗口 右边 后，如果 count = lenws ，那么
	 * 		说明找到了一个符合条件的解
	 * 		append(res, left)，然后
	 * 		移动窗口 左边
	 */

	// reset 重置 remain 和 count
	reset := func() {
		if count == 0 {
			// 统计记录没有被修改，不需要重置
			// 因为有这个预判，所以需要第一次执行 reset 时
			// count 的值不能为 0
			// 即 count 的初始值不能为 0
			return
		}
		for k, v := range record {
			remain[k] = v
		}
		count = 0
	}

	// moveLeft 让 left 指向下一个单词
	moveLeft := func() {
		// left 指向下一个单词前，需要修改统计记录
		// 增加 left 指向的 word 可出现次数一次，
		remain[s[left:left+lenw]]++
		// 连续符合条件的单词数减少一个
		count--
		// left 后移一个 word 的长度
		left += lenw
	}

	// left 需要分别从{0,1,...,lenw-1}这些位置开始检查，才能不遗漏
	for i := 0; i < lenw; i++ {
		left, right = i, i
		reset()

		// s[left:] 的长度 >= words 中所有 word 组成的字符串的长度时，
		// s[left:] 中才有可能存在要找的字符串
		for lens-left >= lenws*lenw {
			word := s[right : right+lenw]
			remainTimes, ok := remain[word]

			switch {
			case !ok:
				// word 不在 words 中
				// 从 right+lenw 处，作为新窗口，重新开始统计
				left, right = right+lenw, right+lenw
				reset()
			case remainTimes == 0:
				// word 的出现次数上次就用完了
				// 说明 word 在 s[left:right] 中出现次数过多
				moveLeft()
				// 这个case会连续出现
				// 直到 s[left:right] 中的统计结果是 remain[word] == 1
				// 这个过程中，right 一直不变
			default:
				// ok && remainTimes > 0，word 符合出现的条件
				// moveRight
				remain[word]--
				count++
				right += lenw
				// 检查 words 能否排列组合成 s[left:right]
				if count == lenws {
					res = append(res, left)
					// moveLeft 可以避免重复统计 s[left+lenw:right] 中的信息
					moveLeft()
				}
			}
		}
	}

	return res
}
func main() {
	fmt.Println(findSubstring2("barfoothefoobarman", []string{"foo","bar"}))
}
