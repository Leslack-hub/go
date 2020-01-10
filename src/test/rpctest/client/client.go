package main

import (
	"fmt"
	"net"
	"net/rpc/jsonrpc"
	"rpctest"
)

func main() {
	conn, err := net.Dial("tcp", ":20001")
	if err != nil {
		panic(err)
	}

	client := jsonrpc.NewClient(conn)
	var result *float64
	dealFunc := func(err error) {
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(*result)
		}
	}

	err = client.Call("DemoService.Div", rpctest.Args{A:10, B:5}, &result)
	dealFunc(err)

	err = client.Call("DemoService.Div", rpctest.Args{10, 0}, &result)
	dealFunc(err)
}
