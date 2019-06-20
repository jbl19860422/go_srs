package app

import(
	"go_srs/srs/protocol/rtmp"
	// "go_srs/srs/codec"
	"fmt"
	"go_srs/srs/codec/flv"
)

type SrsMessageQueue struct {
	ignoreShrink 	bool
	avStartTime		int64
	avEndTime		int64
	queueSizeMs		int

	msgs 			[]*rtmp.SrsRtmpMessage
	msgCount 		chan int
}

func NewSrsMessageQueue() *SrsMessageQueue {
	return &SrsMessageQueue{
		ignoreShrink:true,
		avStartTime:0,
		avEndTime:0,
		queueSizeMs:0,
		msgs:make([]*rtmp.SrsRtmpMessage, 0),
		msgCount:make(chan int, 10000),
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

func (this *SrsMessageQueue) Wait() *rtmp.SrsRtmpMessage {
	<- this.msgCount
	// fmt.Println("msgCount=", len(this.msgs), "&a=", a)
	if len(this.msgs) <= 0 {
		return nil
	}

	msg := this.msgs[0]
	this.msgs = this.msgs[1:]
	if msg == nil {
		fmt.Println("msg is nil")
	}
	return msg
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
		videoSH.SetTimestamp(this.avEndTime)
		this.msgs = append(this.msgs, videoSH)
	}

	if audioSH != nil {
		audioSH.SetTimestamp(this.avEndTime)
		this.msgs = append(this.msgs, audioSH)
	}
}

func (this *SrsMessageQueue) Clear() {
	this.msgs = this.msgs[0:0]
	this.avStartTime = -1
	this.avEndTime = -1
}
