package main

import (
	"fmt"
	"math/rand"
	"time"
)

func generator() chan int {
	c := make(chan int)
	go func() {
		i := 0
		for {
			// time.Sleep(time.Second)
			time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond)
			c <- i
			i++
		}
	}()
	return c
}

func createWorker(id int) chan<- int {
	worker := make(chan int)
	go doWork(id, worker)
	return worker
}

func doWork(id int, w chan int) {
	for n := range w {
		fmt.Printf("id：%d，recevied:%d\n", id, n)
	}
}

func main() {
	c1, c2 := generator(), generator()
	w := createWorker(0)

	var values []int
	tm := time.After(10 * time.Second)
	tmTick := time.Tick(time.Second)
	for {
		activeWorker := make(chan<- int)
		var activeValue int
		if len(values) > 0 {
			activeWorker = w
			activeValue = values[0]
		}
		select {
		case n := <-c1:
			values = append(values, n)
		case n1 := <-c2:
			values = append(values, n1)
		case activeWorker <- activeValue:
			values = values[1:]
		case <-tmTick:
			fmt.Println("lenth:", len(values))
		case <-time.After(800 * time.Millisecond):
			fmt.Println("timeout")
		case <-tm:
			fmt.Println("bye")
			return
		}
	}
}
