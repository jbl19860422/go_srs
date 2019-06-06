package utils

import (
	"encoding/binary"
	"errors"
)

type SrsStream struct {
	// current position at bytes.
	p []byte
	// the bytes data for stream to read or write.
	bytes []byte
	// current position
	pos uint32
}

func NewSrsStream(data []byte) *SrsStream {
	return &SrsStream{
		p:     data,
		bytes: data,
		pos:   0,
	}
}

func (this *SrsStream) Data() []byte {
	return this.bytes
}

func (this *SrsStream) Size() uint32 {
	return uint32(len(this.bytes))
}

func (this *SrsStream) Empty() bool {
	return this.bytes == nil || len(this.p) <= 0
}

func (this *SrsStream) Require(required_size uint32) bool {
	return required_size <= uint32(len(this.p))
}

func (s *SrsStream) Skip(size uint32) {
	s.pos += size
	s.p = s.bytes[s.pos:]
}

func (this *SrsStream) ReadByte() (byte, error) {
	if !this.Require(1) {
		err := errors.New("SrsStream not have enough data")
		return 0, err
	}

	b := this.p[0]
	this.Skip(1)
	return b, nil
}

func (this *SrsStream) WriteByte(data byte) {
	this.p = append(this.p, data)
}

func (this *SrsStream) ReadBytes(count uint32) ([]byte, error) {
	if !this.Require(count) {
		err := errors.New("SrsStream not have enough data")
		return nil, err
	}

	b := this.p[0:count]
	this.Skip(count)
	return b, nil
}

func (this *SrsStream) WriteBytes(data []byte) {
	this.p = append(this.p, data...)
}

func (this *SrsStream) ReadInt16(order binary.ByteOrder) (int16, error) {
	b, err := this.ReadBytes(2)
	if err != nil {
		return 0, err
	}

	v, err := BytesToInt16(b, order)
	return v, err
}

func (this *SrsStream) WriteInt16(data int16, order binary.ByteOrder) {
	b := Int16ToBytes(data, order)
	this.WriteBytes(b)
}

func (this *SrsStream) ReadInt32(order binary.ByteOrder) (int32, error) {
	b, err := this.ReadBytes(4)
	if err != nil {
		return 0, err
	}

	v, err := BytesToInt32(b, order)
	return v, err
}

func (this *SrsStream) WriteInt32(data int32, order binary.ByteOrder) {
	b := Int32ToBytes(data, order)
	this.WriteBytes(b)
}

func (this *SrsStream) ReadInt64(order binary.ByteOrder) (int64, error) {
	b, err := this.ReadBytes(8)
	if err != nil {
		return 0, err
	}

	v, err := BytesToInt64(b, order)
	return v, err
}

func (this *SrsStream) WriteInt64(data int64, order binary.ByteOrder) {
	b := Int64ToBytes(data, order)
	this.WriteBytes(b)
}

func (this *SrsStream) ReadFloat32(order binary.ByteOrder) (float32, error) {
	b, err := this.ReadBytes(4)
	if err != nil {
		return 0, err
	}

	v, err := BytesToFloat32(b, order)
	return v, err
}

func (this *SrsStream) WriteFloat32(data float32, order binary.ByteOrder) {
	b := Float32ToBytes(data, order)
	this.WriteBytes(b)
}

func (this *SrsStream) ReadFloat64(order binary.ByteOrder) (float64, error) {
	b, err := this.ReadBytes(8)
	if err != nil {
		return 0, err
	}

	v, err := BytesToFloat64(b, order)
	return v, err
}

func (this *SrsStream) WriteFloat64(data float64, order binary.ByteOrder) {
	b := Float64ToBytes(data, order)
	this.WriteBytes(b)
}

func (this *SrsStream) ReadString(len uint32) (string, error) {
	if !this.Require(len) {
		err := errors.New("no enough data")
		return "", err
	}

	str := string(this.p[:len])
	this.Skip(len)
	return str, nil
}

func (this *SrsStream) WriteString(str string) {
	this.WriteBytes([]byte(str))
}
