package main

import "fmt"

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func pathSum(root *TreeNode, targetSum int) [][]int {
	var path [][]int
	var dfs func([]int, *TreeNode, int)
	dfs = func(cur []int, root *TreeNode, sum int) {
		if root == nil {
			return
		}
		cur = append(cur, root.Val)
		defer func() { cur = cur[:len(cur)-1] }()

		if root.Left == nil && root.Right == nil {
			if root.Val == sum {
				path = append(path, append([]int{}, cur...))
			}
			return
		}

		s := sum - root.Val
		dfs(cur, root.Left, s)
		dfs(cur, root.Right, s)
	}

	dfs([]int{}, root, targetSum)
	return path
}

var retln [][]int

func check(cur []int, root *TreeNode, s int) {
	if root == nil {
		return
	}

	if root.Left == nil && root.Right == nil {
		if root.Val == s {
			retln = append(retln, append(cur, root.Val))
		}
		return
	}

	cur = append(cur, root.Val)
	s = s - root.Val
	check(cur, root.Left, s)
	check(cur, root.Right, s)
}

func main() {
	//fmt.Println(pathSum(&TreeNode{Val: 1, Left: &TreeNode{Val: 2}, Right: &TreeNode{Val: 3}}, 3))
	fmt.Println(dp1("babgbag", "bag"))
}

func dp1(s, t string) int {
	s1 := len(s)
	s2 := len(t)

	dp := make([][]int, s2+1)
	for i := 0; i <= s2; i++ {
		dp[i] = make([]int, s1+1)
	}
	for i := 0; i <= s1; i++ {
		dp[0][i] = 1
	}
	for i := 1; i <= s2; i++ {
		for j := 1; j <= s1; j++ {
			if s[j-1] == t[i-1] {
				dp[i][j] = dp[i-1][j-1] + dp[i][j-1]
			} else {
				dp[i][j] = dp[i][j-1]
			}
		}
	}
	return dp[s2][s1]
}
