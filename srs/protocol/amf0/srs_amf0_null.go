package amf0

import (
	"go_srs/srs/utils"
	"errors"
	_ "log"
)

type SrsAmf0Null struct {
}

func NewSrsAmf0Null() *SrsAmf0Null {
	return &SrsAmf0Null{}
}

func (this *SrsAmf0Null) Decode(stream *utils.SrsStream) error {
	marker, err := stream.ReadByte()
	if err != nil {
		return err
	}

	if marker != RTMP_AMF0_Null {
		err := errors.New("amf0 check null marker failed.")
		return err
	}
	return nil
}

func (this *SrsAmf0Null) Encode(stream *utils.SrsStream) error {
	stream.WriteByte(RTMP_AMF0_Null)
	return nil
}

func (this *SrsAmf0Null) IsMyType(stream *utils.SrsStream) (bool, error) {
	marker, err := stream.PeekByte()
	if err != nil {
		return false, err
	}

	if marker != RTMP_AMF0_Null {
		return false, nil
	}
	return true, nil
}

func (this *SrsAmf0Null) GetValue() interface{} {
	return nil
}
