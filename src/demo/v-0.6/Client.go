package main

import (
	"fmt"
	"io"
	"net"
	"tcp-server/src/lex/znet"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		return
	}
	var idx int
	for {
		idx++
		dp := znet.NewDataPack()
		binMsg, err := dp.Pack(znet.NetMsgPack(uint32(idx%2), []byte("Hi Server")))
		if err != nil {
			return
		}
		if _, err := conn.Write(binMsg); err != nil {
			fmt.Println("client write error", err)
			return
		}
		// after sending msg the server should send back a ping msg
		binHeader := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binHeader); err != nil {
			break
		}
		if msgHeader, err := dp.Unpack(binHeader); err != nil {
			fmt.Println("client unpack msg header err", err)
			break
		} else if msgHeader.GetDataLen() > 0 {
			msg := msgHeader.(*znet.Message)
			msg.Data = make([]byte, msg.DataLen)
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error", err)
				return
			}
			fmt.Println("-> Recv Server Msg: ID=", msg.GetMsgId(), " Len=", msg.DataLen, " Data=", string(msg.Data))
		}

		time.Sleep(1 * time.Second)
	}
}
