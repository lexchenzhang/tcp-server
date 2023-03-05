package main

import (
	"fmt"
	"tcp-server/src/lex/ziface"
	"tcp-server/src/lex/znet"
)

// ping test - customized router
type PingRouter struct {
	znet.BaseRouter
}

// test pre-handle
func (pr *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("call router prehandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("call back before ping error")
	}
}

// test handle
func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("call router handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping ping ping\n"))
	if err != nil {
		fmt.Println("call back ping error")
	}
}

// test post-handle
func (pr *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("call router posthandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
	if err != nil {
		fmt.Println("call back after ping error")
	}
}

func main() {
	s := znet.NewServer("[server v0.4]")
	s.AddRouter(&PingRouter{})
	s.Serve()
}
