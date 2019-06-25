package app

type SrsHlsSegment struct {
	duration           float64            // duration in seconds in m3u8.
	sequence_no        int                // sequence number in m3u8.
	uri                string             // ts uri in m3u8.
	full_path          string             //ts full file to write.
	writer             *SrsHlsCacheWriter //the muxer to write ts.
	muxer              *SrsTsMuxer
	segment_start_dts  int64 // current segment start dts for m3u8
	is_sequence_header bool  // whether current segement is sequence header.
}

const SRS_AUTO_HLS_SEGMENT_TIMESTAMP_JUMP_MS = 300

func NewSrsHlsSegment(c *SrsTsContext, ac SrsCodecAudio, vc SrsCodecVideo) *SrsHlsSegment {
	w := NewSrsHlsCacheWriter()
	return &SrsHlsSegment{
		muxer:  NewSrsMuxer(w, ac, vc),
		writer: w,
	}
}

func (this *SrsHlsSegment) update_duration(current_frame_dts int64) {
	if current_frame_dts < this.segment_start_dts {
		if current_frame_dts < this.segment_start_dts-SRS_AUTO_HLS_SEGMENT_TIMESTAMP_JUMP_MS*90 {
			this.segment_start_dts = current_frame_dts
		}
	}

	this.duration = (current_frame_dts - this.segment_start_dts) / 90000.0
}
