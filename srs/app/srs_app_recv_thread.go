/*
The MIT License (MIT)

Copyright (c) 2013-2015 GOSRS(gosrs)

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
