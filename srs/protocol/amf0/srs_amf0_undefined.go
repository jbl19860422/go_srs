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

type SrsAmf0Undefined struct {
}

func (this *SrsAmf0Undefined) Decode(stream *utils.SrsStream) error {
	marker, err := stream.ReadByte()
	if err != nil {
		return err
	}

	if marker != RTMP_AMF0_Undefined {
		err := errors.New("amf0 check null marker failed.")
		return err
	}
	return nil
}

func (this *SrsAmf0Undefined) Encode(stream *utils.SrsStream) error {
	stream.WriteByte(RTMP_AMF0_Undefined)
	return nil
}

func (this *SrsAmf0Undefined) IsMyType(stream *utils.SrsStream) (bool, error) {
	marker, err := stream.PeekByte()
	if err != nil {
		return false, err
	}

	if marker != RTMP_AMF0_Undefined {
		return false, nil
	}
	return true, nil
}

func (this *SrsAmf0Undefined) GetValue() interface{} {
	return nil
}
