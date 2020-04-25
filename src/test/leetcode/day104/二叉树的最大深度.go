package day104

import "leslack/src/helper"

/**
* Definition for a binary tree node.
 */
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func maxDepth(root *TreeNode) int {
	return dfs(root, 0)
}

func dfs(root *TreeNode, i int) int {
	if root == nil {
		return i
	}
	return helper.Max(dfs(root.Left, i+1), dfs(root.Right, i+1))
}
