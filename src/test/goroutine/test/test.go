package main

import "fmt"

import "time"

func main() {
	fmt.Println("waiting...")
	c1 := make(chan int)
	var v1 []int
	go func() {
		for {
			time.Sleep(time.Second)
			v1 = append(v1, <-c1)
		}
	}()

	go func() {
		for {
			fmt.Println("hello")
			time.Sleep(time.Millisecond)
		}
	}()

	for {
		c1 <- 1
		time.Sleep(time.Second)
		if len(v1) >= 3 {
			fmt.Println("长度：", len(v1))
			break
		}
	}
}
