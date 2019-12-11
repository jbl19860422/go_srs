/*
The MIT License (MIT)

Copyright (c) 2019 GOSRS(gosrs)

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
	"errors"
	"go_srs/srs/codec"
)

type SrsTsCache struct {
	audio *SrsTsMessage
	video *SrsTsMessage
}

func NewSrsTsCache() *SrsTsCache {
	return &SrsTsCache{}
}

func (this *SrsTsCache) cache_audio(c *SrsAvcAacCodec, dts int64, sampler *SrsCodecSampler) error {
	if this.audio == nil {
		this.audio = NewSrsTsMessage()
		this.audio.writePcr = false
		this.audio.dts = dts
		this.audio.pts = dts
		this.audio.startPts = dts
	}

	this.audio.sid = SrsTsPESStreamIdAudioCommon //used in ts stream_id field
	acodec := codec.SrsCodecAudio(c.audioCodecId)
	if acodec == codec.SrsCodecAudioAAC {
		if err := this.do_cache_aac(c, sampler); err != nil {
			return err
		}
	}
	return nil
}

func (this *SrsTsCache) cache_video(c *SrsAvcAacCodec, dts int64, sampler *SrsCodecSampler) error {
	if this.video == nil {
		this.video = NewSrsTsMessage()
		if sampler.FrameType == codec.SrsCodecVideoAVCFrameKeyFrame {
			this.video.writePcr = true
		} else {
			this.video.writePcr = false
		}
		this.video.startPts = dts
	}

	this.video.dts = dts
	this.video.pts = this.video.dts + int64(sampler.Cts)*90
	this.video.sid = SrsTsPESStreamIdVideoCommon //this is the hint to judge the SrsTsMessage is audio or video
	if err := this.do_cache_avc(c, sampler); err != nil {
		return err
	}
	return nil
}

func (this *SrsTsCache) do_cache_avc(c *SrsAvcAacCodec, sampler *SrsCodecSampler) error {
	audInserted := false
	this.video.payload = make([]byte, 0)
	if sampler.HasAud {
		// the aud(access unit delimiter) before each frame.
		// 7.3.2.4 Access unit delimiter RBSP syntax
		// H.264-AVC-ISO_IEC_14496-10-2012.pdf, page 66.
		//
		// primary_pic_type u(3), the first 3bits, primary_pic_type indicates that the slice_type values
		//      for all slices of the primary coded picture are members of the set listed in Table 7-5 for
		//      the given value of primary_pic_type.
		//      0, slice_type 2, 7
		//      1, slice_type 0, 2, 5, 7
		//      2, slice_type 0, 1, 2, 5, 6, 7
		//      3, slice_type 4, 9
		//      4, slice_type 3, 4, 8, 9
		//      5, slice_type 2, 4, 7, 9
		//      6, slice_type 0, 2, 3, 4, 5, 7, 8, 9
		//      7, slice_type 0, 1, 2, 3, 4, 5, 6, 7, 8, 9
		// 7.4.2.4 Access unit delimiter RBSP semantics
		// H.264-AVC-ISO_IEC_14496-10-2012.pdf, page 102.
		//
		// slice_type specifies the coding type of the slice according to Table 7-6.
		//      0, P (P slice)
		//      1, B (B slice)
		//      2, I (I slice)
		//      3, SP (SP slice)
		//      4, SI (SI slice)
		//      5, P (P slice)
		//      6, B (B slice)
		//      7, I (I slice)
		//      8, SP (SP slice)
		//      9, SI (SI slice)
		// H.264-AVC-ISO_IEC_14496-10-2012.pdf, page 105.

		defaultAudNalu := []byte{0x09, 0x0f}
		this.video.payload = append(this.video.payload, []byte{0x00, 0x00, 0x00, 0x01}...)
		this.video.payload = append(this.video.payload, defaultAudNalu...)
		audInserted = true
	}

	//When it is the first byte stream NAL unit in the bitstream, it may
	//also contain one or more additional leading_zero_8bits syntax elements.
	isSpsPpsAppend := false
	for i := 0; i < len(sampler.SampleUnits); i++ {
		if sampler.SampleUnits[i] == nil || len(sampler.SampleUnits[i]) <= 0 {
			return errors.New("sample unit must not be nil or empty")
		}

		naluUnitType := codec.SrsAvcNaluType(sampler.SampleUnits[i][0] & 0x0f)
		// Insert sps/pps before IDR when there is no sps/pps in samples.
		// The sps/pps is parsed from sequence header(generally the first flv packet).
		if naluUnitType == codec.SrsAvcNaluTypeIDR && !sampler.HasSpsPps && !isSpsPpsAppend {
			if len(c.sequenceParameterSetNALUnit) > 0 {
				if audInserted {
					this.video.payload = append(this.video.payload, []byte{0x00, 0x00, 0x01}...)
				} else {
					this.video.payload = append(this.video.payload, []byte{0x00, 0x00, 0x00, 0x01}...)
					audInserted = true
				}
				this.video.payload = append(this.video.payload, c.sequenceParameterSetNALUnit...)
			}

			if len(c.pictureParameterSetNALUnit) > 0 {
				if audInserted {
					this.video.payload = append(this.video.payload, []byte{0x00, 0x00, 0x01}...)
				} else {
					this.video.payload = append(this.video.payload, []byte{0x00, 0x00, 0x00, 0x01}...)
					audInserted = true
				}
				this.video.payload = append(this.video.payload, c.pictureParameterSetNALUnit...)
			}
			isSpsPpsAppend = true
		}
		this.video.payload = append(this.video.payload, []byte{0x00, 0x00, 0x01}...)
		this.video.payload = append(this.video.payload, sampler.SampleUnits[i]...)
	}
	return nil
}

func (this *SrsTsCache) do_cache_aac(c *SrsAvcAacCodec, sampler *SrsCodecSampler) error {
	this.audio.payload = make([]byte, 0)
	for i := 0; i < len(sampler.SampleUnits); i++ {
		var frame_length int32 = int32(len(sampler.SampleUnits[i]) + 7)
		// AAC-ADTS
		// 6.2 Audio Data Transport Stream, ADTS
		// in aac-iso-13818-7.pdf, page 26.
		// fixed 7bytes header
		var adtsHeader = []byte{0xff, 0xf9, 0x00, 0x00, 0x00, 0x0f, 0xfc}
		// profile, 2bits
		aac_profile := codec.SrsCodecAacRtmp2Ts(c.aacObject)
		adtsHeader[2] = (byte(aac_profile) << 6) & 0xc0
		// sampling_frequency_index 4bits
		adtsHeader[2] |= byte((c.aacSampleRateIndex << 2) & 0x3c)
		// channel_configuration 3bits
		adtsHeader[2] |= byte((c.aacChannels >> 2) & 0x01)
		adtsHeader[3] = byte((int32(c.aacChannels) << 6) & 0xc0)
		// frame_length 13bits
		adtsHeader[3] |= byte((frame_length >> 11) & 0x03)
		adtsHeader[4] = byte((frame_length >> 3) & 0xff)
		adtsHeader[5] = byte(((frame_length << 5) & 0xe0))
		// adts_buffer_fullness; //11bits
		adtsHeader[5] |= 0x1f

		// copy to audio buffer
		this.audio.payload = append(this.audio.payload, adtsHeader...)
		this.audio.payload = append(this.audio.payload, sampler.SampleUnits[i]...)
	}

	return nil
}

func (this *SrsTsCache) do_cache_mp3(c *SrsAvcAacCodec, sampler *SrsCodecSampler) error {
	// for mp3, directly write to cache.
	// TODO: FIXME: implements the ts jitter.
	p := make([]byte, 0)
	for i := 0; i < len(sampler.SampleUnits); i++ {
		p = append(p, sampler.SampleUnits[i]...)
	}

	return nil
}
