package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"tcp-server/src/lex/utils"
	"tcp-server/src/lex/ziface"
)

type Connection struct {
	TcpServer ziface.IServer
	Conn      *net.TCPConn
	ConnID    uint32
	isClosed  bool
	ExitChan  chan bool
	// communication between read-goroutin and write-oroutin
	msgChan      chan []byte
	MsgHandler   ziface.IMsgHandler
	property     map[string]interface{}
	propertyLock sync.RWMutex
}

func NewConnection(s ziface.IServer, conn *net.TCPConn, connID uint32, msghandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:  s,
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msghandler,
		isClosed:   false,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
		property:   make(map[string]interface{}),
	}
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

func (c *Connection) startReader() {
	fmt.Println("[Reader Goroutine is running]")
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

		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandler.DoMsgHandler(&req)
		}

	}
}

func (c *Connection) startWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), " [conn Writer exit.]")
	for {
		select {
		case data := <-c.msgChan:
			// data in chan for client
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error, ", err)
				return
			}
		case <-c.ExitChan:
			// reader exited
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID = ", c.ConnID)
	// Separate the read and write
	go c.startReader()
	go c.startWriter()
	// call framework user's hook func OnConnStart after new connection
	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop().. ConnID = ", c.ConnID)

	if c.isClosed {
		return
	}
	c.isClosed = true
	c.ExitChan <- true
	// call framework user's hook func OnConnStart before disconnect
	c.TcpServer.CallOnConnStop(c)
	c.Conn.Close()
	c.TcpServer.GetConnMgr().Remove(c)
	close(c.ExitChan)
	close(c.msgChan)
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
	// send binary msg to channel
	c.msgChan <- binaryMsg
	return nil
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}
