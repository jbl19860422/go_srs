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
	log "github.com/sirupsen/logrus"
	"go_srs/srs/codec"
	"io"
	"os"
)

/**
* the wrapper of m3u8 segment from specification:
*
* 3.3.2.  EXTINF
* The EXTINF tag specifies the duration of a media segment.
 */

type SrsHlsSegment struct {
	duration           float64   // duration in seconds in m3u8.
	sequence_no        int       // sequence number in m3u8.
	uri                string    // ts uri in m3u8.
	full_path          string    //ts full file to write.
	writer             io.Writer //the muxer to write ts.
	muxer              *SrsTsMuxer
	segment_start_dts  int64 // current segment start dts for m3u8
	is_sequence_header bool  // whether current segement is sequence header.
	context            *SrsTsContext
}

const SRS_AUTO_HLS_SEGMENT_TIMESTAMP_JUMP_MS = 300

func NewSrsHlsSegment(c *SrsTsContext) *SrsHlsSegment {
	return &SrsHlsSegment{
		context: c,
	}
}

func (this *SrsHlsSegment) UpdateDuration(currentFrameDts int64) {
	if currentFrameDts < this.segment_start_dts {
		if currentFrameDts < this.segment_start_dts-SRS_AUTO_HLS_SEGMENT_TIMESTAMP_JUMP_MS*90 {
			this.segment_start_dts = currentFrameDts
		}
	}

	this.duration = float64(currentFrameDts-this.segment_start_dts) / 90000.0
}

func (this *SrsHlsSegment) Open(path string, ac codec.SrsCodecAudio, vc codec.SrsCodecVideo) error {
	var err error
	this.writer, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Error("open full path failed, ", this.full_path)
		return err
	}

	log.Info("open segment path succeed:", this.full_path)
	this.muxer = NewSrsTsMuxer(this.writer, this.context, ac, vc)
	return nil
}

func (this *SrsHlsSegment) Close() error {
	if this.writer != nil {
		this.writer.(*os.File).Close()
	}
	return nil
}

func (this *SrsHlsSegment) WriteAudio(audio *SrsTsMessage) error {
	return this.muxer.WriteAudio(audio)
}

func (this *SrsHlsSegment) WriteVideo(video *SrsTsMessage) error {
	return this.muxer.WriteVideo(video)
}
