package main

import "fmt"

func strStr(haystack string, needle string) int {
	str1 := []byte(haystack)
	str2 := []byte(needle)
	length := len(str2)

	if length < 1 {
		return 0
	}
	if len(str1) < 1 {
		return -1
	}
	for i := 0; i < len(str1)+1-length; i++ {
		temp := 0
		for j := 0; j < length; j++ {
			if str1[i+j] == str2[j] {
				temp++
			}
		}
		if temp == length {
			return i
		}
	}
	return -1
}
func main() {
	fmt.Println(strStr("a", "a"))
}
