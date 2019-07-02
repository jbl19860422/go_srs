package app

import (
	// "fmt"
	// "os"
	"go_srs/srs/protocol/rtmp"
	"go_srs/srs/codec"
	"go_srs/srs/utils"
	"time"
)

type SrsHls struct {
	muxer    *SrsHlsMuxer
	hlsCache *SrsHlsCache

	req		 *SrsRequest
	lastUpdateTime	int64

	source	*SrsSource
	sample   *SrsCodecSample
	codec    *SrsAvcAacCodec
	context	  	*SrsTsContext

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

func (this *SrsHls) Initialize() error {
	this.streamDts = 0
	//todo fix
	return nil
}

func (this *SrsHls) dispose() {
	this.muxer.dispose()
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

func (this *SrsHls) initialize(s *SrsSource, r *SrsRequest) error {
	this.source = s
	this.req = r
	err := this.muxer.initialize()
	if err != nil {
		return err
	}
	return nil
}

func (this *SrsHls) on_publish(req *SrsRequest, fetch_sequence_header bool) error {
	//todo 
	this.lastUpdateTime = utils.GetCurrentMs()

	err := this.hlsCache.on_publish(this.muxer, this.req, this.streamDts) 
	if err != nil {
		return err
	}

	if fetch_sequence_header {
		err = this.source.on_hls_start() 
		if err != nil {
			return err
		}
	}
	return nil
}


func (this *SrsHls) Start() {
	go func() {
	DONE:
		for {
			select {
			case <-time.After(time.Second * 5):
				this.dispose()
			case <-this.exit:
				break DONE
			}
		}
		close(this.done)
	}()
}

func (this *SrsHls) Stop() {
	close(this.exit)
	<-this.done
}




func (this *SrsHls) on_meta_data(metaData *rtmp.SrsRtmpMessage) error {
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

	dts := video.GetHeader().GetTimestamp()*90
	this.streamDts = dts
	
	var add_sps_pps bool = false
	p := make([]byte, 0)
	for i := 0; i < len(this.sample.SampleUnits); i++ {
		nal_unit_type := codec.SrsAvcNaluType(this.sample.SampleUnits[i][0] & 0x1f)
		if !add_sps_pps && nal_unit_type == codec.SrsAvcNaluTypeIDR {
			p = append(p, []byte{0,0,1}...)
			p = append(p, this.codec.sequenceParameterSetNALUnit...)
			p = append(p, []byte{0,0,1}...)
			p = append(p, this.codec.pictureParameterSetNALUnit...)
			add_sps_pps = true
		}
		p = append(p, []byte{0,0,1}...)
		p = append(p, this.sample.SampleUnits[i]...)
	}
	ts := &SrsTsMessage{
		payload:p,
		dts:dts,
		pts:dts + int64(this.sample.Cts*90),
	}
	ts.sid = SrsTsPESStreamIdVideoCommon
	this.context.Encode(ts, codec.SrsCodecVideo(this.codec.videoCodecId), codec.SrsCodecAudio(this.codec.audioCodecId))
	
	return nil
	// err = this.hlsCache.WriteVideo(this.codec, this.muxer, dts, this.sample)
	// if err != nil {
	// 	return err
	// }

	// return nil
}

func (this *SrsHls) on_audio(audio *rtmp.SrsRtmpMessage) error {
	this.lastUpdateTime = utils.GetCurrentMs()

	this.sample.Clear()

	err := this.codec.audio_aac_demux(audio.GetPayload(), this.sample)
	if err != nil {
		return err
	}

	acodec := codec.SrsCodecAudio(this.codec.audioCodecId)
	// fmt.Println("acodec=",acodec)
	//not support
	if acodec != codec.SrsCodecAudioAAC && acodec != codec.SrsCodecAudioMP3 {
		return nil
	}

	p := make([]byte, 0)
	// fmt.Println("audiothis.sample.SampleUnits=", len(this.sample.SampleUnits))
	for i := 0; i < len(this.sample.SampleUnits); i++ {
	 	var frame_length int32 = int32(len(this.sample.SampleUnits[i]) + 7)
        // AAC-ADTS
        // 6.2 Audio Data Transport Stream, ADTS
        // in aac-iso-13818-7.pdf, page 26.
        // fixed 7bytes header
        var adts_header = []byte{0xff, 0xf9, 0x00, 0x00, 0x00, 0x0f, 0xfc}
        // profile, 2bits
        aac_profile := codec.SrsAacProfileLC //srs_codec_aac_rtmp2ts(codec->aac_object);
        adts_header[2] = (byte(aac_profile) << 6) & 0xc0
        // sampling_frequency_index 4bits
        adts_header[2] |= byte((this.codec.aacSampleRateIndex << 2) & 0x3c)
        // channel_configuration 3bits
        adts_header[2] |= byte((this.codec.aacChannels >> 2) & 0x01)
        adts_header[3] = byte((int32(this.codec.aacChannels) << 6) & 0xc0)
        // frame_length 13bits
        adts_header[3] |= byte((frame_length >> 11) & 0x03)
        adts_header[4] = byte((frame_length >> 3) & 0xff)
        adts_header[5] = byte(((frame_length << 5) & 0xe0))
        // adts_buffer_fullness; //11bits
        adts_header[5] |= 0x1f

        // copy to audio buffer
        p = append(p, adts_header...)
        p = append(p, this.sample.SampleUnits[i]...)
	}
	
	dts := audio.GetHeader().GetTimestamp()*90
	ts := &SrsTsMessage{
		payload:p,
		dts:dts,
		pts:dts,
	}

	ts.sid = SrsTsPESStreamIdAudioCommon
	this.context.Encode(ts, codec.SrsCodecVideo(this.codec.videoCodecId), codec.SrsCodecAudio(this.codec.audioCodecId))

	return nil
	// err = this.muxer.update_acodec(acodec)
	// if err != nil {
	// 	return err
	// }
	// return nil
}

func (this *SrsHls) Close() {
}
