package app

type SrsTsCache struct {
	audio *SrsTsMessage
	video *SrsTsMessage
}

func (this *SrsTsCache) cache_audio(codec *SrsAvcAacCodec, dts int64, sample SrsCodecSample) error {
	if this.audio == nil {
		this.audio = NewSrsTsMessage()
		this.audio.writePcr = false
		this.audio.dts = dts
		this.audio.pts = dts
		this.audio.startPts = dts
	}

	this.audio.sid = SrsTsPESStreamIdAudioCommon
	return nil
}
