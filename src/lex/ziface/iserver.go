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
}
