package app

import (
	"os"
	"fmt"
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
	f, err := os.OpenFile("c.ts", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return nil
	}
	f.Truncate(0)
	fmt.Println("*************************NewSrsTsContext********************************")
	return &SrsTsContext{
		ready: false,
		file:f,
		pids:make(map[int]*SrsTsChannel),
	}
}

func (this *SrsTsContext) Reset() {
	this.ready = false
	this.vcodec = codec.SrsCodecVideoReserved
	this.acodec = codec.SrsCodecAudioReserved1
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
	as = SrsTsStreamAudioAAC
	audioPid = TS_AUDIO_AAC_PID
	if as == SrsTsStreamReserved || vs == SrsTsStreamReserved {
		return errors.New("not support as or vs")
	}

	if this.vcodec != vc || this.acodec != ac {
		this.vcodec = vc
		this.acodec = ac
		fmt.Println("videopid=", videoPid, "&audioPid=", audioPid)
		err := this.encodePatPmt(videoPid, vs, audioPid, as)
		if err != nil {
			return err
		}
		// os.Exit(0)
	}

	noVideo := false
	if vs == SrsTsStreamReserved {
		noVideo = true
	}
	
	if msg.IsAudio() {
		// fmt.Println("****************encodePes audio****************")
		this.encodePes(msg, audioPid, as, -1, noVideo)
	} else {
		// fmt.Println("****************encodePes video****************")
		this.encodePes(msg, videoPid, vs, msg.dts, noVideo)
	}
	return nil
}

func (this *SrsTsContext) encodePatPmt(vpid int16, vs SrsTsStream, apid int16, as SrsTsStream) error {
	if vs != SrsTsStreamVideoH264 && as != SrsTsStreamAudioAAC && as != SrsTsStreamAudioMp3 {
		return errors.New("invalid video stream or audio stream type")
	}
	fmt.Println("************encodePatPmt*****************")
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

func (this *SrsTsContext) encodePes(msg *SrsTsMessage, pid int16, sid SrsTsStream, pcr int64, pure_audio bool) error {
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
	
	channel := this.Get(int(pid))
	_ = channel
	// left := len(msg.payload)
	// pcr := msg.dts
	pkts := CreatePes(this, pid, msg.sid, &channel.continuityCounter, 0, pcr, msg.dts, msg.pts, msg.payload)
	
	for i := 0; i < len(pkts); i++ {
		s := utils.NewSrsStream([]byte{})
		pkts[i].Encode(s)
		n, err := this.file.Write(s.Data())
		if len(s.Data()) != 188 {
			fmt.Println("errrrrrrrrrrrrrrrrrrrrrrrrrrrpayload_len=", len(msg.payload), "&pkts_count=", len(pkts), "&data_len=", len(s.Data()))
		}
		// if pid == TS_AUDIO_AAC_PID {
		// 	fmt.Println("payload_len=", len(msg.payload), "&pkts_count=", len(pkts), "&data_len=", len(s.Data()))
		// }
		// fmt.Printf("data[0]=%x, data[1]=%x, data[2]=%x, data[3]=%x, len=%d\n", s.Data()[0], s.Data()[1], s.Data()[2], s.Data()[3], len(s.Data()))
		if err != nil {
			fmt.Println("***********************write file err, n=", n, "********************")
		}
	}
	return nil
}
