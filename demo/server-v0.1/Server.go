package main

import "tpc-server/lex/znet"

func main() {
	s := znet.NewServer("[zinx v0.1]")
	s.Serve()
}
