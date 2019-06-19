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

//todo wait until reqired msg count recv
func (this *SrsConsumer) Wait(minCount uint32, duration uint32) *rtmp.SrsRtmpMessage {
	return this.queue.Wait()
}