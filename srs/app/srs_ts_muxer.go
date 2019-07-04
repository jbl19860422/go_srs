package app

import (
	//"os"
	"go_srs/srs/codec"
	"io"
	"go_srs/srs/utils"
	"errors"
)

type SrsTsMuxer struct {
	as 			SrsTsStream
	vs 			SrsTsStream
	ready		bool
	audioPid	int16
	videoPid	int16
	wrotePatPmt	bool
	pids       	map[int]*SrsTsChannel

	context		*SrsTsContext
	writer		io.Writer
}

func NewSrsTsMuxer(w io.Writer, c *SrsTsContext, ac codec.SrsCodecAudio, vc codec.SrsCodecVideo) *SrsTsMuxer {
	muxer := &SrsTsMuxer{
		writer:w,
		context:c,
		wrotePatPmt:false,
		ready:false,
		pids:make(map[int]*SrsTsChannel),
	}

	muxer.convertACodecToTsStream(ac)
	muxer.convertVCodecToTsStream(vc)
	return muxer
}

func (this *SrsTsMuxer) convertACodecToTsStream(ac codec.SrsCodecAudio) {
	var astream SrsTsStream
	switch ac {
	case codec.SrsCodecAudioAAC:
		astream= SrsTsStreamAudioAAC
		this.audioPid = TS_AUDIO_AAC_PID
	case codec.SrsCodecAudioMP3:
		astream = SrsTsStreamAudioMp3
		this.audioPid = TS_AUDIO_MP3_PID
	default:
		astream = SrsTsStreamReserved
	}

	if astream != this.as {
		this.wrotePatPmt = false 	//rewrite pat pmt
	}
}

func (this *SrsTsMuxer) convertVCodecToTsStream(vc codec.SrsCodecVideo) {
	switch vc {
	case codec.SrsCodecVideoAVC:
		this.vs = SrsTsStreamVideoH264
		this.videoPid = TS_VIDEO_AVC_PID
	default:
		this.vs = SrsTsStreamReserved
	}
}

func (this *SrsTsMuxer) open() error {
	this.close()
	this.context.Reset()
	return nil
}

func (this *SrsTsMuxer) close() {
}

func (this *SrsTsMuxer) WriteAudio(audio *SrsTsMessage) error {
	if err := this.Encode(audio); err != nil {
		return err
	}
	return nil
}

func (this *SrsTsMuxer) WriteVideo(video *SrsTsMessage) error {
	if err := this.Encode(video); err != nil {
		return err
	}
	return nil
}

func (this *SrsTsMuxer) UpdateACodec(ac codec.SrsCodecAudio) error {
	this.convertACodecToTsStream(ac)
	return nil
}

func (this *SrsTsMuxer) Encode(msg *SrsTsMessage) error {
	if this.as == SrsTsStreamReserved || this.vs == SrsTsStreamReserved {
		return errors.New("not support as or vs")
	}

	if !this.wrotePatPmt {
		err := this.encodePatPmt()
		if err != nil {
			return err
		}
		this.wrotePatPmt = true
	}

	if this.vs == SrsTsStreamReserved {
		msg.writePcr = false
	}

	if msg.IsAudio() {
		this.encodePes(msg, this.audioPid, this.as, -1)
	} else {
		this.encodePes(msg, this.videoPid, this.vs, msg.dts)
	}
	return nil
}

func (this *SrsTsMuxer) encodePatPmt() error {
	if this.vs != SrsTsStreamVideoH264 && this.as != SrsTsStreamAudioAAC && this.as != SrsTsStreamAudioMp3 {
		return errors.New("invalid video stream or audio stream type")
	}

	var pmtNumber int16 = TS_PMT_NUMBER
	var pmtPid int16 = TS_PMT_PID
	if true {
		pkt := CreatePAT(this.context, pmtNumber, pmtPid)
		stream := utils.NewSrsStream([]byte{})
		pkt.Encode(stream)
		this.writer.Write(stream.Data())
	}

	if true {
		pkt := CreatePMT(this.context, pmtNumber, pmtPid, this.videoPid, this.vs, this.audioPid, this.as)
		stream := utils.NewSrsStream([]byte{})
		pkt.Encode(stream)
		this.writer.Write(stream.Data())
	}
	// When PAT and PMT are writen, the context is ready now.
	this.ready = true
	return nil
}

func (this *SrsTsMuxer) encodePes(msg *SrsTsMessage, pid int16, sid SrsTsStream, pcr int64) error {
	// Sometimes, the context is not ready(PAT/PMT write failed), error in this situation.
	if !this.ready {
		return errors.New("context not ready")
	}

	if len(msg.payload) <= 0 {
		return errors.New("msg length must not be zero")
	}

	if sid != SrsTsStreamVideoH264 && sid != SrsTsStreamAudioMp3 && sid != SrsTsStreamAudioAAC {
		return errors.New("ts: ignore the unknown stream")
	}

	channel := this.context.Get(int(pid))
	_ = channel
	// left := len(msg.payload)
	// pcr := msg.dts
	pkts := CreatePes(this.context, pid, msg.sid, &channel.continuityCounter, 0, pcr, msg.dts, msg.pts, msg.payload)

	for i := 0; i < len(pkts); i++ {
		s := utils.NewSrsStream([]byte{})
		pkts[i].Encode(s)
		n, err := this.writer.Write(s.Data())
		_ = err
		_ = n
		//if len(s.Data()) != 188 {
		//	fmt.Println("errrrrrrrrrrrrrrrrrrrrrrrrrrrpayload_len=", len(msg.payload), "&pkts_count=", len(pkts), "&data_len=", len(s.Data()))
		//}
	}
	return nil
}