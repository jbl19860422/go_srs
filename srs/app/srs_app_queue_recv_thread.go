package app

import (
	"go_srs/srs/protocol/rtmp"
)

type SrsQueueRecvThread struct {
	queue 		[]*rtmp.SrsRtmpMessage
	consumer 	*SrsConsumer
	rtmp		*rtmp.SrsRtmpServer
	recvThread 	*SrsRecvThread
}

func NewSrsQueueRecvThread(c *SrsConsumer, s *rtmp.SrsRtmpServer) *SrsQueueRecvThread {
	st := &SrsQueueRecvThread{
		queue:make([]*rtmp.SrsRtmpMessage, 1000),
		consumer:c,
	}

	st.recvThread = NewSrsRecvThread(s, st, 1000)
	return st
}

func (this *SrsQueueRecvThread) Start() {
	this.recvThread.Start()
}

func (this *SrsQueueRecvThread) CanHandle() bool {
	return true
}

func (this *SrsQueueRecvThread) Handle(msg *rtmp.SrsRtmpMessage) error {

	//todo fix cid change
	//todo nbmsg++
	this.queue = append(this.queue, msg)
	return nil
}

func (this *SrsQueueRecvThread) Size() int {
	return len(this.queue)
}

func (this *SrsQueueRecvThread) Empty() bool {
	return len(this.queue) == 0
}

func (this *SrsQueueRecvThread) GetMsg() *rtmp.SrsRtmpMessage {
	if this.Empty() {
		return nil
	}

	m := this.queue[0]
	this.queue = this.queue[1:]
	return m
}

func (this *SrsQueueRecvThread) OnRecvError(err error) {
	return
}	

func (this *SrsQueueRecvThread) OnThreadStart() {
	return
}

func (this *SrsQueueRecvThread) OnThreadStop() {
	return
}

