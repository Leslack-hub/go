package day104

import (
	"fmt"
	"testing"
)

func Test_dfs(t *testing.T) {
	fmt.Println(maxDepth(&TreeNode{
		Val: 5,
		Left: &TreeNode{
			Val: 2,
			Left: &TreeNode{
				Val:   3,
				Left:  nil,
				Right: nil,
			},
			Right: nil,
		},
		Right: &TreeNode{
			Val:   2,
			Left:  &TreeNode{
				Val:   3,
				Left:  &TreeNode{
					Val:   5,
					Left:  &TreeNode{
						Val:   6,
						Left:  nil,
						Right: nil,
					},
					Right: nil,
				},
				Right: nil,
			},
			Right: nil,
		},
	}))
}
