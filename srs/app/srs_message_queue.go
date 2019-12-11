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

import(
	"fmt"
	"errors"
	"go_srs/srs/protocol/rtmp"
	"go_srs/srs/codec/flv"
)

type SrsMessageQueue struct {
	ignoreShrink 	bool
	avStartTime		int64
	avEndTime		int64
	queueSizeMs		int

	msgs 			[]*rtmp.SrsRtmpMessage
	msgCount 		chan int
	exit			chan bool
}

func NewSrsMessageQueue() *SrsMessageQueue {
	return &SrsMessageQueue{
		ignoreShrink:true,
		avStartTime:0,
		avEndTime:0,
		queueSizeMs:0,
		msgs:make([]*rtmp.SrsRtmpMessage, 0),
		msgCount:make(chan int, 10000),
		exit:make(chan bool),
	}
}

func (this *SrsMessageQueue) Enqueue(msg *rtmp.SrsRtmpMessage) {
	// if msg == nil {
	// 	fmt.Println("enque nil*******************")
	// } else {
	// 	fmt.Println("enqueue no nil*************")
	// }
	this.msgs = append(this.msgs, msg)
	this.msgCount <- len(this.msgs)
}

func (this *SrsMessageQueue) Size() int {
	return len(this.msgs)
}

func (this *SrsMessageQueue) Duration() int64 {
	return this.avEndTime - this.avStartTime
}

func (this *SrsMessageQueue) SetQueueSize(queueSize float64) {
	this.queueSizeMs = int(queueSize*1000)
} 

func (this *SrsMessageQueue) Empty() bool {
	return len(this.msgs) == 0
}

func (this *SrsMessageQueue) Break() {
	close(this.exit)
}

func (this *SrsMessageQueue) Wait() (*rtmp.SrsRtmpMessage, error) {
	select {
	case <- this.msgCount :
	{
		if len(this.msgs) <= 0 {
			return nil, nil
		}
	
		msg := this.msgs[0]
		this.msgs = this.msgs[1:]
		if msg == nil {
			fmt.Println("msg is nil")
		}
		return msg, nil
	}
	case <- this.exit :
	{
		fmt.Println("**************break from queue****************")
		return nil, errors.New("queue break")
	}
	}
}

//todo dump packets with jitter algorithm

/**
* remove a gop from the front.
* if no iframe found, clear it.
*/
func (this *SrsMessageQueue) Shrink() {
	var videoSH *rtmp.SrsRtmpMessage
	var audioSH *rtmp.SrsRtmpMessage
	for i := 0; i < len(this.msgs); i++ {
		//todo check is raw data?
		if this.msgs[i].GetHeader().IsVideo() && flvcodec.VideoIsSequenceHeader(this.msgs[i].GetPayload()) {
			videoSH = this.msgs[i]
		}

		if this.msgs[i].GetHeader().IsAudio() && flvcodec.AudioIsSequenceHeader(this.msgs[i].GetPayload()) {
			audioSH = this.msgs[i]
		}
	}
	//clear
	this.msgs = this.msgs[0:0]
	
	this.avStartTime = this.avEndTime
	if videoSH != nil {
		videoSH.GetHeader().SetTimestamp(this.avEndTime)
		this.msgs = append(this.msgs, videoSH)
	}

	if audioSH != nil {
		audioSH.GetHeader().SetTimestamp(this.avEndTime)
		this.msgs = append(this.msgs, audioSH)
	}
}

func (this *SrsMessageQueue) Clear() {
	this.msgs = this.msgs[0:0]
	this.avStartTime = -1
	this.avEndTime = -1
}
