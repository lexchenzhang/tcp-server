package znet

import (
	"fmt"
	"net"
	"tcp-server/src/lex/utils"
	"tcp-server/src/lex/ziface"
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
	// msg handler
	MsgHandler ziface.IMsgHandler
}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name : %s, listenner at IP : %s, Port:%d is starting\n", utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version %s\n", utils.GlobalObject.Version)
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

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Router added with msgID=", msgID)
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandler(),
	}
	return s
}

func clientHanlder(s *Server) {
	// 0.start task queue and worker pool
	s.MsgHandler.StartWorkerPool()
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
	fmt.Println("start TCP server success")
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
		connHandler := NewConnection(conn, cid, s.MsgHandler)
		cid++

		go connHandler.Start()
	}
}
