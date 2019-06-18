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