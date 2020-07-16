package main

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

func partition(head *ListNode, x int) *ListNode {
	beforeHead := &ListNode{}
	afterhead := &ListNode{}
	before := beforeHead
	after := afterhead
	for head != nil {
		if head.Val < x {
			before.Next = head
			before = before.Next
		} else {
			after.Next = head
			after = after.Next
		}
		head = head.Next
	}
	after.Next = nil
	before.Next = afterhead.Next
	return beforeHead.Next
}

func main() {
	result := partition(&ListNode{1, &ListNode{4, &ListNode{3, &ListNode{2, &ListNode{5, &ListNode{2, nil}}}}}}, 3)
	for result != nil {
		fmt.Println(result)
		result = result.Next
	}
}
