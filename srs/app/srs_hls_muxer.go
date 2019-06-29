package app

import (
	"os"
	"path"
	"time"
	"go_srs/srs/codec"
	"go_srs/srs/utils"
)

type SrsHlsMuxer struct {
	hls_entry_prefix   string
	hls_path           string
	hls_ts_file        string
	hls_wait_keyframe  bool
	m3u8_dir           string
	hls_aof_ratio      float64
	hls_fragment       float64
	hls_window         float64
	hls_ts_floor       bool
	hls_cleanup		   bool
	m3u8_file		   string
	deviation_ts       int
	accept_floor_ts    float64
	previous_floor_ts  float64
	_sequence_no       int
	m3u8               string
	m3u8_url           string
	max_td             int
	should_write_cache bool
	should_write_file  bool
	segments           []*SrsHlsSegment
	current            *SrsHlsSegment
	acodec             codec.SrsCodecAudio
	context            *SrsTsContext
}

func NewSrsHlsMuxer() *SrsHlsMuxer {
	return &SrsHlsMuxer{
		context:nil,
	}
}

func (this *SrsHlsMuxer) initialize() error {
	return nil
}

func (this *SrsHlsMuxer) dispose() {
	for i := 0; i < len(this.segments); i++ {
		//todo unlink segments'full_path
	}
	this.segments = this.segments[0:0]

	if this.current != nil {
		path := this.current.full_path + ".tmp"
		//todo unlink path
		this.current = nil
		_ = path
	}

	//todo unlink m3u8
}

func (this *SrsHlsMuxer) sequence_no() int {
	return this._sequence_no
}

func (this *SrsHlsMuxer) ts_url() string {
	if this.current != nil {
		return this.current.uri
	}
	return ""
}

func (this *SrsHlsMuxer) duration() float64 {
	if this.current != nil {
		return this.current.duration
	}
	return 0
}

func (this *SrsHlsMuxer) deviation() int {
	if !this.hls_ts_floor {
		return 0
	}
	return this.deviation_ts
}

func (this *SrsHlsMuxer) update_config(entry_prefix string, p string,
								m3u8_file string, ts_file string, fragment float64,window float64, 
								ts_floor bool, aof_ratio float64, cleanup bool, wait_keyframe bool) error {
	this.hls_entry_prefix = entry_prefix
	this.hls_path = p
	this.hls_ts_file = ts_file
	this.hls_fragment = fragment
	this.hls_aof_ratio = aof_ratio
	this.hls_ts_floor = ts_floor
	this.hls_cleanup = cleanup
	this.hls_wait_keyframe = wait_keyframe
	this.previous_floor_ts = 0
	this.accept_floor_ts = 0
	this.hls_window = window
	this.deviation_ts = 0
	this.m3u8_url = utils.Srs_path_build_stream(this.m3u8_file, "aaa", "app", "test")
	this.m3u8 = p + "/" + this.m3u8_url
	//todo set max td
	this.m3u8_dir = path.Dir(this.m3u8)
	err := os.MkdirAll(this.m3u8_dir, os.ModePerm)
	return err
}

func (this *SrsHlsMuxer) on_sequence_header() error {
	// this.current.is_sequence_header = true
	return nil
}

func (this *SrsHlsMuxer) update_acodec(ac codec.SrsCodecAudio) error {
	this.acodec = ac
	return this.current.muxer.UpdateACodec(ac)
}

const SRS_JUMP_WHEN_PIECE_DEVIATION = 20

func (this *SrsHlsMuxer) segment_open(segment_start_dts int64) error {
	if this.current != nil {
		return nil
	}
	//todo
	default_acodec := codec.SrsCodecAudio(codec.SrsCodecAudioAAC)
	default_vcodec := codec.SrsCodecVideo(codec.SrsCodecVideoAVC)

	this.current = NewSrsHlsSegment(this.context, default_acodec, default_vcodec)
	this.current.sequence_no = this._sequence_no
	this._sequence_no++

	this.current.segment_start_dts = segment_start_dts

	ts_file := this.hls_ts_file
	ts_file = utils.Srs_path_build_stream(ts_file, "aaa", "app", "test")

	if this.hls_ts_floor {
		current_floor_ts := int64(((time.Now().UnixNano() / 1e6) / (1000 * 5)))

		if this.accept_floor_ts == 0 {
			this.accept_floor_ts = float64(current_floor_ts - 1)
		} else {
			this.accept_floor_ts++
		}

		if int64(this.accept_floor_ts - float64(current_floor_ts)) > SRS_JUMP_WHEN_PIECE_DEVIATION {
            this.accept_floor_ts = float64(current_floor_ts - 1)
		}
		
		this.deviation_ts = (int)(this.accept_floor_ts - float64(current_floor_ts))

		// dup/jmp detect for ts in floor mode.
        if int64(this.previous_floor_ts) != 0 && int64(this.previous_floor_ts) != current_floor_ts - 1 {

        }
        this.previous_floor_ts = float64(current_floor_ts);
		// we always ensure the piece is increase one by one.
		//todo ts file name replace
	}
	//todo tsfile append seq suffix
	this.current.full_path = this.hls_path + "/" + ts_file
	//add prefix

	// open temp ts file.
	tmp_file := this.current.full_path + ".tmp";
	_ = tmp_file
	//todo	
	// if err := this.current.muxer.open(tmp_file); err != nil {
	// 	return err
	// }
	
	return nil
}

