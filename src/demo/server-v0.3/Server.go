package main

import (
	"fmt"
	"tpc-server/src/lex/ziface"
	"tpc-server/src/lex/znet"
)

// ping test - customized router
type PingRouter struct {
	znet.BaseRouter
}

// test pre-handle
func (this *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("call router prehandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("call back before ping error")
	}
}

// test handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("call router handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping ping ping\n"))
	if err != nil {
		fmt.Println("call back ping error")
	}
}

// test post-handle
func (this *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("call router posthandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
	if err != nil {
		fmt.Println("call back after ping error")
	}
}

func main() {
	s := znet.NewServer("[server v0.3]")
	s.AddRouter(&PingRouter{})
	s.Serve()
}
