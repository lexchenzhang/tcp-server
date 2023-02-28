package ziface

type IServer interface {
	// start
	Start()
	// stop
	Stop()
	// run
	Serve()
}
