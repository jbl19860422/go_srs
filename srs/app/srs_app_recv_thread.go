package app

import (
	"fmt"
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
	rtmp    *rtmp.SrsRtmpServer
	handler ISrsMessageHandler
	timeout int32
	exit    chan bool
	done    chan bool
}

func NewSrsRecvThread(r *rtmp.SrsRtmpServer, h ISrsMessageHandler, timeoutMS int32) *SrsRecvThread {
	return &SrsRecvThread{
		rtmp:    r,
		handler: h,
		timeout: timeoutMS,
		exit:    make(chan bool),
		done:	 make(chan bool),
	}
}

func (this *SrsRecvThread) Start() {
	go this.cycle()
}

func (this *SrsRecvThread) cycle() error {
DONE:
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
		case <-this.exit:
			{
				fmt.Println("********************888quit***************")
				break DONE
			}
		default:
			{
				//continue
			}
		}
	}
	this.done <- true
	return nil
}

func (this *SrsRecvThread) Stop() error {
	this.exit <- true
	// <- this.done
	return nil
}
