package app

import (
	"go_srs/srs/codec/flv"
	"go_srs/srs/protocol/rtmp"
)

type SrsDvr struct {
	source 	*SrsSource
	plan	*flvcodec.SrsDvrPlan
}

func NewSrsDvr() *SrsDvr {
	return &SrsDvr{}
}

func (this *SrsDvr) Initialize(s *SrsSource, r *SrsRequest) error {
	this.source = s
	this.plan = flvcodec.NewSrsDvrPlan("./record.flv")
	//todo fix 
	this.plan.Initialize()
	return nil
}

func (this *SrsDvr) on_meta_data(metaData *rtmp.SrsRtmpMessage) error {
	return this.plan.On_meta_data(metaData)
}

func (this *SrsDvr) on_video(video *rtmp.SrsRtmpMessage) error {
	return this.plan.On_video(video)
}

func (this *SrsDvr) on_audio(audio *rtmp.SrsRtmpMessage) error {
	return this.plan.On_audio(audio)
}

func (this *SrsDvr) Close() {
	this.plan.Close()
}