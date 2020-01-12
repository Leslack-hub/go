package main

import "fmt"

func longestPalindrome(s string) string {
	if len(s) < 2 { // 肯定是回文，直接返回
		return s
	}
	var begin, maxLen int
	for step := 1; step <= len(s); step++ {
		for i := 0; i <= len(s)-step; i++ {
			isSuccess := isSuccess(s[i : i+step])
			if isSuccess {
				begin = i
				maxLen = step
				break
			}
		}
	}
	return s[begin : begin+maxLen]
}

func isSuccess(s string) bool {
	s2 := s
	for i := 0; i < len(s)/2; i++ {
		if s[i] != s2[len(s)-1-i] {
			return false
		}
	}
	return true
}
func main() {
	//fmt.Println(isSuccess("aaxbaa"))
	fmt.Println(longestPalindrome("121-"))
}
