package app

import (
	"os"
	"errors"
	"go_srs/srs/codec"
	"go_srs/srs/utils"
)

type SrsTsContext struct {
	ready      bool
	pids       map[int]*SrsTsChannel
	pure_audio bool
	vcodec     codec.SrsCodecVideo
	acodec     codec.SrsCodecAudio
	file 	   *os.File
}

func NewSrsTsContext() *SrsTsContext {
	f, err := os.OpenFile("aa.ts", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil
	}
	f.Truncate(0)

	return &SrsTsContext{
		ready:false,
		file:f,
	}
}

func (this *SrsTsContext) Get(pid int) *SrsTsChannel {
	c, ok := this.pids[pid]
	if ok {
		return c
	}
	return nil
}

func (this *SrsTsContext) Set(pid int, applyPid SrsTsPidApply, stream SrsTsStream) {
	c, ok := this.pids[pid]
	if ok {
		c.pid = pid
		c.apply = applyPid
		c.stream = stream
	} else {
		c = NewSrsTsChannel()
		c.pid = pid
		c.apply = applyPid
		c.stream = stream
		this.pids[pid] = c
	}
}

func (this *SrsTsContext) Encode(msg *SrsTsMessage, vc codec.SrsCodecVideo, ac codec.SrsCodecAudio) error {
	var vs, as SrsTsStream
	var videoPid, audioPid int16
	switch vc {
	case codec.SrsCodecVideoAVC:
		vs = SrsTsStreamVideoH264
		videoPid = TS_VIDEO_AVC_PID
	default:
		vs = SrsTsStreamReserved
	}

	switch ac {
	case codec.SrsCodecAudioAAC:
		as = SrsTsStreamAudioAAC
		audioPid = TS_AUDIO_AAC_PID
	case codec.SrsCodecAudioMP3:
		as = SrsTsStreamAudioMp3
		audioPid = TS_AUDIO_MP3_PID
	default:
		as = SrsTsStreamReserved
	}

	if as == SrsTsStreamReserved || vs == SrsTsStreamReserved {
		return errors.New("not support as or vs")
	}

	if this.vcodec != vc || this.acodec != ac {
		this.vcodec = vc
		this.acodec = ac
		err := this.encodePatPmt(videoPid, vs, audioPid, as)
		if err != nil {
			return err
		}
	}

	noVideo := false
	if vs == SrsTsStreamReserved {
		noVideo = true
	}

	if msg.IsAudio() {
		this.encodePes(msg, audioPid, as, noVideo)
	} else {
		this.encodePes(msg, videoPid, vs, noVideo)
	}
	return nil
}

func (this *SrsTsContext) encodePatPmt(vpid int16, vs SrsTsStream, apid int16, as SrsTsStream) error {
    if vs != SrsTsStreamVideoH264 && as != SrsTsStreamAudioAAC && as != SrsTsStreamAudioMp3 {
		return errors.New("invalid video stream or audio stream type")
    }

    var pmt_number int16 = TS_PMT_NUMBER
    var pmt_pid int16 = TS_PMT_PID
    if true {
        pkt := CreatePAT(this, pmt_number, pmt_pid)
        stream := utils.NewSrsStream([]byte{})
		pkt.Encode(stream)
		this.file.Write(stream.Data())
	}
	
	if true {
		pkt := CreatePMT(this, pmt_number, pmt_pid, vpid, vs, apid, as)
        stream := utils.NewSrsStream([]byte{})
		pkt.Encode(stream)
		this.file.Write(stream.Data())
	}
    // When PAT and PMT are writen, the context is ready now.
	this.ready = true
	return nil
}	

func (this *SrsTsContext) encodePes(msg *SrsTsMessage, pid int16, stream SrsTsStream, no_video bool) {

}



