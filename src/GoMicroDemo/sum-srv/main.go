package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/util/log"
	"leslack/src/GoMicroDemo/proto/sum"
	"leslack/src/GoMicroDemo/sum-srv/handler"
)

func main() {
	src := micro.NewService(
		micro.Name("go.micro.learning.srv.sum"),
	)
	src.Init(
		micro.BeforeStart(func() error {
			log.Log("启动前的日志")
			return nil
		}),
		micro.AfterStart(func() error {
			log.Log("启动后的日志")
			return nil
		}),
		)
	sum.RegisterSumHandler(src.Server(),handler.Handler())
	err := src.Run()
	if err != nil {
		panic(err)
	}
}
