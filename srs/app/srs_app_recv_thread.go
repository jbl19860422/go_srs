package app

import (
	"fmt"
	"go_srs/srs/protocol/rtmp"
)

type ISrsMessageHandler interface {
	Handle(msg *rtmp.SrsRtmpMessage) error
	OnRecvError(err error)
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
		done:    make(chan bool),
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
			// fmt.Printf("msg[0]=%x, msg[1]=%x, msg[2]=%x, msg[3]=%x", msg.GetPayload()[0], msg.GetPayload()[1], msg.GetPayload()[2], msg.GetPayload()[3])
			err = this.handler.Handle(msg)
		}

		if err != nil {
			this.handler.OnRecvError(err)
			close(this.done)
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
	close(this.done)
	return nil
}

func (this *SrsRecvThread) Stop() {
	close(this.exit) //直接关闭，避免cycle先退出
}

func (this *SrsRecvThread) Join() {
	<-this.done
}
