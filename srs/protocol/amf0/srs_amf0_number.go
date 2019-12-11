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