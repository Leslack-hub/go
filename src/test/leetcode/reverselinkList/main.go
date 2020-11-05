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
		// 当前值
		i := temp
		// 移动指针
		temp = temp.Next
		// 下一个指针指向结果
		i.Next = res
		// 结果就是当前i
		res = i
	}

	return res
}

func main() {
	fmt.Println(reverseLinkList(&ListNode{Val: 3, Next: &ListNode{Val: 2, Next: &ListNode{Val: 1, Next: nil}}}))
}
