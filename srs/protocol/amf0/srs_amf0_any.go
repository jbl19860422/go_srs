package amf0

import (
	"utils"
)

type SrsAmf0Any {
	Decode(stream *utils.SrsStream) error
	Encode(stream *utils.SrsStream) error
	IsMyType(stream *utils.SrsStream) (bool, error)
}

func GenerateSrsAmf0Any(marker byte) SrsAmf0Any {
	switch marker {
	case RTMP_AMF0_Number:
		return &SrsAmf0Number{}
	case RTMP_AMF0_Boolean:
		return &SrsAmf0Boolean{}
	case RTMP_AMF0_String:
		return &SrsAmf0String{}
	case RTMP_AMF0_Object:
		return &SrsAmf0Object{}
	case RTMP_AMF0_Null:
		return &SrsAmf0Null{}
	case RTMP_AMF0_Undefined:
		return &SrsAmf0Undefined{}
	case RTMP_AMF0_EcmaArray:
		return &SrsAmf0EcmaArray{}
	default:
		return nil
	}
}