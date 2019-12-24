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
	"errors"
	"go_srs/srs/utils"
)

type SrsAmf0String struct {
	Value SrsAmf0Utf8
}

func NewSrsAmf0String(str string) *SrsAmf0String {
	return &SrsAmf0String{
		Value: SrsAmf0Utf8{Value: str},
	}
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

	return this.Value.Decode(stream)
}

func (this *SrsAmf0String) Encode(stream *utils.SrsStream) error {
	stream.WriteByte(RTMP_AMF0_String)
	this.Value.Encode(stream)
	return nil
}

func (this *SrsAmf0String) IsMyType(stream *utils.SrsStream) (bool, error) {
	marker, err := stream.PeekByte()
	if err != nil {
		return false, err
	}

	if marker == RTMP_AMF0_String {
		return true, nil
	}
	return false, nil
}

func (this *SrsAmf0String) GetValue() interface{} {
	return this.Value.Value
}
