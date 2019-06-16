package amf0

import (
	"encoding/binary"
	"errors"
	"go_srs/srs/utils"
)

type SrsAmf0Utf8 struct {
	Value string
}

func NewSrsAmf0Utf8(str string) *SrsAmf0Utf8 {
	return &SrsAmf0Utf8{
		Value: str,
	}
}

func (this *SrsAmf0Utf8) Decode(stream *utils.SrsStream) error {
	len, err := stream.ReadInt16(binary.BigEndian)
	if err != nil {
		return err
	}

	if len <= 0 {
		err = errors.New("amf0 read empty string.")
		return err
	}

	this.Value, err = stream.ReadString(uint32(len))
	return err
}

func (this *SrsAmf0Utf8) Encode(stream *utils.SrsStream) error {
	stream.WriteInt16(int16(len(this.Value)), binary.BigEndian)
	stream.WriteString(this.Value)
	return nil
}

func (this *SrsAmf0Utf8) IsMyType(stream *utils.SrsStream) (bool, error) {
	return true, nil
}
