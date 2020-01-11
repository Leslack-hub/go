package main

import (
	"fmt"
	"runtime"
	"time"
)

func chanDemo() {
	c := make(chan int)
	c <- 1
	c <- 2
	n := <-c
	fmt.Println(n)
}

func main() {
	var a [10]int
	for i := 0; i < 10; i++ {
		go func(i int) {
			for {
				a[i]++
				runtime.Gosched()
			}
		}(i)
	}
	time.Sleep(time.Millisecond)
	fmt.Println(a)
}
