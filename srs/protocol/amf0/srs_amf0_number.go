package amf0

import (
	"go_srs/srs/utils"
	"encoding/binary"
	"errors"
)

type SrsAmf0Number struct {
	Value float64
}

func NewSrsAmf0Number(data float64) *SrsAmf0Number {
	return &SrsAmf0Number{
		Value: data,
	}
}

func (this *SrsAmf0Number) Decode(stream *utils.SrsStream) error {
	marker, err := stream.ReadByte()
	if err != nil {
		return err
	}

	if marker != RTMP_AMF0_Number {
		err := errors.New("amf0 check string marker failed.")
		return err
	}

	this.Value, err = stream.ReadFloat64(binary.BigEndian)
	if err != nil {
		return err
	}
	return nil
}

func (this *SrsAmf0Number) Encode(stream *utils.SrsStream) error {
	stream.WriteByte(RTMP_AMF0_Number)
	stream.WriteFloat64(this.Value, binary.BigEndian)
	return nil
}

func (this *SrsAmf0Number) IsMyType(stream *utils.SrsStream) (bool, error) {
	marker, err := stream.PeekByte()
	if err != nil {
		return false, err
	}

	if marker != RTMP_AMF0_Number {
		return false, nil
	}
	return true, nil
}

func (this *SrsAmf0Number) GetValue() interface{} {
	return this.Value
}