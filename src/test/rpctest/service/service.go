package main

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"rpctest"
)

func main() {
	rpc.Register(rpctest.DemoService{})
	listen, err := net.Listen("tcp", ":20001")
	if err != nil {
		panic(listen)
	}

	for {
		accept, err := listen.Accept()
		if err != nil {
			log.Printf("accept err %v", err)
		}

		go jsonrpc.ServeConn(accept)
	}
}
