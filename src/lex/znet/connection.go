package znet

import (
	"fmt"
	"net"
	"tcp-server/src/lex/utils"
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

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID = ", c.ConnID, " Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err", err)
			continue
		}

		// get request
		req := Request{
			conn: c,
			data: buf,
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
	go c.StartReader()
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

func (c *Connection) Send(data []byte) error {
	return nil
}