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
	// "fmt"
	"go_srs/srs/protocol/rtmp"
	"go_srs/srs/codec/flv"
	"go_srs/srs/protocol/packet"
	// "context"
	// "time"
	"errors"
)

type ConsumerStopListener interface {
	OnConsumerStop()
}

type SrsConsumer struct {
	source          *SrsSource
	conn            *SrsRtmpConn
	queue           *SrsMessageQueue
	StreamId		int
	queueRecvThread *SrsQueueRecvThread
}

func NewSrsConsumer(s *SrsSource, c *SrsRtmpConn) *SrsConsumer {
	//todo
	consumer := &SrsConsumer{
		queue:  NewSrsMessageQueue(),
		source: s,
		conn:   c,
		StreamId: 1,
	}
	consumer.queueRecvThread = NewSrsQueueRecvThread(consumer, c.rtmp)
	consumer.queueRecvThread.Start()
	return consumer
}
//有两个协程需要处理，这里的cycle和queueRecvThread
func (this *SrsConsumer) PlayCycle() error {
	for {
		for !this.queueRecvThread.Empty() {//process signal message
			msg := this.queueRecvThread.GetMsg()
			if msg != nil {
				err := this.process_play_control_msg(msg)
				if err != nil {
					return err
				}
			}
		}
		//todo process trd error
		//todo process realtime stream
		msg, err := this.queue.Wait()
		if err != nil {
			return err
		}

		if msg != nil {
			// fmt.Println("send to consumer")
			if msg.GetHeader().IsVideo() {
				//fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxsendmsg video");
				if flvcodec.VideoIsKeyFrame(msg.GetPayload()) {
					// fmt.Println("send key frame")
				}
				// fmt.Println("timestamp=", msg.GetHeader().GetTimestamp())
			} else {
			}
			err := this.conn.rtmp.SendMsg(msg, this.StreamId)
			_ = err
		}
	}

	return nil
}

func (this *SrsConsumer) StopPlay() error {
	this.source.RemoveConsumer(this)
	this.conn.Close()
	this.queueRecvThread.Stop()
	this.queue.Break()
	return nil
}

func (this *SrsConsumer) OnRecvError(err error) {
	this.StopPlay()
}

func (this *SrsConsumer) process_play_control_msg(msg *rtmp.SrsRtmpMessage) error {
	if !msg.GetHeader().IsAmf0Command() && !msg.GetHeader().IsAmf3Command() {
		//ignore 
		return nil
	}
	
	pkt, err := this.conn.rtmp.DecodeMessage(msg)
	if err != nil {
		return err
	}
	//todo add callpacket 
	//todo process pause message
	switch pkt.(type) {
	case *packet.SrsCloseStreamPacket:{
		//todo fix close stream action
		return errors.New("get close stream packet")
	}
	case *packet.SrsPausePacket:{
		return nil
	}
	}
	return nil
}

//todo add rtmp jitter algorithm
func (this *SrsConsumer) Enqueue(msg *rtmp.SrsRtmpMessage, atc bool, jitterAlgorithm *SrsRtmpJitterAlgorithm) {
	this.queue.Enqueue(msg)
}

