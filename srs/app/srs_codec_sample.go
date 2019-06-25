package app

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
	IsIdr			bool
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
		SampleUnits:make([]([]byte), 0)
	}
}

func (this *SrsCodecSample) AddSampleUnit(data []byte) error {
	if len(this.SampleUnits) > SRS_SRS_MAX_CODEC_SAMPLE {
		return errors.New("hls decode samples error, exceed the max count")
	}

	this.SampleUnits = append(this.SampleUnits, data)
	if this.IsVideo {
		nalUnitType := SrsAvcNaluType(data[0] & 0x1f)
		if nalUnitType == SrsAvcNaluTypeIDR {
			this.HasIdr = true
		} else if nalUnitType == SrsAvcNaluTypeSPS || nalUnitType == SrsAvcNaluTypePPS {
			this.HasSpsPps = true
		} else if nalUnitType == SrsAvcNaluTypeAccessUnitDelimiter {
			this.HasAud = true
		}

		if this.FirstNaluType == SrsAvcNaluTypeReserved {
			this.FirstNaluType = nalUnitType
		}
	}
	return nil
}

func (this *SrsCodecSample) SetIsVideo(v bool) {
	this.IsVideo = v
}

