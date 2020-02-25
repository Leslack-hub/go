package main

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

func reverseKGroup2(head *ListNode, k int) *ListNode {
	temp := head
	var length int
	for temp != nil {
		temp = temp.Next
		length++
	}
	if length < k {
		return head
	}

	res := &ListNode{}
	first := head
	for i := 0; i < k; i++ {
		temp := first
		first = first.Next
		temp.Next = res.Next
		res.Next = temp
	}
	secondLink := res
	for i := 0; i < k; i++ {
		secondLink = secondLink.Next
	}
	secondLink.Next = reverseKGroup2(first, k)

	return res.Next
}
func reverseKGroup(head *ListNode, k int) *ListNode {
	if k == 1 {
		return head
	}
	dummy := &ListNode{}
	temp := head
	var length int
	for temp != nil {
		temp = temp.Next
		length++
	}
	if length < k {
		return nil
	}

	first := head
	sortLink := dummy
	for i := 0; i < length; i += k {
		temp := &ListNode{}
		temp2 := temp
		for j := 0; j < k; j++ {
			temp2.Next = first
			first = first.Next
			temp2 = temp2.Next
		}
		temp2.Next = nil
		sortLink.Next = reverseLink(temp.Next)
		for x := 0; x < k; x++ {
			sortLink = sortLink.Next
		}
	}

	positive := head
	n := length % k
	length -= n
	for i := 0; n != 0 && i < length; i++ {
		positive = positive.Next
	}
	sortLink.Next = positive
	return dummy.Next
}

func reverseLink(head *ListNode) *ListNode {
	res := &ListNode{}
	temp := head
	for temp != nil {
		i := temp
		temp = temp.Next
		i.Next = res.Next
		res.Next = i
	}

	return res.Next
}

func main() {
	fmt.Println(reverseKGroup2(&ListNode{Val: 1, Next: &ListNode{Val: 2, Next: &ListNode{Val: 3, Next: &ListNode{Val: 4, Next: nil}}}}, 2))
}
