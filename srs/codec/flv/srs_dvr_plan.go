package flvcodec

import (
	"go_srs/srs/protocol/rtmp"
	"go_srs/srs/protocol/packet"
)

type SrsDvrPlan struct {
	segment *SrsFlvSegment
}

func NewSrsDvrPlan() *SrsDvrPlan {
	return &SrsDvrPlan{
		segment: NewSrsFlvSegment(),
	}
}

func (this *SrsDvrPlan) Initialize() {
	this.segment.Initialize()
}

func (this *SrsDvrPlan) OnVideo(msg *rtmp.SrsRtmpMessage) {
	this.segment.WriteVideo(msg)
}

func (this *SrsDvrPlan) OnAudio(msg *rtmp.SrsRtmpMessage) {
	this.segment.WriteAudio(msg)
}

func (this *SrsDvrPlan) OnMetaData(pkt *packet.SrsOnMetaDataPacket) {
	this.segment.WriteMetaData(pkt)
}
