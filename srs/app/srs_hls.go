package app

import (
	"go_srs/srs/protocol/rtmp"
	"time"
)

type SrsHls struct {
	muxer    *SrsHlsMuxer
	hlsCache *SrsHlsCache
	sample   *SrsCodecSample
	codec    *SrsAvcAacCodec
	startDts int64
	exit     chan bool
	done     chan bool
}

func NewSrsHls(fname string) *SrsHls {
	return &SrsHls{
		muxer:    NewSrsHlsMuxer(),
		hlsCache: NewSrsHlsCache(),
		sample:   NewSrsCodecSample(),
		codec:	  NewSrsAvcAacCodec(),
		exit:     make(chan bool),
		done:     make(chan bool),
	}
}

func (this *SrsHls) Initialize() error {
	this.startDts = 0
	//todo fix
	return nil
}

func (this *SrsHls) Start() {
	go func() {
	DONE:
		for {
			select {
			case <-time.After(time.Second * 5):
				this.dispose()
			case <-this.exit:
				break DONE
			}
		}
		close(this.done)
	}()
}

func (this *SrsHls) Stop() {
	close(this.exit)
	<-this.done
}

func (this *SrsHls) dispose() {
	this.muxer.dispose()
}

func (this *SrsHls) on_meta_data(metaData *rtmp.SrsRtmpMessage) error {
	return nil
}

func (this *SrsHls) on_video(video *rtmp.SrsRtmpMessage) error {
	this.sample.Clear()
	this.codec.video_avc_demux(video.GetPayload(), this.sample)
	return nil
}

func (this *SrsHls) on_audio(audio *rtmp.SrsRtmpMessage) error {
	return nil
}

func (this *SrsHls) Close() {
}
