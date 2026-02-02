package main

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

func deleteDuplicates(head *ListNode) *ListNode {
	current := head
	for current != nil && current.Next != nil {
		if current.Val == current.Next.Val {
			current.Next = current.Next.Next
		} else {
			current = current.Next
		}
	}
	return head
}

func main() {
	// result := deleteDuplicates(&ListNode{1, &ListNode{2, &ListNode{3, &ListNode{3, &ListNode{4, &ListNode{4, &ListNode{5, nil}}}}}}})
	result := deleteDuplicates(&ListNode{1, &ListNode{1, &ListNode{1, nil}}})
	for result != nil {
		fmt.Println(result)
		result = result.Next
	}
}
