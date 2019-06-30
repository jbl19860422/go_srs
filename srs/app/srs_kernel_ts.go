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

const SRS_TS_PACKET_SIZE = 188
const SRS_CONSTS_HLS_PURE_AUDIO_AGGREGATE = 720 * 90

/**
* the adaption_field_control of ts packet,
* Table 2-5 - Adaptation field control values, hls-mpeg-ts-iso13818-1.pdf, page 38
 */
type SrsTsAdaptationFieldType int

const (
	_ SrsTsAdaptationFieldType = iota
	// Reserved for future use by ISO/IEC
	SrsTsAdaptationFieldTypeReserved = 0x00
	// No adaptation_field, payload only
	SrsTsAdaptationFieldTypePayloadOnly = 0x01
	// Adaptation_field only, no payload
	SrsTsAdaptationFieldTypeAdaptionOnly = 0x02
	// Adaptation_field followed by payload
	SrsTsAdaptationFieldTypeBoth = 0x03
)

/**
* the actually parsed ts pid,
* @see SrsTsPid, some pid, for example, PMT/Video/Audio is specified by PAT or other tables.
 */
type SrsTsPidApply int

const (
	_                     SrsTsPidApply = iota
	SrsTsPidApplyReserved               = 0 // TSPidTypeReserved, nothing parsed, used reserved.

	SrsTsPidApplyPAT = 1 // Program associtate table
	SrsTsPidApplyPMT = 2 // Program map table.

	SrsTsPidApplyVideo = 3 // for video
	SrsTsPidApplyAudio = 4 // vor audio
)

/**
* Table 2-29 - Stream type assignments
 */
type SrsTsStream int

const (
	_ SrsTsStream = iota
	// ITU-T | ISO/IEC Reserved
	SrsTsStreamReserved = 0x00
	// ISO/IEC 11172 Video
	// ITU-T Rec. H.262 | ISO/IEC 13818-2 Video or ISO/IEC 11172-2 constrained parameter video stream
	// ISO/IEC 11172 Audio
	// ISO/IEC 13818-3 Audio
	SrsTsStreamAudioMp3 = 0x04
	// ITU-T Rec. H.222.0 | ISO/IEC 13818-1 private_sections
	// ITU-T Rec. H.222.0 | ISO/IEC 13818-1 PES packets containing private data
	// ISO/IEC 13522 MHEG
	// ITU-T Rec. H.222.0 | ISO/IEC 13818-1 Annex A DSM-CC
	// ITU-T Rec. H.222.1
	// ISO/IEC 13818-6 type A
	// ISO/IEC 13818-6 type B
	// ISO/IEC 13818-6 type C
	// ISO/IEC 13818-6 type D
	// ITU-T Rec. H.222.0 | ISO/IEC 13818-1 auxiliary
	// ISO/IEC 13818-7 Audio with ADTS transport syntax
	SrsTsStreamAudioAAC = 0x0f
	// ISO/IEC 14496-2 Visual
	SrsTsStreamVideoMpeg4 = 0x10
	// ISO/IEC 14496-3 Audio with the LATM transport syntax as defined in ISO/IEC 14496-3 / AMD 1
	SrsTsStreamAudioMpeg4 = 0x11
	// ISO/IEC 14496-1 SL-packetized stream or FlexMux stream carried in PES packets
	// ISO/IEC 14496-1 SL-packetized stream or FlexMux stream carried in ISO/IEC14496_sections.
	// ISO/IEC 13818-6 Synchronized Download Protocol
	// ITU-T Rec. H.222.0 | ISO/IEC 13818-1 Reserved
	// 0x15-0x7F
	SrsTsStreamVideoH264 = 0x1b
	// User Private
	// 0x80-0xFF
	SrsTsStreamAudioAC3 = 0x81
	SrsTsStreamAudioDTS = 0x8a
)

func srs_ts_stream2string(stream SrsTsStream) string {
	switch stream {
	case SrsTsStreamReserved:
		return "Reserved"
	case SrsTsStreamAudioMp3:
		return "MP3"
	case SrsTsStreamAudioAAC:
		return "AAC"
	case SrsTsStreamAudioAC3:
		return "AC3"
	case SrsTsStreamAudioDTS:
		return "AudioDTS"
	case SrsTsStreamVideoH264:
		return "H.264"
	case SrsTsStreamVideoMpeg4:
		return "MP4"
	case SrsTsStreamAudioMpeg4:
		return "MP4A"
	default:
		return "Other"
	}
}

/**
* the stream_id of PES payload of ts packet.
* Table 2-18 - Stream_id assignments, hls-mpeg-ts-iso13818-1.pdf, page 52.
 */
type SrsTsPESStreamId int

const (
	_ SrsTsPESStreamId = iota
	// program_stream_map
	SrsTsPESStreamIdProgramStreamMap = 0xbc // 0b10111100
	// private_stream_1
	SrsTsPESStreamIdPrivateStream1 = 0xbd // 0b10111101
	// padding_stream
	SrsTsPESStreamIdPaddingStream = 0xbe // 0b10111110
	// private_stream_2
	SrsTsPESStreamIdPrivateStream2 = 0xbf // 0b10111111

	// 110x xxxx
	// ISO/IEC 13818-3 or ISO/IEC 11172-3 or ISO/IEC 13818-7 or ISO/IEC
	// 14496-3 audio stream number x xxxx
	// ((sid >> 5) & 0x07) == SrsTsPESStreamIdAudio
	// @remark, use SrsTsPESStreamIdAudioCommon as actually audio, SrsTsPESStreamIdAudio to check whether audio.
	SrsTsPESStreamIdAudioChecker = 0x06 // 0b110
	SrsTsPESStreamIdAudioCommon  = 0xc0

	// 1110 xxxx
	// ITU-T Rec. H.262 | ISO/IEC 13818-2 or ISO/IEC 11172-2 or ISO/IEC
	// 14496-2 video stream number xxxx
	// ((stream_id >> 4) & 0x0f) == SrsTsPESStreamIdVideo
	// @remark, use SrsTsPESStreamIdVideoCommon as actually video, SrsTsPESStreamIdVideo to check whether video.
	SrsTsPESStreamIdVideoChecker = 0x0e // 0b1110
	SrsTsPESStreamIdVideoCommon  = 0xe0

	// ECM_stream
	SrsTsPESStreamIdEcmStream = 0xf0 // 0b11110000
	// EMM_stream
	SrsTsPESStreamIdEmmStream = 0xf1 // 0b11110001
	// DSMCC_stream
	SrsTsPESStreamIdDsmccStream = 0xf2 // 0b11110010
	// 13522_stream
	SrsTsPESStreamId13522Stream = 0xf3 // 0b11110011
	// H_222_1_type_A
	SrsTsPESStreamIdH2221TypeA = 0xf4 // 0b11110100
	// H_222_1_type_B
	SrsTsPESStreamIdH2221TypeB = 0xf5 // 0b11110101
	// H_222_1_type_C
	SrsTsPESStreamIdH2221TypeC = 0xf6 // 0b11110110
	// H_222_1_type_D
	SrsTsPESStreamIdH2221TypeD = 0xf7 // 0b11110111
	// H_222_1_type_E
	SrsTsPESStreamIdH2221TypeE = 0xf8 // 0b11111000
	// ancillary_stream
	SrsTsPESStreamIdAncillaryStream = 0xf9 // 0b11111001
	// SL_packetized_stream
	SrsTsPESStreamIdSlPacketizedStream = 0xfa // 0b11111010
	// FlexMux_stream
	SrsTsPESStreamIdFlexMuxStream = 0xfb // 0b11111011
	// reserved data stream
	// 1111 1100 ... 1111 1110
	// program_stream_directory
	SrsTsPESStreamIdProgramStreamDirectory = 0xff // 0b11111111
)
