package main

import (
	"github.com/leSlcak/go/crawlor/engine"
	"github.com/leSlcak/go/crawlor/perslset"
	"github.com/leSlcak/go/crawlor/scheduler"
	"github.com/leSlcak/go/crawlor/zhenai/parser"
)

func main() {
	e := engine.ConcurrentEngine{
		Scheduler:   &scheduler.SimpleScheduler{},
		WorkerCount: 10,
		ItemChan:    perslset.ItemSaver(),
	}

	e.Run(engine.Request{
		Url:        "http://www.zhenai.com/zhenghun/",
		ParserFunc: parser.ParserCityList,
	})
}
