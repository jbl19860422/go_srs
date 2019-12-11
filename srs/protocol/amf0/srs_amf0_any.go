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
)

type SrsAmf0Any interface {
	Decode(stream *utils.SrsStream) error
	Encode(stream *utils.SrsStream) error
	IsMyType(stream *utils.SrsStream) (bool, error)
	GetValue() interface{}
}

func GenerateSrsAmf0Any(marker byte) SrsAmf0Any {
	switch marker {
	case RTMP_AMF0_Number:
		return &SrsAmf0Number{}
	case RTMP_AMF0_Boolean:
		return &SrsAmf0Boolean{}
	case RTMP_AMF0_String:
		return &SrsAmf0String{}
	case RTMP_AMF0_Object:
		return &SrsAmf0Object{}
	case RTMP_AMF0_Null:
		return &SrsAmf0Null{}
	case RTMP_AMF0_Undefined:
		return &SrsAmf0Undefined{}
	case RTMP_AMF0_EcmaArray:
		return &SrsAmf0EcmaArray{}
	default:
		return nil
	}
}
