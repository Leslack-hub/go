package main

import "fmt"

func lengthOfLastWord(s string) int {
	var res int
	end := len(s) - 1
	for end >= 0 && s[end] == ' ' {
		end--
	}
	for end >= 0 {
		if s[end] != ' ' {
			res++
			end--
		} else {
			break
		}
	}

	return res
}

func main() {
	fmt.Println(lengthOfLastWord("a "))
}
