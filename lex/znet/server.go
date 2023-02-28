package znet

import (
	"errors"
	"fmt"
	"net"
	"tpc-server/lex/ziface"
)

type Server struct {
	// name
	Name string
	// ipv4
	IPVersion string
	// ip
	IP string
	// port
	Port int
}

func (s *Server) Start() {
	fmt.Printf("[Start] server Listening at IP :%s, Port %d \n", s.IP, s.Port)
	go clientHanlder(s)
}

func (s *Server) Stop() {
	// TODO :: release or GC
}

func (s *Server) Serve() {
	s.Start()
	// TODO :: extra services
	select {}
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "127.0.0.1",
		Port:      8999,
	}
	return s
}

func clientHanlder(s *Server) {
	// 1.gain TCP addr
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println(err)
		return
	}
	// 2.listen
	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("start Zinx server success")
	var cid uint32
	cid = 0
	// 3.clients hanlder
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("Accept err", err)
			continue
		}
		var remoteAddr = conn.RemoteAddr()
		fmt.Println("Accept ", remoteAddr)
		// connect to client
		// echo back (512)
		connHandler := NewConnection(conn, cid, func(c *net.TCPConn, data []byte, cnt int) error {
			fmt.Println("Conn Handler...")
			if _, err := c.Write(data[:cnt]); err != nil {
				fmt.Println("write back buf err", err)
				return errors.New("callback to client error")
			}
			return nil
		})
		cid++

		go connHandler.Start()
	}
}
