package main

import (
	"fmt"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

func swapPairs(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	head.Next, head.Next.Next, head = head.Next.Next, head, head.Next
	head.Next.Next = swapPairs(head.Next.Next)
	return head
}

func main() {
	fmt.Println(swapPairs(&ListNode{Val: 2, Next: &ListNode{Val: 1, Next: &ListNode{Val: 4, Next: &ListNode{Val: 3, Next: nil}}}}))
}
