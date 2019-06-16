package amf0

import (
	"errors"
	"go_srs/srs/utils"
)

type SrsAmf0Undefined struct {
}

func (this *SrsAmf0Undefined) Decode(stream *utils.SrsStream) error {
	marker, err := stream.ReadByte()
	if err != nil {
		return err
	}

	if marker != RTMP_AMF0_Undefined {
		err := errors.New("amf0 check null marker failed.")
		return err
	}
	return nil
}

func (this *SrsAmf0Undefined) Encode(stream *utils.SrsStream) error {
	stream.WriteByte(RTMP_AMF0_Undefined)
	return nil
}

func (this *SrsAmf0Undefined) IsMyType(stream *utils.SrsStream) (bool, error) {
	marker, err := stream.PeekByte()
	if err != nil {
		return false, err
	}

	if marker != RTMP_AMF0_Undefined {
		return false, nil
	}
	return true, nil
}

func (this *SrsAmf0Undefined) GetValue() interface{} {
	return nil
}
