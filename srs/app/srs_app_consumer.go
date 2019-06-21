package app

import (
	"fmt"
	"go_srs/srs/protocol/rtmp"
)

type ConsumerStopListener interface {
	OnConsumerStop()
}

type SrsConsumer struct {
	source          *SrsSource
	conn            *SrsRtmpConn
	queue           *SrsMessageQueue
	queueRecvThread *SrsQueueRecvThread
}

func NewSrsConsumer(s *SrsSource, c *SrsRtmpConn) *SrsConsumer {
	//todo
	consumer := &SrsConsumer{
		queue:  NewSrsMessageQueue(),
		source: s,
		conn:   c,
	}
	consumer.queueRecvThread = NewSrsQueueRecvThread(consumer, c.rtmp)
	consumer.queueRecvThread.Start()
	return consumer
}

//todo add rtmp jitter algorithm
func (this *SrsConsumer) Enqueue(msg *rtmp.SrsRtmpMessage, atc bool) {
	this.queue.Enqueue(msg)
}

//todo wait until reqired msg count recv
func (this *SrsConsumer) Wait(minCount uint32, duration uint32) *rtmp.SrsRtmpMessage {
	return this.queue.Wait()
}

func (this *SrsConsumer) Stop() {
	fmt.Println("this.queueRecvThread.Stop() start")
	this.conn.Stop()
	//this.queueRecvThread.Stop()
	this.conn.RemoveSelf()
	fmt.Println("this.queueRecvThread.Stop() end")
}
