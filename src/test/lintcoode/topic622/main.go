package main

import (
	"fmt"
	"strconv"
)

/**
题意：
有一条小河上有N个石头，位置一次在a0<a1<... a(n-1)
有一只青蛙在第一个石头上
青蛙一开始可以向右跳距离为1
它必须一直向右调，并且落在石头上
如果上次跳的距离是L,这次跳的距离可以是L-1,L或者L+1
问能否到达最后一个石头
子问题：
	设f[i][j] 表示是否能最后一跳j跳到石头a(i)
*/
func FrogJump(A []int) bool {
	n := len(A)
	f := make(map[int]int, n)
	for i := 0; i < n; i++ {
		f[A[i]] = i
	}
	return false
}

// 1921681233

func isIPFour(s string) int {
	n := len(s)
	f := make([][]int, 5)
	for i := 0; i <= 4; i++ {
		f[i] = make([]int, n+1)
	}
	f[0][0] = 1
	for i := 1; i <= 4; i++ {
		for j := i; j <= n; j++ {
			for x := 1; x <= 3; x++ {
				if j-x >= 0 {
					str := s[j-x : j]
					if str == "0" {
						f[i][j] += f[i-1][j-x]
						continue
					}
					if str[0] == '0' {
						continue
					}
					varI, _ := strconv.Atoi(str)
					if varI <= 255 {
						f[i][j] += f[i-1][j-x]
					}
				}
			}
		}
	}
	fmt.Println(f)
	return 0
}

func dfs(level int, index int, str string, res []string) {
	if level == 5 || index == len(str)-1 {
		if level == 5 && index == len(str)-1 {
			fmt.Println(res)
		}
		return
	}
	for i := 1; i <= 3; i++ {
		if index+i+1 > len(str) {
			continue
		}
		s := str[index+1 : index+i+1]
		//if s != "0" && s[index+1] != '0' {
		if s != "0" && s[0] != '0' {
			varI, _ := strconv.Atoi(s)
			if varI <= 255 {
				res = append(res, s)
				dfs(level+1, index+i, str, res)
				res = res[:level-1]
			}
		}
	}
}

func main() {
	//fmt.Println(FrogJump([]int{0, 1, 3, 5, 6, 8, 12, 17}))
	//fmt.Println(isIPFour("19216801"))
	res := []string{}
	dfs(1, -1, "19216801", res)

}
