package main

import "fmt"

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func isSymmetric(root *TreeNode) bool {
	check2(root)
	fmt.Println(check3(root))
	fmt.Println(check4(root))
	fmt.Println(maxDepth(root))
	return check(root, root)
}

func check(p, q *TreeNode) bool {
	if p == nil && q == nil {
		return true
	}
	if p == nil || q == nil {
		return false
	}

	return p.Val == q.Val && check(p.Left, q.Right) && check(p.Right, q.Left)
}

func check2(p *TreeNode) {
	var q []*TreeNode
	q = append(q, p)
	for len(q) > 0 {
		t := q[0]
		fmt.Println(t.Val)
		if t.Left != nil {
			q = append(q, t.Left)
		}
		if t.Right != nil {
			q = append(q, t.Right)
		}
		q = q[1:]
	}
}

func check3(root *TreeNode) [][]int {
	var q []*TreeNode
	var result [][]int
	q = append(q, root)
	for i := 0; len(q) > 0; i++ {
		result = append(result, []int{})
		var p []*TreeNode
		for j := 0; j < len(q); j++ {
			node := q[j]
			result[i] = append(result[i], node.Val)
			if node.Left != nil {
				p = append(p, node.Left)
			}
			if node.Right != nil {
				p = append(p, node.Right)
			}
		}
		q = p
	}
	return result
}

func check4(root *TreeNode) [][]int {
	var q []*TreeNode
	var result [][]int
	q = append(q, root)
	for i := 0; len(q) > 0; i++ {
		result = append(result, []int{})
		var p []*TreeNode
		for j := 0; j < len(q); j++ {
			node := q[j]
			if i%2 == 1 {
				result[i] = append([]int{node.Val}, result[i]...)
			} else {
				result[i] = append(result[i], node.Val)
			}
			if node.Left != nil {
				p = append(p, node.Left)
			}
			if node.Right != nil {
				p = append(p, node.Right)
			}
		}
		q = p
	}
	return result
}

func maxDepth(root *TreeNode) int {
	return check5(root, 0)
}

func check5(root *TreeNode, deep int) int {
	if root == nil {
		return deep
	}
	deep++
	return max(check5(root.Left, deep), check5(root.Right, deep))
}

func main() {

	fmt.Println(isSymmetric(&TreeNode{
		Val: 1,
		Left: &TreeNode{
			Val: 2,
			Left: &TreeNode{
				Val: 4,
			},
			Right: &TreeNode{
				Val: 5,
				Left: &TreeNode{
					Val: 8,
				},
				Right: &TreeNode{
					Val: 9,
				},
			},
		},
		Right: &TreeNode{
			Val: 3,
			Left: &TreeNode{
				Val: 6,
			},
			Right: &TreeNode{
				Val: 7,
				Left: &TreeNode{
					Val: 10,
				},
			},
		},
	}))
}
