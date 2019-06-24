package flvcodec

import (
	"fmt"
	"go_srs/srs/protocol/rtmp"
	// "go_srs/srs/protocol/packet"
)

type SrsDvrPlan struct {
	segment *SrsFlvSegment
}

func NewSrsDvrPlan(fname string) *SrsDvrPlan {
	return &SrsDvrPlan{
		segment: NewSrsFlvSegment(fname),
	}
}

func (this *SrsDvrPlan) Initialize() {
	this.segment.Initialize()
}

func (this *SrsDvrPlan) On_video(msg *rtmp.SrsRtmpMessage) error {
	this.segment.WriteVideo(msg)
	return nil
}

func (this *SrsDvrPlan) On_audio(msg *rtmp.SrsRtmpMessage) error {
	this.segment.WriteAudio(msg)
	return nil
}

func (this *SrsDvrPlan) On_meta_data(msg *rtmp.SrsRtmpMessage) error {
	fmt.Println("**********************xxxx On_meta_dataxxxxxxxxxxxxxxxxxxxx")
	err := this.segment.WriteMetaData(msg)
	if err != nil {
		fmt.Println("fffffffffffffffffff ", err)
	}
	return nil
}

func (this *SrsDvrPlan) Close() {
	this.segment.Close()
}
