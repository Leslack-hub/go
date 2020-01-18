package main

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

func reverseLinkList(head *ListNode) *ListNode {
	res := &ListNode{}
	temp := head
	for temp != nil {
		i := temp
		temp = temp.Next
		i.Next = res
		res = i
	}

	return res
}

func main() {
	fmt.Println(reverseLinkList(&ListNode{Val: 3, Next: &ListNode{Val: 2, Next: &ListNode{Val: 1, Next: nil}}}))
}
