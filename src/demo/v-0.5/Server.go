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

// test handle
func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("call router handle...")
	fmt.Println("recv from client: msgID=", request.GetMsgID(), " data=", string(request.GetData()))
	err := request.GetConnection().SendMsg(1, []byte("ping..."))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	s := znet.NewServer("[server v0.5]")
	s.AddRouter(&PingRouter{})
	s.Serve()
}
