package main

import "fmt"

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func bfs(root *TreeNode, res *[]int, num int) {
	if root == nil {
		return
	}
	var queue []TreeNode
	queue = append(queue, *root)
	for len(queue) > 0 {
		var temp []TreeNode
		for i := 0; i < len(queue); i++ {
			*res = append(*res, queue[i].Val)
			if queue[i].Left != nil {
				temp = append(temp, *queue[i].Left)
			}
			if queue[i].Left != nil {
				temp = append(temp, *queue[i].Right)
			}
		}
		queue = temp
	}
	fmt.Println(res)
}

func main() {
	var res []int
	bfs(&TreeNode{
		Val: 6,
		Left: &TreeNode{
			Val: 2,
			Left: &TreeNode{
				Val:   0,
				Left:  nil,
				Right: nil,
			},
			Right: &TreeNode{
				Val: 4,
				Left: &TreeNode{
					Val:   3,
					Left:  nil,
					Right: nil,
				},
				Right: &TreeNode{
					Val:   5,
					Left:  nil,
					Right: nil,
				},
			},
		},
		Right: &TreeNode{
			Val: 8,
			Left: &TreeNode{
				Val:   7,
				Left:  nil,
				Right: nil,
			},
			Right: &TreeNode{
				Val:   9,
				Left:  nil,
				Right: nil,
			},
		},
	}, &res, 0)
}
