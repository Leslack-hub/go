package main

import (
	"fmt"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

func reverseBetween(head *ListNode, m int, n int) *ListNode {
	dummy, i, j := &ListNode{Next: head}, m, n-m
	d := dummy
	for i > 1 {
		d = d.Next
		i--
	}
	cur := d.Next.Next
	pre := d.Next
	for j > 0 {
		pre.Next = cur.Next
		cur.Next = d.Next
		d.Next = cur
		cur = pre.Next
		j--
	}
	return dummy.Next
}

func main() {
	result := reverseBetween(&ListNode{1,
		&ListNode{2,
			&ListNode{3,
				&ListNode{4,
					&ListNode{5, nil},
				},
			},
		},
	}, 2, 4)
	for result != nil {
		fmt.Println(result)
		result = result.Next
	}
}

func printListNode(node *ListNode) {
	result := node
	for result != nil {
		fmt.Println(result)
		result= result.Next
	}
}