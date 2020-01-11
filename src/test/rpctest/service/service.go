package main

import (
	"leslack/src/test/rpctest"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

// sendMessage: {"method": "DemoService.Div", "params":[{"A":3,"B":4}],"id":1}
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
