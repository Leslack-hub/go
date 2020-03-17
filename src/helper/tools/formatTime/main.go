package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	var t = os.Args[1]
	timestamp, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		fmt.Println("格式化失败")
	}
	s := time.Unix(timestamp, 0).String()
	fmt.Println(s)
}
