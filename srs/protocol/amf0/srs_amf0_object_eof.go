/*
The MIT License (MIT)

Copyright (c) 2019 GOSRS(gosrs)

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
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

func (this *SrsAmf0ObjectEOF) GetValue() interface{} {
	return nil
}
