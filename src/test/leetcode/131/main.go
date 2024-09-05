package main

import "fmt"

func partition(s string) [][]string {
	table := palindromeTable(s)
	var ret [][]string
	var dfs func(int)
	var sp []string
	n := len(s)
	dfs = func(i int) {
		if i == n {
			ret = append(ret, append([]string{}, sp...))
		}

		for j := i; j < n; j++ {
			if table[i][j] {
				sp = append(sp, s[i:j+1])
				dfs(j + 1)
				sp = sp[:len(sp)-1]
			}
		}
	}
	dfs(0)
	return ret
}

func palindromeTable(s string) [][]bool {
	n := len(s)
	ret := make([][]bool, n)
	for i := 0; i < n; i++ {
		ret[i] = make([]bool, n)
		for j := 0; j < n; j++ {
			ret[i][j] = true
		}
	}
	for i := n - 1; i >= 0; i-- {
		for j := i + 1; j < n; j++ {
			ret[i][j] = s[i] == s[j] && ret[i+1][j-1]
		}
	}
	return ret
}

func palindromeTable2(s string) [][]bool {
	n := len(s)
	ret := make([][]bool, n)
	for i := 0; i < n; i++ {
		ret[i] = make([]bool, n)
		for j := 0; j < n; j++ {
			ret[i][j] = true
		}
	}
	for j := 1; j < n; j++ {
		for i := j - 1; i >= 0; i-- {
			ret[i][j] = s[j] == s[i] && ret[i+1][j-1]
		}
	}
	return ret
}

func main() {
	fmt.Println(palindromeTable("abccba"))
	fmt.Println(palindromeTable2("abccba"))
}
