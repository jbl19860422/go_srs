/*
The MIT License (MIT)

Copyright (c) 2013-2015 GOSRS(gosrs)

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

type SrsAmf0Boolean struct {
	Value bool
}

func NewSrsAmf0Boolean(data bool) *SrsAmf0Boolean {
	return &SrsAmf0Boolean{
		Value: data,
	}
}

func (this *SrsAmf0Boolean) Decode(stream *utils.SrsStream) error {
	marker, err := stream.ReadByte()
	if err != nil {
		return err
	}

	if marker != RTMP_AMF0_Boolean {
		err := errors.New("amf0 check bool marker failed.")
		return err
	}

	this.Value, err = stream.ReadBool()
	if err != nil {
		return err
	}
	return nil
}

func (this *SrsAmf0Boolean) Encode(stream *utils.SrsStream) error {
	stream.WriteByte(RTMP_AMF0_Boolean)
	var d byte
	if this.Value {
		d = 1
	} else {
		d = 0
	}
	stream.WriteByte(d)
	return nil
}

func (this *SrsAmf0Boolean) IsMyType(stream *utils.SrsStream) (bool, error) {
	marker, err := stream.PeekByte()
	if err != nil {
		return false, err
	}

	if marker != RTMP_AMF0_Boolean {
		return false, nil
	}
	return true, nil
}

func (this *SrsAmf0Boolean) GetValue() interface{} {
	return this.Value
}
