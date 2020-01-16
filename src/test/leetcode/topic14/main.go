package main

import (
	"fmt"
)

func longestCommonPrefix(strs []string) string {
	if len(strs) < 1 {
		return ""
	}

	var str, str2 []byte
	for i := 0; i < len(strs[0]); i++ {
		str = append(str, strs[0][i])
	}

	min := len(strs[0])
	for i := 1; i < len(strs); i++ {
		if len(strs[i]) < min {
			min = len(strs[i])
		}
	}

loop:
	for k := 0; k < min; k++ {
		for i := 1; i < len(strs); i++ {
			if strs[i][k] != str[k] {
				break loop
			}
		}
		str2 = append(str2, str[k])
	}

	return string(str2[:])
}

func main() {
	fmt.Println(longestCommonPrefix([]string{"aa", "a"}))
}
