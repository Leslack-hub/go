package main

import (
	"leslack/src/crawler/engine"
	"leslack/src/crawler/perslset"
	"leslack/src/crawler/scheduler"
	"leslack/src/crawler/zhenai/parser"
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
