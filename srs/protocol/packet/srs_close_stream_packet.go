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

package packet
import (
	"go_srs/srs/protocol/amf0"
	"go_srs/srs/utils"
	"go_srs/srs/global"
)

type SrsCloseStreamPacket struct {
	CommandName		amf0.SrsAmf0String
	TransactionId	amf0.SrsAmf0Number
	NullObj			*amf0.SrsAmf0Object
}

func NewSrsCloseStreamPacket() *SrsCloseStreamPacket {
	return &SrsCloseStreamPacket{
		CommandName: amf0.SrsAmf0String{Value:amf0.SrsAmf0Utf8{Value:amf0.RTMP_AMF0_COMMAND_CLOSE_STREAM}},
		TransactionId: amf0.SrsAmf0Number{Value:0},
	}
}

func (s *SrsCloseStreamPacket) GetMessageType() int8 {
	return global.RTMP_MSG_AMF0CommandMessage
}

func (s *SrsCloseStreamPacket) GetPreferCid() int32 {
	return global.RTMP_CID_OverStream
}

func (this *SrsCloseStreamPacket) Decode(stream *utils.SrsStream) error {
	var err error
	err = this.TransactionId.Decode(stream)
	if err != nil {
		return err
	}

	err = this.NullObj.Decode(stream)
	if err != nil {
		return err
	}

	return nil
}

func (this *SrsCloseStreamPacket) Encode(stream *utils.SrsStream) error {
	return nil
}