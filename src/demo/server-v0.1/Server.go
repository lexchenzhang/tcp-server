package main

import "tpc-server/src/lex/znet"

func main() {
	s := znet.NewServer("[server v0.1]")
	s.Serve()
}
