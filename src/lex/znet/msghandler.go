package znet

import (
	"fmt"
	"strconv"
	"tcp-server/src/lex/utils"
	"tcp-server/src/lex/ziface"
)

type MsgHandler struct {
	APIs           map[uint32]ziface.IRouter
	TaskQueue      []chan ziface.IRequest
	WorkerPoolSize uint32
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		APIs:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	// get msg id from request
	router, ok := mh.APIs[request.GetMsgID()]
	if !ok {
		fmt.Println("API msgID=", request.GetMsgID(), " is not found! Register needed!")
	}
	// call router according to msg id
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

func (mh *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	if _, ok := mh.APIs[msgID]; ok {
		panic("duplicate api, msgID=" + strconv.Itoa(int(msgID)))
	}
	mh.APIs[msgID] = router
	fmt.Println("Add API MsgID=", msgID, " succ!")
}

// start worker pool only once
func (mh *MsgHandler) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskSize)
		go mh.startSingleWorker(i)
	}
}

func (mh *MsgHandler) startSingleWorker(id int) {
	for {
		select {
		case request := <-mh.TaskQueue[id]:
			fmt.Println("grab form task queue and workerID=", id)
			mh.DoMsgHandler(request)
		}
	}
}

func (mh *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	workderID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	mh.TaskQueue[workderID] <- request
}
