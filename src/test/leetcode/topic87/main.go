package main

import "fmt"

func isScramble(s1 string, s2 string) bool {
	m, n := len(s1), len(s2)
	if m != n {
		return false
	}
	f := make([][][]bool, n)
	for i := range f {
		f[i] = make([][]bool, n)
		for j := range f[i] {
			f[i][j] = make([]bool, n+1)
		}
	}
	// 初始化单个字符的情况
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			f[i][j][1] = s1[i] == s2[j] // 字符变字符
		}
	}
	// 枚举区间长度 2~n
	for k := 2; k <= n; k++ {
		// 枚举 S 中的起点位置
		for i := 0; i <= n-k; i++ {
			// 枚举 T 中的起点位置
			for j := 0; j <= n-k; j++ {
				f[i][j][k] = false
				// 枚举划分位置
				for w := 1; w <= k-1; w++ {
					// 第一种情况: S1->T1, S2->T2
					if f[i][j][w] && f[i+w][j+w][k-w] {
						f[i][j][k] = true
						break
					}
				}
				for w := 1; w <= k-1; w++ {
					// 第二种情况: S1 -> T2, S2 -> T1
					// S1 起点 i，T2 起点 j + 前面那段长度 len-k ，S2 起点 i + 前面长度k
					if f[i][j+k-w][w] && f[i+w][j][k-w] {
						f[i][j][k] = true
						break
					}
				}
			}
		}
	}
	return f[0][0][n]
}

func main() {
	fmt.Println(isScramble("great", "rgeat"))
}
