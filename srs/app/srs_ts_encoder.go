package app

import (
	"io"
	"go_srs/srs/codec"
	"errors"
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
	return &SrsTsEncoder{
		codec:c,
		sample:s,
		tsCache:cache,
		context:context,
		writer:w,
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
		return 0, err
	}

	acodec := codec.SrsCodecAudio(this.codec.audioCodecId)
	if acodec != codec.SrsCodecAudioAAC && acodec != codec.SrsCodecAudioMP3 {
		return 0, errors.New("audio format error, need aac or mp3")
	}

	this.muxer.update_acodec(acodec)
	if acodec == codec.SrsCodecAudioAAC && this.sample.AacPacketType == codec.SrsCodecAudioTypeSequenceHeader {
		return 0, nil	//ignore aac sequence header
	}

	dts := int64(timestamp * 90)
	if err := this.tsCache.cache_audio(this.codec, dts, this.sample); err != nil {
		return 0, err
	}

	return 0, nil
}

func (this *SrsTsEncoder) WriteVideo(timestamp uint32, data []byte) (uint32, error) {
	this.sample.Clear()
	return 0, nil
}