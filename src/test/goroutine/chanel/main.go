package main

import (
	"fmt"
	"sync"
)

type worker struct {
	in   chan int
	done func()
}

func createWorker(id int, wg *sync.WaitGroup) worker {
	worker := worker{
		in: make(chan int),
		done: func() {
			wg.Done()
		},
	}
	go doWork(id, worker)
	return worker
}

func doWork(id int, w worker) {
	for n := range w.in {
		fmt.Printf("id：%d，recevied:%c\n", id, n)
		w.done()
	}
}

// func bufferedChanel() {
// 	// 缓冲区，如果chan 没有收数据的程序回deadlock, 加入缓冲区后 超过了 才需要收数据
// 	c := make(chan int, 3)
// 	go worker(0, c)
// 	c <- 'a'
// 	c <- 'b'
// 	c <- 'c'
// 	c <- 'd'
// 	time.Sleep(time.Millisecond)
// }

func chanDemo() {
	var workers [10]worker
	var wg sync.WaitGroup
	wg.Add(20)
	for i := 0; i < 10; i++ {
		workers[i] = createWorker(i, &wg)
	}

	for i, worker := range workers {
		worker.in <- 'a' + i
	}

	for i, worker := range workers {
		worker.in <- 'A' + i
	}

	wg.Wait()
	// wait for all workers 问题:  第一组channel 的done 没有go func收数据 又送了新的数据，会出现deadlock的现象，
	// 解决方式1: 开一个go routine 送done 解决方式2: 分别接受数据 解决方式3: 使用wait group
	// for _, worker := range workers {
	// 	<-worker.done
	// 	<-worker.done
	// }
}

// func channelClose() {
// 	c := make(chan int, 3)
// 	go worker(0, c)
// 	c <- 'a'
// 	c <- 'b'
// 	c <- 'c'
// 	c <- 'd'
// 	close(c)
// 	time.Sleep(time.Millisecond)
// }

func main() {
	chanDemo()
	// bufferedChanel()
	// channelClose()
}
