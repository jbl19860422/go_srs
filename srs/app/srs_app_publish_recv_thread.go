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

// import (
// 	"go_srs/srs/protocol/rtmp"
// )

// type SrsAppPublishRecvThread struct {
// 	recvThread *SrsRecvThread
// 	rtmp       *rtmp.SrsRtmpServer
// 	req        *SrsRequest
// 	conn       *SrsRtmpConn
// 	source     *SrsSource
// 	isFmle     bool
// 	isEdge     bool
// }

// func NewSrsAppPublishRecvThread(s *rtmp.SrsRtmpServer, r *SrsRequest, c *SrsRtmpConn, source_ *SrsSource, isFmle_ bool, isEdge_ bool) *SrsAppPublishRecvThread {
// 	st := &SrsAppPublishRecvThread{
// 		rtmp:   s,
// 		req:    r,
// 		conn:   c,
// 		source: source_,
// 		isFmle: isFmle_,
// 		isEdge: isEdge_,
// 	}
// 	st.recvThread = NewSrsRecvThread(s, st, 1000)
// 	return st
// }

// func (this *SrsAppPublishRecvThread) Start() {
// 	this.recvThread.Start()
// }

// func (this *SrsAppPublishRecvThread) CanHandle() bool {
// 	return true
// }
// func (this *SrsAppPublishRecvThread) Handle(msg *rtmp.SrsRtmpMessage) error {

// 	//todo fix cid change
// 	//todo nbmsg++
// 	// err := this.conn.HandlePublishMessage(this.source, msg, this.isFmle, this.isEdge)
// 	// return err
// 	return nil
// }

// func (this *SrsAppPublishRecvThread) OnRecvError(err error) {
// 	this.rtmp.OnRecvError(err)
// 	return
// }

// func (this *SrsAppPublishRecvThread) OnThreadStart() {
// 	return
// }

// func (this *SrsAppPublishRecvThread) OnThreadStop() {
// 	return
// }
