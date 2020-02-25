package main

type ListNode struct {
	Val  int
	Next *ListNode
}

func rotateRight(head *ListNode, k int) *ListNode {
	var length int
	temp := head
	for temp != nil {
		length++
		temp = temp.Next
	}
	if length <= 1 {
		return head
	}
	key := k % length
	// 反转链表
	first := head
	res := &ListNode{}
	for i := 0; i < length; i++ {
		if length-i == key {
			res.Next = first
			break
		}
		first = first.Next
	}
	secondLink := res
	for i := 0; i < key; i++ {
		secondLink = secondLink.Next
	}
	second := head
	for i := 0; i < length-key; i++ {
		secondLink.Next = second
		secondLink = secondLink.Next
		second = second.Next
	}
	secondLink.Next = nil

	return res.Next
}
func main() {
	rotateRight(&ListNode{1, &ListNode{2, &ListNode{3, &ListNode{4, &ListNode{5, nil}}}}}, 2)
}
