package main

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func main() {
	n := &ListNode{Val: 1, Next: &ListNode{Val: 2, Next: &ListNode{Val: 3, Next: &ListNode{Val: 4, Next: &ListNode{Val: 5, Next: &ListNode{Val: 6, Next: &ListNode{Val: 7, Next: nil}}}}}}}
	mid := getMid(n, nil)
	fmt.Println(mid.Val)
}

func getMid(n, r *ListNode) *ListNode {
	slow, fast := n, n
	for fast != r && fast.Next != r {
		slow = slow.Next
		fast = fast.Next.Next
	}
	return slow
}
