package main

import "fmt"

func canCompleteCircuit(gas []int, cost []int) int {
	n := len(gas)
	for i := 0; i < n; {
		sumG, sumC, cur := 0, 0, 0
		for cur < n {
			j := (i + cur) % n
			sumG += gas[j]
			sumC += cost[j]
			if sumC > sumG {
				break
			}
			cur++
		}
		if cur == n {
			return i
		} else {
			i += cur + 1
		}
	}
	return -1
}

func main() {
	fmt.Println(canCompleteCircuit([]int{1, 2, 3, 4, 5}, []int{3, 4, 5, 1, 2}))
}
