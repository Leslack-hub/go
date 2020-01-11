package day02

import (
	"fmt"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	res := &ListNode{}
	cur := res
	carry := 0

	for l1 != nil || l2 != nil || carry > 0 {
		sum := carry
		if l1 != nil {
			sum += l1.Val
			l1 = l1.Next
		}
		if l2 != nil {
			sum += l2.Val
			l2 = l2.Next
		}

		carry = sum / 10
		cur.Next = &ListNode{Val: sum % 10}
		cur = cur.Next
	}

	return res.Next
}

func main() {
	node1 := &ListNode{Val: 3, Next: nil}
	node2 := &ListNode{Val: 4, Next: node1}
	node3 := &ListNode{Val: 2, Next: node2}
	numbers := addTwoNumbers(node3, &ListNode{Val: 5, Next: &ListNode{Val: 6, Next: &ListNode{Val: 4, Next: nil}}})
	for numbers != nil {
		fmt.Println(numbers.Val)
		numbers = numbers.Next
	}
}
