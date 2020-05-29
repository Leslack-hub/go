package main

import "fmt"

func plusOne(digits []int) []int {
	length := len(digits)
	for i :=  length -1; i >= 0; i-- {
		num := digits[i] + 1
		if num / 10 == 0 {
			digits[i] = num
			return digits
		} else {
			digits[i] = num % 10 
		}
	}
	temp := []int{1}
	digits = append(temp,digits...)
	return digits
}

func main() {
	fmt.Println(plusOne([]int{9}))
}