package main

import "fmt"

func subsets(nums []int) [][]int {
	output := make([][]int, 0)
	output = append(output, []int{})
	for _, v := range nums {
		size := len(output)
		for i := 0; i < size; i++ {
			temp := output[i]
			temp = append(temp, v)
			output = append(output, temp)
		}
	}
	return output
}

func main() {
	fmt.Println(subsets([]int{1, 2, 3}))
}
