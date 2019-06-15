package amf0

import (
	"encoding/binary"
	"errors"
	"utils"
)

type SrsAmf0String struct {
	value SrsAmf0Utf8
}

func NewSrsAmf0String(str string) *SrsAmf0String {
	return &SrsAmf0String{
		value:SrsAmf0Utf8{value:str}
	}
}

func (this *SrsAmf0String) Decode(stream *utils.SrsStream) error {
	marker, err := stream.ReadByte()
	if err != nil {
		return err
	}

	if marker != RTMP_AMF0_String {
		err := errors.New("amf0 check string marker failed.")
		return err
	}

	err := this.value.Decode(stream)
	return err
}

func (this *SrsAmf0String) Encode(stream *utils.SrsStream) error {
	stream.WriteByte(RTMP_AMF0_String)
	this.value.Encode(stream)
	return nil
}

func (this *SrsAmf0String) IsMyType(stream *utils.SrsStream) (bool, error) {
	marker, err := stream.PeekByte()
	if err != nil {
		return err
	}

	if marker != RTMP_AMF0_String {
		return false, nil
	}
	return true, nil
}
