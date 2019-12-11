/*
The MIT License (MIT)

Copyright (c) 2019 GOSRS(gosrs)

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

import (
	"time"
	"go_srs/srs/protocol/rtmp"
	"go_srs/srs/codec"
	"go_srs/srs/utils"
)

type SrsHls struct {
	muxer    *SrsHlsMuxer
	hlsCache *SrsHlsCache

	req		 *SrsRequest
	lastUpdateTime	int64

	source	 *SrsSource
	sample   *SrsCodecSample
	codec    *SrsAvcAacCodec
	context	 *SrsTsContext

	hlsDispose	int64			//second
	/**
    * we store the stream dts,
    * for when we notice the hls cache to publish,
    * it need to know the segment start dts.
    * 
    * for example. when republish, the stream dts will 
    * monotonically increase, and the ts dts should start 
    * from current dts.
    * 
    * or, simply because the HlsCache never free when unpublish,
    * so when publish or republish it must start at stream dts,
    * not zero dts.
    */
	streamDts int64
	exit     chan bool
	done     chan bool
}

func NewSrsHls(c *SrsTsContext) *SrsHls {
	return &SrsHls{
		muxer:    NewSrsHlsMuxer(),
		hlsCache: NewSrsHlsCache(),
		sample:   NewSrsCodecSample(),
		codec:	  NewSrsAvcAacCodec(),
		context:  c,
		streamDts:0,
		lastUpdateTime:0,
		hlsDispose:5000,//dispose every five second
		exit:     make(chan bool),
		done:     make(chan bool),
	}
}

func (this *SrsHls) dispose() {
	//this.muxer.dispose()
}

func (this *SrsHls) cycle() error {
	if this.lastUpdateTime <= 0 {
		this.lastUpdateTime = time.Now().UnixNano() / 1e6
	}

	if this.req == nil {
		return nil
	}

	//todo read hls dispose from config
	if utils.GetCurrentMs() - this.lastUpdateTime <= this.hlsDispose {

	}

	this.lastUpdateTime = utils.GetCurrentMs()

	this.dispose()
	return nil
}

func (this *SrsHls) Initialize(s *SrsSource, r *SrsRequest) error {
	this.source = s
	this.req = r
	err := this.muxer.initialize()
	if err != nil {
		return err
	}
	return nil
}

func (this *SrsHls) onPublish(req *SrsRequest, fetch_sequence_header bool) error {
	//todo 
	this.lastUpdateTime = utils.GetCurrentMs()

	err := this.hlsCache.onPublish(this.muxer, this.req, this.streamDts) 
	if err != nil {
		return err
	}

	if fetch_sequence_header {
		err = this.source.onHlsStart() 
		if err != nil {
			return err
		}
	}
	return nil
}


func (this *SrsHls) Start() {
	//go func() {
	//DONE:
	//	for {
	//		select {
	//		case <-time.After(time.Second * 5):
	//			this.dispose()
	//		case <-this.exit:
	//			break DONE
	//		}
	//	}
	//	close(this.done)
	//}()
}

func (this *SrsHls) Stop() {
	close(this.exit)
	<-this.done
}

func (this *SrsHls) onMetaData(metaData *rtmp.SrsRtmpMessage) error {
	return nil
}

func (this *SrsHls) on_video(video *rtmp.SrsRtmpMessage) error {
	this.lastUpdateTime = utils.GetCurrentMs()

	this.sample.Clear()
	err := this.codec.video_avc_demux(video.GetPayload(), this.sample)
	if err != nil {
		return err
	}

	if this.sample.FrameType == codec.SrsCodecVideoAVCFrameVideoInfoFrame {
		return nil
	}

	if this.codec.videoCodecId != codec.SrsCodecVideoAVC {
		return nil
	}
	
	if this.sample.FrameType == codec.SrsCodecVideoAVCFrameKeyFrame && this.sample.AvcPacketType == codec.SrsCodecVideoAVCTypeSequenceHeader {
		return this.hlsCache.on_sequence_header(this.muxer)
	}
	//todo add jitter
	dts := video.GetHeader().GetTimestamp()*90
	this.streamDts = dts
	
	if err := this.hlsCache.WriteVideo(this.codec, this.muxer, dts, this.sample); err != nil {
		return err
	}
	return nil
}

func (this *SrsHls) on_audio(audio *rtmp.SrsRtmpMessage) error {
	this.lastUpdateTime = utils.GetCurrentMs()

	this.sample.Clear()
	err := this.codec.audio_aac_demux(audio.GetPayload(), this.sample)
	if err != nil {
		return err
	}

	acodec := codec.SrsCodecAudio(this.codec.audioCodecId)
	//not support
	if acodec != codec.SrsCodecAudioAAC && acodec != codec.SrsCodecAudioMP3 {
		return nil
	}

	if err := this.muxer.update_acodec(acodec); err != nil {
		return err
	}

	if acodec == codec.SrsCodecAudioAAC && this.sample.AacPacketType == codec.SrsCodecAudioTypeSequenceHeader {
		return this.hlsCache.on_sequence_header(this.muxer)
	}
	//todo config jitter
	dts := int64(audio.GetHeader().GetTimestamp()*90)
	// for pure audio, we need to update the stream dts also.
	this.streamDts = dts

	if err := this.hlsCache.write_audio(this.codec, this.muxer, dts, this.sample); err != nil {
		return err
	}
	return nil
}

func (this *SrsHls) Close() {
}
