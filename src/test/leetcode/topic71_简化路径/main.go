package main

import (
	"fmt"
	"strings"
)

func simplifyPath(path string) string {
	buf := strings.Split(path, "/")
	var strack []string
	for i := 0; i < len(buf); i++ {
		if buf[i] == "" || buf[i] == "." {
			continue
		}
		if buf[i] == ".." {
			strack = strack[0 : len(strack)-1]
		} else {
			strack = append(strack, buf[i])
		}
	}
	return "/" + strings.Join(strack, "/")
}

func main() {
	fmt.Println(simplifyPath("/a/../../b/../c//.//"))
}
