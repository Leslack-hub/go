package main

import "fmt"

/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func levelOrder(root *TreeNode) [][]int {
	res := make([][]int, 20)
	for i := 0; i < 20; i++ {
		res[i] = make([]int, 0)
	}
	helper(root, 0, res)

	return res
}

func helper(head *TreeNode, level int, res [][]int) {
	if head == nil {
		return
	}

	res[level] = append(res[level], head.Val)
	helper(head.Left, level+1, res)
	helper(head.Right, level+1, res)
	fmt.Println(res)
}

func main() {
	levelOrder(&TreeNode{
		Val:   3,
		Left:  &TreeNode{
			Val:   9,
			Left:  nil,
			Right: nil,
		},
		Right: &TreeNode{
			Val:   20,
			Left:  &TreeNode{
				Val:   15,
				Left:  nil,
				Right: nil,
			},
			Right: &TreeNode{
				Val:   7,
				Left:  nil,
				Right: nil,
			},
		},
	})
}
