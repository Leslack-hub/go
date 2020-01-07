package main

import (
	"fmt"
	"regexp"
)

const text = `my name is haha@lishuaishuai.com 
fddsafds@fdsaf.com.cn
`

func main() {
	// 正则 测试
	re := regexp.MustCompile(`([a-zA-Z0-9]+)@([a-zA-Z0-9]+)(\.[a-zA-Z0-9.]+)`)
	match := re.FindAllStringSubmatch(text, -1)
	for _, m := range match {
		fmt.Println(m)
	}
}
