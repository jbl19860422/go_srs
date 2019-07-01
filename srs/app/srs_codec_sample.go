package app

import (
	// "fmt"
	"errors"
	"go_srs/srs/codec"
)
type SrsCodecSample struct {
	IsVideo			bool
	SampleUnits		[]([]byte)
	/**
    * CompositionTime, video_file_format_spec_v10_1.pdf, page 78.
    * cts = pts - dts, where dts = flvheader->timestamp.
    */
	Cts				int32
	//video specify
	FrameType		codec.SrsCodecVideoAVCFrame
	AvcPacketType	codec.SrsCodecVideoAVCType
	//
	HasIdr			bool
	HasAud			bool
	HasSpsPps		bool

	//
	FirstNaluType	codec.SrsAvcNaluType
	//audio specify
	ACodec			codec.SrsCodecAudio
	SoundRate		codec.SrsCodecAudioSampleRate
	SoundSize		codec.SrsCodecAudioSampleSize
	SoundType		codec.SrsCodecAudioSoundType
	AacPacketType	codec.SrsCodecAudioType
}

const SRS_SRS_MAX_CODEC_SAMPLE = 128

func NewSrsCodecSample() *SrsCodecSample {
	return &SrsCodecSample{
		SampleUnits:make([]([]byte), 0),
	}
}

func (this *SrsCodecSample) AddSampleUnit(data []byte) error {
	if len(this.SampleUnits) > SRS_SRS_MAX_CODEC_SAMPLE {
		return errors.New("hls decode samples error, exceed the max count")
	}

	this.SampleUnits = append(this.SampleUnits, data)
	if this.IsVideo {
		nalUnitType := codec.SrsAvcNaluType(data[0] & 0x1f)
		if nalUnitType == codec.SrsAvcNaluTypeIDR {
			this.HasIdr = true
		} else if nalUnitType == codec.SrsAvcNaluTypeSPS || nalUnitType == codec.SrsAvcNaluTypePPS {
			this.HasSpsPps = true
		} else if nalUnitType == codec.SrsAvcNaluTypeAccessUnitDelimiter {
			this.HasAud = true
		}

		if this.FirstNaluType == codec.SrsAvcNaluTypeReserved {
			this.FirstNaluType = nalUnitType
		}
	}
	return nil
}

func (this *SrsCodecSample) Clear() {
	this.IsVideo = false
	this.SampleUnits = this.SampleUnits[0:0]
    this.Cts = 0
    this.FrameType = codec.SrsCodecVideoAVCFrameReserved
	this.AvcPacketType = codec.SrsCodecVideoAVCTypeReserved
	this.HasIdr = false
	this.HasAud  = false
	this.HasSpsPps = false
	this.FirstNaluType = codec.SrsAvcNaluTypeReserved
    
    this.ACodec = codec.SrsCodecAudioReserved1
    this.SoundRate = codec.SrsCodecAudioSampleRateReserved
    this.SoundSize = codec.SrsCodecAudioSampleSizeReserved
    this.SoundType = codec.SrsCodecAudioSoundTypeReserved
    this.AacPacketType = codec.SrsCodecAudioTypeReserved
}

func (this *SrsCodecSample) SetIsVideo(v bool) {
	this.IsVideo = v
}

