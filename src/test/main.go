package main

import "fmt"

const text = `my name is haha@lishuaishuai.com 
fddsafds@fdsaf.com.cn
`

func main() {
	strings := []string{"hello.com"}
	fmt.Println(strings)
	//str := []int{1,2,3,4}
	//for _, i := range str {
	//	fmt.Println(i)
	//}

	//re := regexp.MustCompile(`([a-zA-Z0-9]+)@([a-zA-Z0-9]+)(\.[a-zA-Z0-9.]+)`)
	//match := re.FindAllStringSubmatch(text, -1)
	//for _, m := range match {
	//	fmt.Println(m)
	//}
}
