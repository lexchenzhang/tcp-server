package main

import "tcp-server/src/lex/znet"

func main() {
	s := znet.NewServer("mmo game server")

	s.Serve()
}
