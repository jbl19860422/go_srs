/*
The MIT License (MIT)

Copyright (c) 2013-2015 GOSRS(gosrs)

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package app

import (
	"os"
	"path"
	"time"
	"go_srs/srs/codec"
	"go_srs/srs/utils"
	"strconv"
	"errors"
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
		context:NewSrsTsContext(),
	}
}

func (this *SrsHlsMuxer) initialize() error {
	return nil
}

func (this *SrsHlsMuxer) is_segment_overflow() bool {
	if this.current.duration * 1000 < 2 * 100 {
		return false
	}

	var deviation float64
	if this.hls_ts_floor {
		deviation = float64(0.3 * float64(this.deviation_ts) * float64(this.hls_fragment))
	} else {
		deviation = 0.0
	}

	if this.current.duration >= this.hls_fragment + deviation {
		return true
	}
	return false
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
	this.max_td = 10000
	this.m3u8_dir = path.Dir(this.m3u8)
	err := os.MkdirAll(this.m3u8_dir, os.ModePerm)
	return err
}

func (this *SrsHlsMuxer) flush_video(cache *SrsTsCache) error {
	if cache.video == nil || len(cache.video.payload) <= 0 {
		return errors.New("the len of video must not be 0")
	}

	this.current.UpdateDuration(cache.video.dts)

	if err := this.current.WriteVideo(cache.video); err != nil {
		return err
	}
	return nil
}

func (this *SrsHlsMuxer) flush_audio(cache *SrsTsCache) error {
	if this.current == nil {
		return nil
	}

	if cache.audio == nil || len(cache.audio.payload) <= 0 {
		return errors.New("error len of audio")
	}

	if err := this.current.WriteAudio(cache.audio); err != nil {
		return err
	}

	return nil
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

	this.current = NewSrsHlsSegment(this.context)
	this.current.sequence_no = this._sequence_no
	this._sequence_no++

	this.current.segment_start_dts = segment_start_dts

	//ts_file := this.hls_ts_file
	//ts_file = utils.Srs_path_build_stream(ts_file, "aaa", "app", "test")

	//if this.hls_ts_floor {
	//	current_floor_ts := int64(((time.Now().UnixNano() / 1e6) / (1000 * 5)))
	//
	//	if this.accept_floor_ts == 0 {
	//		this.accept_floor_ts = float64(current_floor_ts - 1)
	//	} else {
	//		this.accept_floor_ts++
	//	}
	//
	//	if int64(this.accept_floor_ts - float64(current_floor_ts)) > SRS_JUMP_WHEN_PIECE_DEVIATION {
     //       this.accept_floor_ts = float64(current_floor_ts - 1)
	//	}
	//
	//	this.deviation_ts = (int)(this.accept_floor_ts - float64(current_floor_ts))
	//
	//	// dup/jmp detect for ts in floor mode.
     //   if int64(this.previous_floor_ts) != 0 && int64(this.previous_floor_ts) != current_floor_ts - 1 {
	//
     //   }
     //   this.previous_floor_ts = float64(current_floor_ts);
	//	// we always ensure the piece is increase one by one.
	//	//todo ts file name replace
	//}
	////todo tsfile append seq suffix
	ts_file := strconv.FormatInt(int64(((time.Now().UnixNano() / 1e6) / (1000 * 5))), 10) + ".ts"
	this.current.full_path = this.hls_path + "/" + ts_file
	//add prefix
	this.current.uri = this.hls_entry_prefix + "/" + ts_file
	// open temp ts file.
	tmp_file := this.current.full_path + ".tmp";
	if err := this.current.Open(tmp_file, default_acodec, default_vcodec); err != nil {
		return err
	}

	if default_acodec != codec.SrsCodecAudioReserved1 {
		this.current.muxer.UpdateACodec(default_acodec)
	}
	_ = tmp_file
	//todo	
	// if err := this.current.muxer.open(tmp_file); err != nil {
	// 	return err
	// }
	
	return nil
}

func (this *SrsHlsMuxer) segment_close() error {
	if this.current == nil {
		return nil
	}
	// valid, add to segments if segment duration is ok
	// when too small, it maybe not enough data to play.
	// when too large, it maybe timestamp corrupt.
	// make the segment more acceptable, when in [min, max_td * 2], it's ok.
	if this.current.duration * 1000 >= 100 && this.current.duration <= float64(this.max_td*2){
		this.segments = append(this.segments, this.current)

		full_path := this.current.full_path
		this.current = nil

		tmp_file := full_path + ".tmp"
		if err := os.Rename(tmp_file, full_path); err != nil {
			return err
		}
	} else {
		this._sequence_no--
		tmp_file := this.current.full_path + ".tmp"
		if err := os.Remove(tmp_file); err != nil {
			return err
		}
	}
	//这里主要限制hls的总时长，超过hls_window的话，就把前面部分文件删除掉，这一大堆代码就是干这个事情
	//shrink the segments.
	var duration float64 = 0
	var removeIndex = 0
	for i := len(this.segments) - 1; i >= 0; i-- {
		duration += this.segments[i].duration
		if duration > this.hls_window {
			removeIndex = i
			break
		}
	}

	segment_to_remove := make([]*SrsHlsSegment, 0)
	for i := 0; i < removeIndex && i < len(this.segments); i++ {
		segment_to_remove = append(segment_to_remove, this.segments[0])
		this.segments = this.segments[1:]
	}

	for i := 0; i < len(segment_to_remove); i++ {
		if this.hls_cleanup {
			if err := os.Remove(segment_to_remove[i].full_path); err != nil {
				return err
			}
		}
	}
	return nil
}

