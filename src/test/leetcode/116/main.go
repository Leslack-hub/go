package main

import (
	"fmt"
)

type Node struct {
	Val   int
	Left  *Node
	Right *Node
	Next  *Node
}

func connect(root *Node) *Node {
	if root == nil {
		return nil
	}
	queue := []*Node{root}
	for len(queue) > 0 {
		var tmp []*Node
		for i := 0; i < len(queue); i++ {
			if i < len(queue)-1 {
				queue[i].Next = queue[i+1]
			}
			if queue[i].Left != nil {
				tmp = append(tmp, queue[i].Left)
			}
			if queue[i].Right != nil {
				tmp = append(tmp, queue[i].Right)
			}
		}
		queue = tmp
	}
	return root
}

func main() {
	f := []int{1, 2, 3, 4}
	fmt.Println(min(0, f...))
}
