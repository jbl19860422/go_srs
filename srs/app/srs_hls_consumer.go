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
	"go_srs/srs/protocol/rtmp"
	"go_srs/srs/utils"
	"go_srs/srs/codec"
)

type SrsHlsConsumer struct {
	source          *SrsSource
	req 			*SrsRequest
	queue           *SrsMessageQueue
	codec    		*SrsAvcAacCodec
	sampler   		*SrsCodecSampler
	muxer    		*SrsHlsMuxer
	hlsCache 		*SrsHlsCache
	context	 		*SrsTsContext

	lastUpdateTime	int64
	streamDts 		int64
}

func NewSrsHlsConsumer(s *SrsSource, req *SrsRequest) *SrsHlsConsumer {
	return &SrsHlsConsumer{
		source:s,
		req:req,
		queue:NewSrsMessageQueue(),
		codec:NewSrsAvcAacCodec(),
		sampler:NewSrsCodecSampler(),
		muxer:NewSrsHlsMuxer(),
		hlsCache:NewSrsHlsCache(),
		context:NewSrsTsContext(),
	}
}

func (this *SrsHlsConsumer) OnPublish() error {
	this.muxer.initialize()

	this.lastUpdateTime = utils.GetCurrentMs()
	err := this.hlsCache.onPublish(this.muxer, this.req, this.streamDts)
	if err != nil {
		return err
	}

	return nil
}

func (this *SrsHlsConsumer) OnUnpublish() error {

	return nil
}

func (this *SrsHlsConsumer) ConsumeCycle() error {
	for {
		msg, err := this.queue.Wait()
		if err != nil {
			return err
		}

		if msg != nil {
			if msg.GetHeader().IsVideo() {
				if err := this.onVideo(msg); err != nil {
					return err
				}
			} else if msg.GetHeader().IsAudio() {
				if err := this.onAudio(msg); err != nil {
					return err
				}
			} else {
			}
		}
	}
	return nil
}

func (this *SrsHlsConsumer) onVideo(video *rtmp.SrsRtmpMessage) error {
	this.lastUpdateTime = utils.GetCurrentMs()

	this.sampler.Clear()
	err := this.codec.videoAvcDemux(video.GetPayload(), this.sampler)
	if err != nil {
		return err
	}

	if this.sampler.FrameType == codec.SrsCodecVideoAVCFrameVideoInfoFrame {
		return nil
	}

	if this.codec.videoCodecId != codec.SrsCodecVideoAVC {
		return nil
	}

	if this.sampler.FrameType == codec.SrsCodecVideoAVCFrameKeyFrame && this.sampler.AvcPacketType == codec.SrsCodecVideoAVCTypeSequenceHeader {
		return this.hlsCache.onSequenceHeader(this.muxer)
	}
	//todo add jitter
	dts := video.GetHeader().GetTimestamp()*90
	this.streamDts = dts

	if err := this.hlsCache.WriteVideo(this.codec, this.muxer, dts, this.sampler); err != nil {
		return err
	}
	return nil
}

func (this *SrsHlsConsumer) onAudio(audio *rtmp.SrsRtmpMessage) error {
	this.lastUpdateTime = utils.GetCurrentMs()

	this.sampler.Clear()
	err := this.codec.audioAACDemux(audio.GetPayload(), this.sampler)
	if err != nil {
		return err
	}

	acodec := codec.SrsCodecAudio(this.codec.audioCodecId)
	//not support
	if acodec != codec.SrsCodecAudioAAC && acodec != codec.SrsCodecAudioMP3 {
		return nil
	}

	if err := this.muxer.updateACodec(acodec); err != nil {
		return err
	}

	if acodec == codec.SrsCodecAudioAAC && this.sampler.AacPacketType == codec.SrsCodecAudioTypeSequenceHeader {
		return this.hlsCache.onSequenceHeader(this.muxer)
	}
	//todo config jitter
	dts := int64(audio.GetHeader().GetTimestamp()*90)
	// for pure audio, we need to update the stream dts also.
	this.streamDts = dts

	if err := this.hlsCache.writeAudio(this.codec, this.muxer, dts, this.sampler); err != nil {
		return err
	}
	return nil
}

func (this *SrsHlsConsumer) onMetadata(metaData *rtmp.SrsRtmpMessage) error {
	return nil
}

func (this *SrsHlsConsumer) StopConsume() error {
	this.source.RemoveConsumer(this)
	this.queue.Break()
	return nil
}

func (this *SrsHlsConsumer) OnRecvError(err error) {
	this.StopConsume()
}

func (this *SrsHlsConsumer) Enqueue(msg *rtmp.SrsRtmpMessage, atc bool, jitterAlgorithm *SrsRtmpJitterAlgorithm) {
	this.queue.Enqueue(msg)
}