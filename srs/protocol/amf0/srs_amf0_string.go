package amf0

import (
	"encoding/binary"
	"errors"
	"utils"
)

type SrsAmf0String struct {
	value string
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

func (this *SrsAmf0String) Encode(stream *utils.SrsStream) error {
	stream.WriteByte(RTMP_AMF0_String)
	stream.WriteUInt16(uint16(len(this.value)))
	stream.WriteString(this.value)
}
