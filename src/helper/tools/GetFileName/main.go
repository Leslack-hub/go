package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	files := os.Args[1]
	//paths, fileName := filepath.Split(files)
	//fmt.Println(paths, fileName)      //获取路径中的目录及文件名 E:\data\  test.txt
	base := filepath.Base(files)
	// 正则 测试
	re := regexp.MustCompile(`(.*)\.`)
	match := re.FindStringSubmatch(base)
	fmt.Println(match[1])
	//for _, m := range match {
	//	fmt.Println(m)
	//}
}
