package main

import "fmt"

func mergeSort(arr []int, start, end int) {
	if start >= end {
		return
	}
	mid := (start + end) / 2
	mergeSort(arr, start, mid)
	mergeSort(arr, mid+1, end)
	merge(arr, start, mid, end)
}

func merge(arr []int, start, mid, end int) {
	var tmparr = []int{}
	var s1, s2 = start, mid + 1
	for s1 <= mid && s2 <= end {
		if arr[s1] > arr[s2] {
			tmparr = append(tmparr, arr[s2])
			s2++
		} else {
			tmparr = append(tmparr, arr[s1])
			s1++
		}
	}
	if s1 <= mid {
		tmparr = append(tmparr, arr[s1:mid+1]...)
	}
	if s2 <= end {
		tmparr = append(tmparr, arr[s2:end+1]...)
	}
	for pos, item := range tmparr {
		arr[start+pos] = item
	}
}

var a = []int{3, 4, 0, 1, 5, 6, 7, 8}

func main() {
	mergeSort(a, 0, len(a)-1)
	fmt.Println(a)
}
