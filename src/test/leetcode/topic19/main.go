package main

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

func removeNthFromEnd(head *ListNode, n int) *ListNode {
	var length int
	res := head
	temp := head
	for {
		if temp.Next == nil {
			break
		}
		temp = temp.Next
		length++
	}
	if length < 1 {
		return res
	}

	length -= n
	first := res
	// 找到倒数n 个节点
	for length > 0 {
		length--
		first = first.Next
	}
	// 将它的下个节点设置为下下个节点
	first.Next = first.Next.Next

	return res
}

/**
 * 链表删除节点就是将删除的上一个节点设置为它的下下个节点
 */
func main() {
	// 1->2->3->4->5
	fmt.Println(removeNthFromEnd(&ListNode{Val: 1, Next: &ListNode{Val: 2, Next: nil}}, 1))
}
