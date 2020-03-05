package main
import "fmt"
func main() {
	var balance = [4]int{0, 0, 0, 0}
	queen(balance, 0)
}
func queen(a [4]int, cur int) {
	if cur == len(a) {
		fmt.Print(a)
		fmt.Println()
		return
	}
	for i := 0; i < len(a); i++ {
		a[cur] = i
		flag := true
		for j := 0; j < cur; j++ {
			// i表示当前的数字 a[j] 表示前面已经出现的数字
			// 0 2 i = 1  cur = 3; i = 1; j=0 时 ab = 1, a[j] =0  temp == 3 - 0; 进入下一次递归
			// 0 2 1 0 cur =4
			ab := i - a[j]
			temp := 0
			if ab > 0 {
				temp = ab
			} else {
				temp = -ab
			}
			// 如果前面已经出现了 i 直接跳过，
			if a[j] == i || temp == cur-j {
				flag = false
				break
			}
		}
		if flag {
			queen(a, cur+1)
		}
	}
}