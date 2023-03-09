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

// hi test - customized router
type HiRouter struct {
	znet.BaseRouter
}

// ping router
func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("call router handle...")
	fmt.Println("recv from client: msgID=", request.GetMsgID(), " data=", string(request.GetData()))
	err := request.GetConnection().SendMsg(0, []byte("ping..."))
	if err != nil {
		fmt.Println(err)
	}
}

// hello router
func (hr *HiRouter) Handle(request ziface.IRequest) {
	err := request.GetConnection().SendMsg(1, []byte("hi..."))
	fmt.Println("recv from client: msgID=", request.GetMsgID(), " data=", string(request.GetData()))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	s := znet.NewServer("[server v0.9]")
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HiRouter{})
	s.Serve()
}
