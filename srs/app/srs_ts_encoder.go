package app

import (
	"io"
	"go_srs/srs/codec"
	"errors"
	"fmt"
)

type SrsTsEncoder struct {
	codec 	*SrsAvcAacCodec
	sample	*SrsCodecSample
	tsCache	*SrsTsCache
	context *SrsTsContext
	muxer 	*SrsTsMuxer

	writer	io.Writer
}

func NewSrsTsEncoder(w io.Writer) *SrsTsEncoder {
	c := NewSrsAvcAacCodec()
	s := NewSrsCodecSample()
	cache := NewSrsTsCache()
	context := NewSrsTsContext()
	m := NewSrsTsMuxer(w, context, codec.SrsCodecAudioAAC, codec.SrsCodecVideoAVC)
	return &SrsTsEncoder{
		codec:c,
		sample:s,
		tsCache:cache,
		context:context,
		writer:w,
		muxer:m,
	}
}

func (this *SrsTsEncoder) WriteHeader() error {
	return nil
}

func (this *SrsTsEncoder) WriteAudio(timestamp uint32, data []byte) (uint32, error) {
	this.sample.Clear()
	if err := this.codec.audio_aac_demux(data, this.sample); err != nil {
		//if err := this.codec.audio_mp3_demux(data, this.sample); err != nil {
		//	return 0, err
		//}
		fmt.Println("demux aac error", err)
		return 0, err
	}

	acodec := codec.SrsCodecAudio(this.codec.audioCodecId)
	if acodec != codec.SrsCodecAudioAAC && acodec != codec.SrsCodecAudioMP3 {
		fmt.Println("audio format error, need aac or mp3")
		return 0, errors.New("audio format error, need aac or mp3")
	}

	this.muxer.UpdateACodec(acodec)
	if acodec == codec.SrsCodecAudioAAC && this.sample.AacPacketType == codec.SrsCodecAudioTypeSequenceHeader {
		return 0, nil	//ignore aac sequence header
	}

	dts := int64(timestamp * 90)
	if err := this.tsCache.cache_audio(this.codec, dts, this.sample); err != nil {
		return 0, err
	}

	this.muxer.WriteAudio(this.tsCache.audio)

	return 0, nil
}

func (this *SrsTsEncoder) WriteVideo(timestamp uint32, data []byte) (uint32, error) {
	this.sample.Clear()
	if err := this.codec.video_avc_demux(data, this.sample); err != nil {
		return 0, err
	}

	// ignore info frame,
	// @see https://github.com/ossrs/srs/issues/288#issuecomment-69863909
	if this.sample.FrameType == codec.SrsCodecVideoAVCFrameVideoInfoFrame {
		return 0, nil
	}

	if (this.codec.videoCodecId != codec.SrsCodecVideoAVC) {
		return 0, nil
	}

	// ignore sequence header
	if this.sample.FrameType == codec.SrsCodecVideoAVCFrameKeyFrame && this.sample.AvcPacketType == codec.SrsCodecVideoAVCTypeSequenceHeader {
		return 0, nil
	}

	dts := int64(timestamp * 90)
	if err := this.tsCache.cache_video(this.codec, dts, this.sample); err != nil {
		return 0, nil
	}
	// write video to cache.
	if err := this.muxer.WriteVideo(this.tsCache.video); err != nil {
		return 0, err
	}
	return 0, nil
}