package amf0

import (
	"encoding/binary"
	"errors"
	"go_srs/srs/utils"
)

type SrsAmf0ObjectEOF struct {
}

func NewSrsAmf0ObjectEOF() *SrsAmf0ObjectEOF {
	return &SrsAmf0ObjectEOF{}
}

func (this *SrsAmf0ObjectEOF) Decode(stream *utils.SrsStream) error {
	tmp, err := stream.ReadInt16(binary.BigEndian)
	if err != nil {
		return err
	}

	if tmp != 0x00 {
		err = errors.New("amf0 read object eof value check failed.")
		return err
	}

	marker, err := stream.ReadByte()
	if err != nil {
		return err
	}

	if marker != RTMP_AMF0_ObjectEnd {
		err := errors.New("amf0 check string marker failed.")
		return err
	}
	return nil
}

func (this *SrsAmf0ObjectEOF) Encode(stream *utils.SrsStream) error {
	stream.WriteInt16(0, binary.BigEndian)
	stream.WriteByte(RTMP_AMF0_ObjectEnd)
	return nil
}

func (this *SrsAmf0ObjectEOF) IsMyType(stream *utils.SrsStream) (bool, error) {
	b, err := stream.PeekBytes(3)
	if err != nil {
		return false, err
	}

	if b[0] != 0x00 || b[1] != 0x00 || b[2] != 0x09 {
		return false, nil
	}

	return true, nil
}
