package app

import(
	"go_srs/srs/protocol/rtmp"
)
type ISrsMessageHandler interface {
	CanHandle() bool
	Handle(msg *rtmp.SrsRtmpMessage) error
	OnRecvError(err error)
	OnThreadStart()
	OnThreadStop()
}

type SrsRecvThread struct {
	rtmp 		*rtmp.SrsRtmpServer
	handler 	ISrsMessageHandler
	timeout		int32
	exit		chan bool
}

func NewSrsRecvThread(r *rtmp.SrsRtmpServer, h ISrsMessageHandler, timeoutMS int32) *SrsRecvThread {
	return &SrsRecvThread{
		rtmp:r,
		handler:h,
		timeout:timeoutMS,
		exit:make(chan bool),
	}
}

func (this *SrsRecvThread) Start() {
	go this.cycle()
}

func (this *SrsRecvThread) cycle() error {
	for {
		msg, err := this.rtmp.RecvMessage()
		if err == nil {
			err = this.handler.Handle(msg)
		}

		if err != nil {
			this.handler.OnRecvError(err)
			return err
		}

		select {
		case <-this.exit:{
			break
		}
		default:{
			//continue
		}
		}
	}
}

func (this *SrsRecvThread) Stop() error {
	this.exit <- true
	return nil
}

