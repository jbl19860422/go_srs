package app

import (
	"go_srs/srs/protocol/rtmp"
)

type SrsConsumer struct {
	source 	*SrsSource
	conn 	*SrsRtmpConn
	queue 	*SrsMessageQueue
}

func NewSrsConsumer(s *SrsSource, c *SrsRtmpConn) *SrsConsumer {
	return &SrsConsumer{
		queue:NewSrsMessageQueue(),
		source:s,
		conn:c,
	}
}

//todo add rtmp jitter algorithm
func (this *SrsConsumer) Enqueue(msg *rtmp.SrsRtmpMessage, atc bool) {
	this.queue.Enqueue(msg)
}