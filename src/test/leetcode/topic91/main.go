package main

import "fmt"

func numDecodings(s string) int {
	length := len(s)
	if length == 0 || s[0] == '0' {
		return 0
	}
	pre, curr := 1, 1
	for i := 1; i < length; i++ {
		temp := curr
		if s[i] == '0' {
			if s[i-1] == '1' || s[i-1] == '2' {
				curr = pre
			} else {
				return 0
			}
		} else if s[i-1] == '1' || (s[i-1] == '2' && s[i] >= '1' && s[i] <= '6') {
			curr += pre
		}
		pre = temp
	}
	return curr
}

func main() {
	fmt.Println(numDecodings("226"))
}
