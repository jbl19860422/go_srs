package amf0

import (
	"encoding/binary"
	"errors"
	"utils"
)

type SrsAmf0Utf8 struct {
	value string
}

func NewSrsAmf0Utf8(str string) *SrsAmf0Utf8 {
	return &SrsAmf0Utf8{
		value:str
	}
}

func (this *SrsAmf0Utf8) Decode(stream *utils.SrsStream) error {
	len, err := stream.ReadUInt16(binary.BigEndian)
	if err != nil {
		return err
	}

	if len <= 0 {
		err = errors.New("amf0 read empty string.")
		return err
	}

	this.value, err = stream.ReadString(len)
	return err
}

func (this *SrsAmf0Utf8) Encode(stream *utils.SrsStream) error {
	stream.WriteUInt16(uint16(len(this.value)))
	stream.WriteString(this.value)
	return nil
}

func (this *SrsAmf0Utf8) IsMyType(stream *utils.SrsStream) (bool, error) {
	return true, nil
}
