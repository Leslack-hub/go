package main

import "fmt"

var m = map[byte]byte{
	'(': ')',
	'[': ']',
	'{': '}',
}

func isValid(s string) bool {
	if len(s) == 0 {
		return true
	}
	if len(s)%2 != 0 {
		return false
	}

	var stack []byte
	top := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '(', '[', '{':
			stack = append(stack, c)
			top++
		case ')', '}', ']':
			if len(stack) > 0 && m[stack[top-1]] == c {
				stack = stack[:top-1]
				top--
			}
		}
	}

	if len(stack) > 0 {
		return false
	} else {
		return true
	}
}

func main() {
	fmt.Println(isValid("(([]){})"))
}
