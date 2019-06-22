package app

import (
	"fmt"
	"go_srs/srs/protocol/rtmp"
	"go_srs/srs/codec/flv"
	"go_srs/srs/protocol/packet"
	// "context"
	// "time"
	"errors"
)

type ConsumerStopListener interface {
	OnConsumerStop()
}

type SrsConsumer struct {
	source          *SrsSource
	conn            *SrsRtmpConn
	queue           *SrsMessageQueue
	StreamId		int
	queueRecvThread *SrsQueueRecvThread
}

func NewSrsConsumer(s *SrsSource, c *SrsRtmpConn) *SrsConsumer {
	//todo
	consumer := &SrsConsumer{
		queue:  NewSrsMessageQueue(),
		source: s,
		conn:   c,
		StreamId: 1,
	}
	consumer.queueRecvThread = NewSrsQueueRecvThread(consumer, c.rtmp)
	consumer.queueRecvThread.Start()
	return consumer
}
//有两个协程需要处理，这里的cycle和queueRecvThread
func (this *SrsConsumer) PlayingCycle() error {
	for {
		for !this.queueRecvThread.Empty() {//process signal message
			msg := this.queueRecvThread.GetMsg()
			if msg != nil {
				err := this.process_play_control_msg(msg)
				if err != nil {
					return err
				}
			}
		}
		//todo process trd error
		//todo process realtime stream
		msg, err := this.queue.Wait()
		if err != nil {
			fmt.Println("quit PlayingCycle")
			return err
		}

		if msg != nil {
			// fmt.Println("send to consumer")
			if msg.GetHeader().IsVideo() {
				//fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxsendmsg video");
				if flvcodec.VideoIsKeyframe(msg.GetPayload()) {
					// fmt.Println("send key frame")
				}
			} else {
				//fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxsendmsg audio");
			}
			
			err := this.conn.rtmp.SendMsg(msg, this.StreamId)
			_ = err
		}
	}

	return nil
}

func (this *SrsConsumer) StopPlay() error {
	fmt.Println("******************StopPlay*****************")
	this.source.RemoveConsumer(this)
	this.conn.Close()
	fmt.Println("******************StopPlay1*****************")
	this.queueRecvThread.Stop()
	fmt.Println("******************StopPlay2*****************")
	this.queue.Break()
	fmt.Println("******************StopPlay3*****************")
	// this.conn.Stop()
	fmt.Println("******************StopPlay4*****************")
	//this.queueRecvThread.Stop()
	
	//this.conn.RemoveSelf()
	fmt.Println("******************StopPlay5*****************")
	return nil
}

func (this *SrsConsumer) OnRecvError(err error) {
	this.StopPlay()
}

func (this *SrsConsumer) process_play_control_msg(msg *rtmp.SrsRtmpMessage) error {
	if !msg.GetHeader().IsAmf0Command() && !msg.GetHeader().IsAmf3Command() {
		//ignore 
		return nil
	}
	
	pkt, err := this.conn.rtmp.DecodeMessage(msg)
	if err != nil {
		return err
	}
	//todo add callpacket 
	//todo process pause message
	switch pkt.(type) {
	case *packet.SrsCloseStreamPacket:{
		//todo fix close stream action
		return errors.New("get close stream packet")
	}
	case *packet.SrsPausePacket:{
		return nil
	}
	}
	return nil
}


//todo add rtmp jitter algorithm
func (this *SrsConsumer) Enqueue(msg *rtmp.SrsRtmpMessage, atc bool) {
	this.queue.Enqueue(msg)
}

