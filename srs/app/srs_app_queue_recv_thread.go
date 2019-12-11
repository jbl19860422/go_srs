/*
The MIT License (MIT)

Copyright (c) 2019 GOSRS(gosrs)

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package app

import (
	"go_srs/srs/protocol/rtmp"
)

type SrsQueueRecvThread struct {
	queue      []*rtmp.SrsRtmpMessage
	consumer   *SrsConsumer
	rtmp       *rtmp.SrsRtmpServer
	recvThread *SrsRecvThread
}

func NewSrsQueueRecvThread(c *SrsConsumer, s *rtmp.SrsRtmpServer) *SrsQueueRecvThread {
	st := &SrsQueueRecvThread{
		queue:    make([]*rtmp.SrsRtmpMessage, 1000),
		consumer: c,
		rtmp:     s,
	}

	st.recvThread = NewSrsRecvThread(s, st, 1000)
	return st
}

func (this *SrsQueueRecvThread) Start() {
	this.recvThread.Start()
}

func (this *SrsQueueRecvThread) Stop() {
	this.recvThread.Stop()
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
	this.consumer.OnRecvError(err)
	return
}

func (this *SrsQueueRecvThread) OnThreadStart() {
	return
}

func (this *SrsQueueRecvThread) OnThreadStop() {
	return
}
