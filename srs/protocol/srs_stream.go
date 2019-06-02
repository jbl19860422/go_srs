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
	// the total number of bytes.
	n_bytes int32
	// current position
	pos int32
}

func NewSrsStream(data []byte, len int32) *SrsStream {
	return &SrsStream{
		p:       data,
		bytes:   data,
		n_bytes: len,
		pos:     0,
	}
}

func (s *SrsStream) data() []byte {
	return s.bytes
}

func (s *SrsStream) size() int32 {
	return s.n_bytes
}

func (s *SrsStream) empty() bool {
	return s.bytes == nil || len(s.p) <= 0
}

func (s *SrsStream) require(required_size int32) bool {
	return int(required_size) <= len(s.p)
}

func (s *SrsStream) skip(size int32) {
	s.pos += size
	s.p = s.bytes[s.pos:]
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
	log.Print("...................len=", len(bin_buf), "...................")
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
