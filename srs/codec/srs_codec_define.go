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
package codec

type SrsCodecVideoAVCFrame int

const (
	_ SrsCodecVideoAVCFrame = iota
	SrsCodecVideoAVCFrameReserved = 0
	SrsCodecVideoAVCFrameReserved1 = 6
	SrsCodecVideoAVCFrameKeyFrame = 1
	SrsCodecVideoAVCFrameInterFrame = 2
	SrsCodecVideoAVCFrameDisposableInterFrame = 3
	SrsCodecVideoAVCFrameGeneratedKeyFrame = 4
	SrsCodecVideoAVCFrameVideoInfoFrame = 5
)

/**
* the aac profile, for ADTS(HLS/TS)
* @see https://github.com/ossrs/srs/issues/310
*/
type SrsAacProfile int
const (
    _ SrsAacProfile = iota
    SrsAacProfileReserved = 3
    
    // @see 7.1 Profiles, aac-iso-13818-7.pdf, page 40
    SrsAacProfileMain = 0
    SrsAacProfileLC = 1
    SrsAacProfileSSR = 2
)

// E.4.3.1 VIDEODATA
// CodecID UB [4]
// Codec Identifier. The following values are defined:
//     2 = Sorenson H.263
//     3 = Screen video
//     4 = On2 VP6
//     5 = On2 VP6 with alpha channel
//     6 = Screen video version 2
//     7 = AVC
type SrsCodecVideo int
const (
	_ SrsCodecVideo = iota
	// set to the zero to reserved, for array map.
    SrsCodecVideoReserved                = 0
    SrsCodecVideoReserved1                = 1
	SrsCodecVideoReserved2                = 9
	
	// for user to disable video, for example, use pure audio hls.
    SrsCodecVideoDisabled                = 8
    
    SrsCodecVideoSorensonH263             = 2
    SrsCodecVideoScreenVideo             = 3
    SrsCodecVideoOn2VP6                 = 4
    SrsCodecVideoOn2VP6WithAlphaChannel = 5
    SrsCodecVideoScreenVideoVersion2     = 6
    SrsCodecVideoAVC                     = 7
)

// SoundFormat UB [4] 
// Format of SoundData. The following values are defined:
//     0 = Linear PCM, platform endian
//     1 = ADPCM
//     2 = MP3
//     3 = Linear PCM, little endian
//     4 = Nellymoser 16 kHz mono
//     5 = Nellymoser 8 kHz mono
//     6 = Nellymoser
//     7 = G.711 A-law logarithmic PCM
//     8 = G.711 mu-law logarithmic PCM
//     9 = reserved
//     10 = AAC
//     11 = Speex
//     14 = MP3 8 kHz
//     15 = Device-specific sound
// Formats 7, 8, 14, and 15 are reserved.
// AAC is supported in Flash Player 9,0,115,0 and higher.
// Speex is supported in Flash Player 10 and higher.
type SrsCodecAudio int
const (
	_ SrsCodecAudio = iota
	// set to the max value to reserved, for array map.
    SrsCodecAudioReserved1                = 16
    
    // for user to disable audio, for example, use pure video hls.
    SrsCodecAudioDisabled                   = 17
    
    SrsCodecAudioLinearPCMPlatformEndian             = 0
    SrsCodecAudioADPCM                                 = 1
    SrsCodecAudioMP3                                 = 2
    SrsCodecAudioLinearPCMLittleEndian                 = 3
    SrsCodecAudioNellymoser16kHzMono                 = 4
    SrsCodecAudioNellymoser8kHzMono                 = 5
    SrsCodecAudioNellymoser                         = 6
    SrsCodecAudioReservedG711AlawLogarithmicPCM        = 7
    SrsCodecAudioReservedG711MuLawLogarithmicPCM    = 8
    SrsCodecAudioReserved                             = 9
    SrsCodecAudioAAC                                 = 10
    SrsCodecAudioSpeex                                 = 11
    SrsCodecAudioReservedMP3_8kHz                     = 14
    SrsCodecAudioReservedDeviceSpecificSound         = 15
)

// AVCPacketType IF CodecID == 7 UI8
// The following values are defined:
//     0 = AVC sequence header
//     1 = AVC NALU
//     2 = AVC end of sequence (lower level NALU sequence ender is
//         not required or supported)
type SrsCodecVideoAVCType int
const (
	_ SrsCodecVideoAVCType = iota
	SrsCodecVideoAVCTypeReserved                     = 3
    
    SrsCodecVideoAVCTypeSequenceHeader               = 0
    SrsCodecVideoAVCTypeNALU                         = 1
    SrsCodecVideoAVCTypeSequenceHeaderEOF            = 2
)

// AACPacketType IF SoundFormat == 10 UI8
// The following values are defined:
//     0 = AAC sequence header
//     1 = AAC raw
type SrsCodecAudioType int
const (
	_ SrsCodecAudioType = iota
	// set to the max value to reserved, for array map.
    SrsCodecAudioTypeReserved                        = 2
    
    SrsCodecAudioTypeSequenceHeader                  = 0
    SrsCodecAudioTypeRawData                         = 1
)

/**
 * Table 7-1 - NAL unit type codes, syntax element categories, and NAL unit type classes
 * H.264-AVC-ISO_IEC_14496-10-2012.pdf, page 83.
 */
 type SrsAvcNaluType int
 const (
     _ SrsAvcNaluType = iota
     // Unspecified
     SrsAvcNaluTypeReserved = 0
     
     // Coded slice of a non-IDR picture slice_layer_without_partitioning_rbsp( )
     SrsAvcNaluTypeNonIDR = 1
     // Coded slice data partition A slice_data_partition_a_layer_rbsp( )
     SrsAvcNaluTypeDataPartitionA = 2
     // Coded slice data partition B slice_data_partition_b_layer_rbsp( )
     SrsAvcNaluTypeDataPartitionB = 3
     // Coded slice data partition C slice_data_partition_c_layer_rbsp( )
     SrsAvcNaluTypeDataPartitionC = 4
     // Coded slice of an IDR picture slice_layer_without_partitioning_rbsp( )
     SrsAvcNaluTypeIDR = 5
     // Supplemental enhancement information (SEI) sei_rbsp( )
     SrsAvcNaluTypeSEI = 6
     // Sequence parameter set seq_parameter_set_rbsp( )
     SrsAvcNaluTypeSPS = 7
     // Picture parameter set pic_parameter_set_rbsp( )
     SrsAvcNaluTypePPS = 8
     // Access unit delimiter access_unit_delimiter_rbsp( )
     SrsAvcNaluTypeAccessUnitDelimiter = 9
     // End of sequence end_of_seq_rbsp( )
     SrsAvcNaluTypeEOSequence = 10
     // End of stream end_of_stream_rbsp( )
     SrsAvcNaluTypeEOStream = 11
     // Filler data filler_data_rbsp( )
     SrsAvcNaluTypeFilterData = 12
     // Sequence parameter set extension seq_parameter_set_extension_rbsp( )
     SrsAvcNaluTypeSPSExt = 13
     // Prefix NAL unit prefix_nal_unit_rbsp( )
     SrsAvcNaluTypePrefixNALU = 14
     // Subset sequence parameter set subset_seq_parameter_set_rbsp( )
     SrsAvcNaluTypeSubsetSPS = 15
     // Coded slice of an auxiliary coded picture without partitioning slice_layer_without_partitioning_rbsp( )
     SrsAvcNaluTypeLayerWithoutPartition = 19
     // Coded slice extension slice_layer_extension_rbsp( )
     SrsAvcNaluTypeCodedSliceExt = 20
 )

/**
* the FLV/RTMP supported audio sample rate.
* Sampling rate. The following values are defined:
* 0 = 5.5 kHz = 5512 Hz
* 1 = 11 kHz = 11025 Hz
* 2 = 22 kHz = 22050 Hz
* 3 = 44 kHz = 44100 Hz
*/
type SrsCodecAudioSampleRate int
const (
    _ SrsCodecAudioSampleRate = iota
    // set to the max value to reserved, for array map.
    SrsCodecAudioSampleRateReserved                 = 4
    SrsCodecAudioSampleRate5512                     = 0
    SrsCodecAudioSampleRate11025                    = 1
    SrsCodecAudioSampleRate22050                    = 2
    SrsCodecAudioSampleRate44100                    = 3
)

/**
* the FLV/RTMP supported audio sample size.
* Size of each audio sample. This parameter only pertains to
* uncompressed formats. Compressed formats always decode
* to 16 bits internally.
* 0 = 8-bit samples
* 1 = 16-bit samples
*/
type SrsCodecAudioSampleSize int
const (
    _ SrsCodecAudioSampleSize = iota
    // set to the max value to reserved, for array map.
    SrsCodecAudioSampleSizeReserved                 = 2
    
    SrsCodecAudioSampleSize8bit                     = 0
    SrsCodecAudioSampleSize16bit                    = 1
)


/**
* the FLV/RTMP supported audio sound type/channel.
* Mono or stereo sound
* 0 = Mono sound
* 1 = Stereo sound
*/
type  SrsCodecAudioSoundType int
const (
    _ SrsCodecAudioSoundType = iota
    // set to the max value to reserved, for array map.
    SrsCodecAudioSoundTypeReserved                  = 2
    
    SrsCodecAudioSoundTypeMono                      = 0
    SrsCodecAudioSoundTypeStereo                    = 1
)
 
const SRS_AAC_SAMPLE_RATE_UNSET = 15

/**
* the profile for avc/h.264.
* @see Annex A Profiles and levels, H.264-AVC-ISO_IEC_14496-10.pdf, page 205.
*/
type SrsAvcProfile int
const (
	_ SrsAvcProfile = iota
	SrsAvcProfileReserved = 0
    
    // @see ffmpeg, libavcodec/avcodec.h:2713
    SrsAvcProfileBaseline = 66
    // FF_PROFILE_H264_CONSTRAINED  (1<<9)  // 8+1; constraint_set1_flag
    // FF_PROFILE_H264_CONSTRAINED_BASELINE (66|FF_PROFILE_H264_CONSTRAINED)
    SrsAvcProfileConstrainedBaseline = 578
    SrsAvcProfileMain = 77
    SrsAvcProfileExtended = 88
    SrsAvcProfileHigh = 100
    SrsAvcProfileHigh10 = 110
    SrsAvcProfileHigh10Intra = 2158
    SrsAvcProfileHigh422 = 122
    SrsAvcProfileHigh422Intra = 2170
    SrsAvcProfileHigh444 = 144
    SrsAvcProfileHigh444Predictive = 244
    SrsAvcProfileHigh444Intra = 2192
)

/**
* the level for avc/h.264.
* @see Annex A Profiles and levels, H.264-AVC-ISO_IEC_14496-10.pdf, page 207.
*/
type SrsAvcLevel int
const (
	_ SrsAvcLevel = iota
	SrsAvcLevelReserved = 0
    
    SrsAvcLevel_1 = 10
    SrsAvcLevel_11 = 11
    SrsAvcLevel_12 = 12
    SrsAvcLevel_13 = 13
    SrsAvcLevel_2 = 20
    SrsAvcLevel_21 = 21
    SrsAvcLevel_22 = 22
    SrsAvcLevel_3 = 30
    SrsAvcLevel_31 = 31
    SrsAvcLevel_32 = 32
    SrsAvcLevel_4 = 40
    SrsAvcLevel_41 = 41
    SrsAvcLevel_5 = 50
    SrsAvcLevel_51 = 51
)

/**
* the avc payload format, must be ibmf or annexb format.
* we guess by annexb first, then ibmf for the first time,
* and we always use the guessed format for the next time.
*/
type SrsAvcPayloadFormat int
const (
	_ SrsAvcPayloadFormat = iota
	SrsAvcPayloadFormatGuess = 0
    SrsAvcPayloadFormatAnnexb = 1
    SrsAvcPayloadFormatIbmf = 2
)
/**
* the aac object type, for RTMP sequence header
* for AudioSpecificConfig, @see aac-mp4a-format-ISO_IEC_14496-3+2001.pdf, page 33
* for audioObjectType, @see aac-mp4a-format-ISO_IEC_14496-3+2001.pdf, page 23
*/
type SrsAacObjectType int
const (
	_ SrsAacObjectType = iota
    SrsAacObjectTypeReserved = 0
    
    // Table 1.1 - Audio Object Type definition
    // @see @see aac-mp4a-format-ISO_IEC_14496-3+2001.pdf, page 23
    SrsAacObjectTypeAacMain = 1
    SrsAacObjectTypeAacLC = 2
    SrsAacObjectTypeAacSSR = 3
    
    // AAC HE = LC+SBR
    SrsAacObjectTypeAacHE = 5
    // AAC HEv2 = LC+SBR+PS
    SrsAacObjectTypeAacHEV2 = 29
)
