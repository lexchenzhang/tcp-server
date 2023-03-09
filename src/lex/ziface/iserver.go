package ziface

type IServer interface {
	// start
	Start()
	// stop
	Stop()
	// run
	Serve()
	// router
	AddRouter(msdID uint32, router IRouter)
	// get conn manager
	GetConnMgr() IConnManager
	// register OnConnStart
	SetOnConnStart(func(conn IConnection))
	// register OnConnStop
	SetOnConnStop(func(conn IConnection))
	// trigger OnConnStart
	CallOnConnStart(conn IConnection)
	// trigger OnConnStop
	CallOnConnStop(conn IConnection)
}
