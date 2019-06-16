package amf0

import (
	"errors"
	"go_srs/srs/utils"
)

type SrsAmf0Boolean struct {
	Value bool
}

func NewSrsAmf0Boolean(data bool) *SrsAmf0Boolean {
	return &SrsAmf0Boolean{
		Value: data,
	}
}

func (this *SrsAmf0Boolean) Decode(stream *utils.SrsStream) error {
	marker, err := stream.ReadByte()
	if err != nil {
		return err
	}

	if marker != RTMP_AMF0_Boolean {
		err := errors.New("amf0 check bool marker failed.")
		return err
	}

	this.Value, err = stream.ReadBool()
	if err != nil {
		return err
	}
	return nil
}

func (this *SrsAmf0Boolean) Encode(stream *utils.SrsStream) error {
	stream.WriteByte(RTMP_AMF0_Boolean)
	var d byte
	if this.Value {
		d = 1
	} else {
		d = 0
	}
	stream.WriteByte(d)
	return nil
}

func (this *SrsAmf0Boolean) IsMyType(stream *utils.SrsStream) (bool, error) {
	marker, err := stream.PeekByte()
	if err != nil {
		return false, err
	}

	if marker != RTMP_AMF0_Boolean {
		return false, nil
	}
	return true, nil
}
