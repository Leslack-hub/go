package main

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

func deleteDuplicates(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	flag := false
	for head.Next != nil && head.Val == head.Next.Val {
		head, flag = head.Next, true
	}
	head.Next = deleteDuplicates(head.Next)
	if flag {
		return head.Next
	}
	return head
}

func main() {
	result := deleteDuplicates(&ListNode{1, &ListNode{2, &ListNode{3, &ListNode{3, &ListNode{4, &ListNode{4, &ListNode{5, nil}}}}}}})
	for result != nil {
		fmt.Println(result)
		result = result.Next
	}
}
