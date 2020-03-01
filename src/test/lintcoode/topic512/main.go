package main

import "fmt"

func numDecoding(ss string) int {
	n := len(ss)
	if n == 0 {
		return 1
	}
	f := make([]int, n+1)
	f[0] = 1
	for i := 1; i <= n; i++ {
		if ss[i-1] >= '1' && ss[i-1] <= '9' {
			f[i] += f[i-1]
		}
		if i > 1 {
			j := (ss[i-2]-'0')*10 + (ss[i-1] - '0')
			if ss[i-2] != '0' && j <= 26 && j >= 10 {
				f[i] += f[i-2]
			}
		}
	}
	return f[n]
}

func main() {
	fmt.Println(decodeWay("12"))
}
