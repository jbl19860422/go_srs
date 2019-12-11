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

func (this *SrsAmf0Utf8) GetValue() interface{} {
	return this.Value
}
