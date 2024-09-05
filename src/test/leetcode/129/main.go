package main

import (
	"fmt"
	"strconv"
)

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func sumNumbers(root *TreeNode) int {
	return helper(root, "")
}

func helper(root *TreeNode, s string) int {
	if root == nil {
		return 0
	}
	s = fmt.Sprintf("%s%d", s, root.Val)

	if root.Left == nil && root.Right == nil {
		num, _ := strconv.Atoi(s)
		return num
	}
	return helper(root.Left, s) + helper(root.Right, s)
}

func main() {
	fmt.Println(sumNumbers(&TreeNode{1, &TreeNode{2, nil, nil}, &TreeNode{3, nil, nil}}))
}
