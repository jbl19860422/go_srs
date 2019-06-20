package app

import(
	"go_srs/srs/protocol/rtmp"
	"fmt"
)
type SrsAppPublishRecvThread struct {
	recvThread 	*SrsRecvThread
	rtmp 		*rtmp.SrsRtmpServer
	req			*SrsRequest
	conn		*SrsRtmpConn
	source		*SrsSource
	isFmle		bool
	isEdge		bool
}

func NewSrsAppPublishRecvThread(s *rtmp.SrsRtmpServer, r *SrsRequest, c *SrsRtmpConn, source_ *SrsSource, isFmle_ bool, isEdge_ bool) *SrsAppPublishRecvThread {
	st := &SrsAppPublishRecvThread{
		rtmp:s,
		req:r,
		conn:c,
		source:source_,
		isFmle:isFmle_,
		isEdge:isEdge_,
	}
	st.recvThread = NewSrsRecvThread(s, st, 1000)
	return st
}

func (this *SrsAppPublishRecvThread) Start() {
	this.recvThread.Start()
}

func (this *SrsAppPublishRecvThread) CanHandle() bool {
	return true
}
func (this *SrsAppPublishRecvThread) Handle(msg *rtmp.SrsRtmpMessage) error {

	//todo fix cid change
	//todo nbmsg++
	err := this.conn.HandlePublishMessage(this.source, msg, this.isFmle, this.isEdge)
	return err
}

func (this *SrsAppPublishRecvThread) OnRecvError(err error) {
	fmt.Println("OnRecvErr=", err)
	this.conn.OnRecvError(err)
	return
}	

func (this *SrsAppPublishRecvThread) OnThreadStart() {
	return
}

func (this *SrsAppPublishRecvThread) OnThreadStop() {
	return
}

