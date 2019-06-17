package flvcodec

import (
	"go_srs/srs/codec"
)

func VideoIsKeyFrame(data []byte) bool {
	if len(data) < 1 {
		return false
	}

	frameType := (data[0] >> 4) & 0x0F
	return frameType == codec.SrsCodecVideoAVCFrameKeyFrame
}

func VideoIsSequenceHeader(data []byte) bool {
	if !VideoIsH264(data) {
		return false
	}

	if len(data) < 2 {
		return false
	}

	formatType := (data[0] >> 4) & 0x0F
	avcPacketType := data[1]
	return formatType == codec.SrsCodecVideoAVCFrameKeyFrame && avcPacketType == codec.SrsCodecVideoAVCTypeSequenceHeader
}

func AudioIsSequenceHeader(data []byte) bool {
	if !AudioIsAAC(data) {
		return false
	}

	if len(data) < 2 {
		return false
	}

	aacPacketType := data[1]
	return aacPacketType == codec.SrsCodecAudioTypeSequenceHeader
}

func VideoIsH264(data []byte) bool {
	if len(data) < 1 {
		return false
	}

	codecId := data[0] & 0x0F
	return codecId == codec.SrsCodecVideoAVC
}

func AudioIsAAC(data []byte) bool {
	if len(data) < 1 {
		return false
	}

	soundFormat := (data[0] >> 4) & 0x0F
	return soundFormat == codec.SrsCodecAudioAAC
}

func VideoIsAcceptable(data []byte) bool {
	if len(data) < 1 {
		return false
	}

	formatType := data[0]
	codecId := formatType & 0x0F
	formatType = (formatType >> 4) & 0x0F

	if formatType < 1 || formatType > 5 {
		return false
	}

	if codecId < 2 || codecId > 7 {
		return false
	}

	return true
}
