package app

import (
	"os"
	"go_srs/srs/codec"
	"io"
	"fmt"
)

/**
* the wrapper of m3u8 segment from specification:
*
* 3.3.2.  EXTINF
* The EXTINF tag specifies the duration of a media segment.
*/

type SrsHlsSegment struct {
	duration           float64            // duration in seconds in m3u8.
	sequence_no        int                // sequence number in m3u8.
	uri                string             // ts uri in m3u8.
	full_path          string             //ts full file to write.
	writer             io.Writer 		//the muxer to write ts.
	muxer              *SrsTsMuxer
	segment_start_dts  int64 // current segment start dts for m3u8
	is_sequence_header bool  // whether current segement is sequence header.
	context				*SrsTsContext
}

const SRS_AUTO_HLS_SEGMENT_TIMESTAMP_JUMP_MS = 300

func NewSrsHlsSegment(c *SrsTsContext) *SrsHlsSegment {
	return &SrsHlsSegment{
		context:c,
	}
}

func (this *SrsHlsSegment) UpdateDuration(currentFrameDts int64) {
	if currentFrameDts < this.segment_start_dts {
		if currentFrameDts < this.segment_start_dts-SRS_AUTO_HLS_SEGMENT_TIMESTAMP_JUMP_MS*90 {
			this.segment_start_dts = currentFrameDts
		}
	}

	this.duration = float64(currentFrameDts - this.segment_start_dts) / 90000.0
}

func (this *SrsHlsSegment) Open(path string, ac codec.SrsCodecAudio, vc codec.SrsCodecVideo) error {
	var err error
	this.writer, err = os.OpenFile(path, os.O_RDWR | os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("open full path failed, ", this.full_path)
		return err
	}

	fmt.Println("open segment path succeed:", this.full_path)
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