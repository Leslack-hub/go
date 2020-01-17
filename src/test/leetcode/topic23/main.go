package main

import (
	"fmt"
	"sort"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

func mergeKLists(lists []*ListNode) *ListNode {
	res := &ListNode{}
	var length int
	var ValList []int
	var headList []*ListNode
	for _, i := range lists {
		temp := i
		for temp != nil {
			length++
			ValList = append(ValList, temp.Val)
			headList = append(headList, temp)
			temp = temp.Next
		}
	}
	sort.Ints(ValList)
	first := res
	for i := 0; i < length; i++ {
		for j := 0; j < len(headList); j++ {
			if ValList[i] == headList[j].Val {
				first.Next = headList[j]
				first = first.Next
				headList = append(headList[:j], headList[j+1:]...)
				break
			}
		}
	}

	return res.Next
}

func main() {
	fmt.Println(mergeKLists([]*ListNode{
		{Val: 3, Next: &ListNode{Val: 2, Next: nil}},
		{Val: 1, Next: &ListNode{Val: 4, Next: nil}},
	}))
}
