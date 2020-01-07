package main

import (
	"json_rpc"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main() {
	rpc.Register(json_rpc.DemoService{})
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
