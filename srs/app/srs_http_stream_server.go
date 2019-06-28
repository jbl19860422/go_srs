package app

import (
	"fmt"
	"io"
	"net/http"
	"go_srs/srs/protocol/rtmp"
	"go_srs/srs/codec/flv"
)

type SrsHttpStreamServer struct {
	sources map[string]*SrsSource
}

func NewSrsHttpStreamServer() *SrsHttpStreamServer {
	return &SrsHttpStreamServer{
		sources:make(map[string]*SrsSource),
	}
}

func (this *SrsHttpStreamServer) Mount(r *SrsRequest, s *SrsSource) {
	path := r.GetStreamUrl()
	path += ".flv"
	fmt.Println("????????????????path=", path, "??????????????????")
	this.sources[path] = s
}

func (this *SrsHttpStreamServer) CreateConsumer(s *SrsSource, w http.ResponseWriter, r *http.Request) Consumer {
	c := NewSrsHttpFlvConsumer(s, w, r)
	if err := s.AppendConsumer(c); err != nil {
		return nil
	}
	return c
}

type SrsHttpFlvConsumer struct {
	source          *SrsSource
	queue           *SrsMessageQueue
	StreamId		int
	writer			http.ResponseWriter
}

func NewSrsHttpFlvConsumer(s *SrsSource, w http.ResponseWriter, r *http.Request) *SrsHttpFlvConsumer {
	return &SrsHttpFlvConsumer{
		source:s,
		writer:w,
		queue:NewSrsMessageQueue(),
		StreamId:0,
	}
}

func (this *SrsHttpFlvConsumer) PlayCycle() error {
	for {
		// fmt.Println("***********SrsHttpFlvConsumer playing")
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
			// fmt.Println("send to consumer")
			if msg.GetHeader().IsVideo() {
				//fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxsendmsg video");
				if flvcodec.VideoIsKeyFrame(msg.GetPayload()) {
					fmt.Println("send key frame")
				}
				// fmt.Println("timestamp=", msg.GetHeader().GetTimestamp())
			} else {
			}
			//todo send msg to response writer
			// err := this.conn.rtmp.SendMsg(msg, this.StreamId)
			// _ = err
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

func (this *SrsHttpStreamServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("url=", r.URL.Path)
	source, ok := this.sources[r.URL.Path]
	if !ok {
		fmt.Println("not find for", r.URL.Path)
		for k, _ := range this.sources {
			fmt.Println("k=", k)
		}
		io.WriteString(w, "404")
		return
	}

	fmt.Println("*****************create consumer for", r.URL.Path)
	consumer := this.CreateConsumer(source, w, r)
	err := consumer.PlayCycle()
	_ = err
}
