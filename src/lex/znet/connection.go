package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"tcp-server/src/lex/ziface"
)

type Connection struct {
	Conn     *net.TCPConn
	ConnID   uint32
	isClosed bool
	ExitChan chan bool
	Router   ziface.IRouter
}

func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}
	return c
}

func (c *Connection) startReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID = ", c.ConnID, " Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// create pack obj in order to unpack request's msg
		dp := NewDataPack()
		// read message header
		headerData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headerData); err != nil {
			fmt.Println("read msg header error", err)
			break
		}
		_msg, err := dp.Unpack(headerData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}
		var data []byte
		if _msg.GetDataLen() > 0 {
			data = make([]byte, _msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error", err)
				break
			}
		}
		_msg.SetData(data)

		// get request
		req := Request{
			conn: c,
			msg:  _msg,
		}

		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID = ", c.ConnID)
	// Separate the read and write
	go c.startReader()
	// TODO:: Start Writer
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop().. ConnID = ", c.ConnID)

	if c.isClosed {
		return
	}
	c.isClosed = true

	c.Conn.Close()
	close(c.ExitChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// pack data before sending to clients
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed when send msg")
	}
	// pack the data before sending (Len/ID/Data)
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NetMsgPack(msgId, data))
	if err != nil {
		fmt.Println("pack error msg id = ", msgId)
		return errors.New("pack msg error")
	}
	// send binary msg to clients
	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("Write msg id ", msgId, " error :", err)
		return errors.New("conn write error")
	}
	return nil
}
