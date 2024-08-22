package main

import "fmt"

const limit = 6

var result [][limit]int

func main() {
	var balance = [limit]int{}
	queen(balance, 0)
	fmt.Println(result)
	fmt.Println(len(result))
}
func queen(ret [limit]int, level int) {
	if level == len(ret) {
		result = append(result, ret)
		return
	}
	// 第几排 开始尝试数字
	for number := 0; number < len(ret); number++ {
		ret[level] = number
		flag := true
		// 判断第一排的位置 到 现在的位置 是否冲突
		for i := 0; i < level; i++ {
			tmp := number - ret[i]
			var differ int
			if tmp > 0 {
				differ = tmp
			} else {
				differ = -tmp
			}
			// 如果前面已经出现了 number 直接跳过， differ == level-i 表示在同一个对角线了
			if ret[i] == number || differ == level-i {
				flag = false
				break
			}
		}
		if flag {
			// 表示进入下一行
			queen(ret, level+1)
		}
	}
}
