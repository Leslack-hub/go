package main

type ListNode struct {
	Val  int
	Next *ListNode
}

func removeNthFromEnd(head *ListNode, n int) *ListNode {
	var length int
	res := head
	temp := head
	var list []ListNode
	for {
		list = append(list, *temp)
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
	for length > 0 {
		length--

	}

	return res

}
func main() {
	// 1->2->3->4->5
	removeNthFromEnd(&ListNode{Val: 1, Next: &ListNode{Val: 2, Next: nil}}, 1)
}
