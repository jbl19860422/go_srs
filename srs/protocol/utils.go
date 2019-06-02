package protocol

import (
	"bytes"
	"encoding/binary"
	"math"
)

func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, x)
	return bytesBuffer.Bytes()
}

func Float64ToByte(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, bits)

	return bytes
}

type SrsAmf0Size struct {
}

func (s SrsAmf0Size) utf8(value string) int {
	return 2 + len(value)
}

func (s SrsAmf0Size) str(value string) int {
	return 1 + s.utf8(value)
}

func (s SrsAmf0Size) number() int {
	return 1 + 8
}

func (s SrsAmf0Size) date() int {
	return 1 + 8 + 2
}

func (s SrsAmf0Size) null() int {
	return 1
}

func (s SrsAmf0Size) undefined() int {
	return 1
}

func (s SrsAmf0Size) boolean() int {
	return 1 + 1
}

func (s SrsAmf0Size) object(obj *SrsAmf0Object) int {
	if obj == nil {
		return 0
	}

	return obj.total_size()
}

func (s SrsAmf0Size) object_eof() int {
	return 2 + 1
}

func (s SrsAmf0Size) any(v interface{}) int {
	var size int = 1
	switch v.(type) {
	case string:

	}
	return size
}
