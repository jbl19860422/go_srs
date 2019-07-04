package app

import (
	"net/http"
	"go_srs/srs/protocol/rtmp"
)

type SrsHttpTsConsumer struct {
	source          *SrsSource
	queue           *SrsMessageQueue
	StreamId		int
	writer			http.ResponseWriter
	tsEncoder		*SrsTsEncoder
}

func NewSrsHttpTsConsumer(s *SrsSource, w http.ResponseWriter, r *http.Request) *SrsHttpTsConsumer {
	return &SrsHttpTsConsumer{
		source:s,
		writer:w,
		queue:NewSrsMessageQueue(),
		StreamId:0,
		tsEncoder:NewSrsTsEncoder(w),
	}
}

func (this *SrsHttpTsConsumer) PlayCycle() error {
	this.tsEncoder.WriteHeader()
	go func() {
		notify := this.writer.(http.CloseNotifier).CloseNotify()
		<- notify
		this.StopPlay()
	}()

	for {
		// fmt.Println("***********SrsHttpTsConsumer playing")
		//todo process http message
		// for !this.queueRecvThread.Empty() {//process signal message
		// 	msg := this.queueRecvThread.GetMsg()
		// 	if msg != nil {
		// 		err := this.process_play_control_msg(msg)
		// 		if err != nil {
		// 			return err
		// 		}
		// 	}
		// }
		//todo process trd error
		//todo process realtime stream
		msg, err := this.queue.Wait()
		if err != nil {
			return err
		}

		if msg != nil {
			if msg.GetHeader().IsVideo() {
				this.tsEncoder.WriteVideo(uint32(msg.GetHeader().GetTimestamp()), msg.GetPayload())
			} else if msg.GetHeader().IsAudio() {
				this.tsEncoder.WriteAudio(uint32(msg.GetHeader().GetTimestamp()), msg.GetPayload())
			} else {
				//this.tsEncoder.WriteMetaData(msg.GetPayload())
			}
			//todo send msg to response writer
			// err := this.conn.rtmp.SendMsg(msg, this.StreamId)
			// _ = err
		}
	}
}

func (this *SrsHttpTsConsumer) StopPlay() error {
	this.source.RemoveConsumer(this)
	//send connection close to response writer
	this.queue.Break()
	return nil
}

func (this *SrsHttpTsConsumer) OnRecvError(err error) {
	this.StopPlay()
}

func (this *SrsHttpTsConsumer) Enqueue(msg *rtmp.SrsRtmpMessage, atc bool, jitterAlgorithm *SrsRtmpJitterAlgorithm) {
	this.queue.Enqueue(msg)
}