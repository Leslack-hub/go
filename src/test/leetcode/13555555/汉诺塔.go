package main

import "fmt"

func main() {
	hanoi("A", "B", "C", 2)
}

func hanoi(A string, B string, C string, n int) {
	if n > 0 {
		hanoi(A, C, B, n-1)
		fmt.Println(A, "->", C)
		hanoi(B, A, C, n-1)
	}
}
