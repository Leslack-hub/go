package main

import (
	"bufio"
	"github.com/atotto/clipboard"
	"os"
)

type Handler interface {
	run()
}

func worker(c chan func()) {
	for n := range c {
		n()
		// 复制内容到剪切板
		clipboard.WriteAll(clipboardList[point])
	}
}

var clipboardList = [...]string{
	"风急天高猿啸哀，渚清沙白鸟飞回",
	"无边落木萧萧下，不尽长江滚滚来",
	"万里悲秋常作客，百年多病独登台",
	"艰难苦恨繁霜鬓，潦倒新停浊酒杯",
}

var point = 0

func main() {
	chanFunc := make(chan func())
	length := len(clipboardList)
	go worker(chanFunc)
	for {
		inputReader := bufio.NewReader(os.Stdin)   //创建一个读取器，并将其与标准输入绑定。
		input, err := inputReader.ReadString('\n') //读取器对象提供一个方法 ReadString(delim byte) ，该方法从输入中读取内容，直到碰到 delim 指定的字符，然后将读取到的内容连同 delim 字符一起放到缓冲区。
		if err != nil {
			panic("程序错误")
		}
		callable := func() {}
		switch input {
		case "n\n":
			callable = func() {
				point++
				point = point % length
			}
		case "p\n":
			callable = func() {
				point--
				if point < 0 {
					point = length - 1
				}
			}
		}
		chanFunc <- callable
	}
}
