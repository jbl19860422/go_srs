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
	"go_srs/srs/codec/flv"
	"fmt"
	"go_srs/srs/protocol/rtmp"
	"net/http"
)

type SrsHttpFlvConsumer struct {
	source          *SrsSource
	queue           *SrsMessageQueue
	StreamId		int
	writer			http.ResponseWriter
	flvEncoder		*flvcodec.SrsFlvEncoder
}

func NewSrsHttpFlvConsumer(s *SrsSource, w http.ResponseWriter, r *http.Request) *SrsHttpFlvConsumer {
	return &SrsHttpFlvConsumer{
		source:s,
		writer:w,
		queue:NewSrsMessageQueue(),
		StreamId:0,
		flvEncoder:flvcodec.NewSrsFlvEncoder(w),
	}
}

func (this *SrsHttpFlvConsumer) OnPublish() error {
	return nil
}

func (this *SrsHttpFlvConsumer) OnUnpublish() error {
	return nil
}

func (this *SrsHttpFlvConsumer) ConsumeCycle() error {
	this.flvEncoder.WriteHeader()
	go func() {
		notify := this.writer.(http.CloseNotifier).CloseNotify()
		<- notify
		this.StopConsume()
	}()
	this.writer.Header().Set("Content-Type", "video/x-flv")
	for {
		msg, err := this.queue.Wait()
		if err != nil {
			return err
		}

		if msg != nil {
			if msg.GetHeader().IsVideo() {
				if flvcodec.VideoIsKeyFrame(msg.GetPayload()) {
					fmt.Println("send key frame")
				}
				this.flvEncoder.WriteVideo(uint32(msg.GetHeader().GetTimestamp()), msg.GetPayload())
			} else if msg.GetHeader().IsAudio() {
				this.flvEncoder.WriteAudio(uint32(msg.GetHeader().GetTimestamp()), msg.GetPayload())
			} else {
				this.flvEncoder.WriteMetaData(msg.GetPayload())
			}
		}
	}
}

func (this *SrsHttpFlvConsumer) StopConsume() error {
	this.source.RemoveConsumer(this)
	//send connection close to response writer
	this.queue.Break()
	return nil
}

func (this *SrsHttpFlvConsumer) OnRecvError(err error) {
	this.StopConsume()
}

func (this *SrsHttpFlvConsumer) Enqueue(msg *rtmp.SrsRtmpMessage, atc bool, jitterAlgorithm *SrsRtmpJitterAlgorithm) {
	this.queue.Enqueue(msg)
}