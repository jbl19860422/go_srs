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

func (this *SrsHttpFlvConsumer) PlayCycle() error {
	this.flvEncoder.WriteHeader()
	go func() {
		notify := this.writer.(http.CloseNotifier).CloseNotify()
		<- notify
		this.StopPlay()
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

func (this *SrsHttpFlvConsumer) StopPlay() error {
	this.source.RemoveConsumer(this)
	//send connection close to response writer
	this.queue.Break()
	return nil
}

func (this *SrsHttpFlvConsumer) OnRecvError(err error) {
	this.StopPlay()
}

func (this *SrsHttpFlvConsumer) Enqueue(msg *rtmp.SrsRtmpMessage, atc bool, jitterAlgorithm *SrsRtmpJitterAlgorithm) {
	this.queue.Enqueue(msg)
}