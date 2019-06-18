package flvcodec

import (
	"go_srs/srs/protocol/packet"
	"go_srs/srs/protocol/rtmp"
	"go_srs/srs/global"
	"go_srs/srs/utils"
)

type SrsFlvSegment struct {
	path 			string
	flvEncoder		*SrsFlvEncoder
	durationOffset	int32
	filesizeOffset	int32
	startTime		int64
	previousPktTime int64
	duration		int64
	streamDuration	int64
}

func NewSrsFlvSegment() *SrsFlvSegment {
	return &SrsFlvSegment{
		path:"./record.flv",
		flvEncoder:NewSrsFlvEncoder("./record.flv"),
		startTime:-1,
		previousPktTime:-1,
		duration:0,
		streamDuration:0,
	}
}

func (this *SrsFlvSegment) Initialize() {
	_ = this.flvEncoder.WriteHeader()
}

func (this *SrsFlvSegment) WriteMetaData(pkt *packet.SrsOnMetaDataPacket) error {
	// remove duration and filesize.
	pkt.Set("filesize", int32(0))
	pkt.Set("duration", float64(0))
	pkt.Set("service", global.RTMP_SIG_SRS_SERVER)

	stream := utils.NewSrsStream([]byte{})
	_ = pkt.Encode(stream)
	_, _ = this.flvEncoder.WriteMetaData(stream.Data())
	return nil
}

func (this *SrsFlvSegment) onUpdateDuration(msg *rtmp.SrsRtmpMessage) error {
	if this.startTime < 0 {
		this.startTime = msg.GetTimestamp()
	}

	if this.previousPktTime < 0 || this.previousPktTime > msg.GetTimestamp() {
		this.previousPktTime = msg.GetTimestamp()
	}

	this.duration += msg.GetTimestamp() - this.previousPktTime
	this.streamDuration += msg.GetTimestamp() - this.previousPktTime
	return nil
}

func (this *SrsFlvSegment) WriteAudio(msg *rtmp.SrsRtmpMessage) error {
	this.flvEncoder.WriteAudio(uint32(msg.GetTimestamp()), msg.GetPayload())
	this.onUpdateDuration(msg)
	return nil
}

func (this *SrsFlvSegment) WriteVideo(msg *rtmp.SrsRtmpMessage) error {
	this.flvEncoder.WriteVideo(uint32(msg.GetTimestamp()), msg.GetPayload())
	this.onUpdateDuration(msg)
	return nil
}

func (this *SrsFlvSegment) Close() error {
	return nil
}

func (this *SrsFlvSegment) updateMetaData() error {
	return nil
}