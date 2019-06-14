package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
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
		p:       data,
		bytes:   data,
		pos:     0,
	}
}

func (this *SrsStream) Data() []byte {
	return this.bytes
}

func (this *SrsStream) Size() uint32 {
	return len(this.bytes)
}

func (this *SrsStream) Empty() bool {
	return this.bytes == nil || len(this.p) <= 0
}

func (this *SrsStream) Require(required_size uint32) bool {
	return required_size <= len(this.p)
}

func (s *SrsStream) Skip(size uint32) {
	s.pos += size
	s.p = s.bytes[s.pos:]
}

func (this *SrsStream) ReadByte() (byte, error) {
	if !this.Require(1) {
		err := errors.New("SrsStream not have enough data")
		return nil, err
	}

	b = this.p[0]
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

	b := this.p[0:n]
	this.Skip(count)
	return b, nil
}

func (this *SrsStream) WriteBytes(data []byte) {
	this.p = append(this.p, data...)
}

func (this *SrsStream) ReadInt16(order binary.ByteOrder) (int16, error) {
	b, err := this.ReadBytes(2)
	if err != nil {
		return err
	}

	v, err := utils.BytesToInt16(b, order)
	return v, err
}

func (this *SrsStream) WriteInt16(data int16, order binary.ByteOrder) {
	b := utils.Int16ToBytes(data, order)
	this.WriteBytes(b)
}

func (this *SrsStream) ReadInt32(order binary.ByteOrder) (int16, error) {
	b, err := this.ReadBytes(4)
	if err != nil {
		return err
	}

	v, err := utils.BytesToInt32(b, order)
	return v, err
}

func (this *SrsStream) WriteInt32(data int32, order binary.ByteOrder) {
	b := utils.Int32ToBytes(data, order)
	this.WriteBytes(b)
}




func (s *SrsStream) read_nbytes(n int32) (b []byte, err error) {
	if !s.require(n) {
		err = errors.New("no enough data")
		return
	}

	b = s.p[0 : n+1]
	s.skip(n)
	return
}

func (s *SrsStream) read_int8() (v int8, err error) {
	b, err := s.read_nbytes(1)
	if err != nil {
		return
	}
	bin_buf := bytes.NewBuffer(b)
	binary.Read(bin_buf, binary.BigEndian, &v)
	return
}

func (s *SrsStream) read_bool() (v bool, err error) {
	b, err := s.read_nbytes(1)
	if err != nil {
		return
	}
	bin_buf := bytes.NewBuffer(b)
	binary.Read(bin_buf, binary.BigEndian, &v)
	return
}

func (s *SrsStream) read_int16() (v int16, err error) {
	b, err := s.read_nbytes(2)
	if err != nil {
		return
	}

	bin_buf := bytes.NewBuffer(b)
	log.Printf("read int16 %x %x", b[0], b[1])
	binary.Read(bin_buf, binary.BigEndian, &v)
	return
}

func (s *SrsStream) read_int32() (v int32, err error) {
	b, err := s.read_nbytes(4)
	if err != nil {
		return
	}

	bin_buf := bytes.NewBuffer(b)
	binary.Read(bin_buf, binary.BigEndian, &v)
	return
}

func (s *SrsStream) write_int32(v int32) error {
	b := IntToBytes(int(v))
	d := make([]byte, 4)
	d[0] = b[3]
	d[1] = b[2]
	d[2] = b[1]
	d[3] = b[0]
	s.write_bytes(d)
	return nil
}

func (s *SrsStream) read_int64() (v int64, err error) {
	b, err := s.read_nbytes(8)
	if err != nil {
		return
	}

	bin_buf := bytes.NewBuffer(b)
	binary.Read(bin_buf, binary.BigEndian, &v)
	return
}

func (s *SrsStream) read_float64() (v float64, err error) {
	b, err := s.read_nbytes(8)
	if err != nil {
		return
	}

	bin_buf := bytes.NewBuffer(b)
	binary.Read(bin_buf, binary.BigEndian, &v)
	return
}

func (s *SrsStream) write_float64(v float64) (err error) {
	bin_buf := Float64ToByte(v)
	s.write_bytes(bin_buf)
	return nil
}

func (s *SrsStream) read_string(len int32) (str string, err error) {
	if !s.require(len) {
		err = errors.New("no enough data")
		return
	}

	str = string(s.p[:len])
	s.skip(len)
	err = nil
	return
}

func (s *SrsStream) write_1byte(b byte) {
	s.p = append(s.p, b)
}
func (s *SrsStream) write_bytes(d []byte) {
	s.p = append(s.p, d...)
}

func (s *SrsStream) write_string(v string) {
	s.p = append(s.p, []byte(v)...)
}
