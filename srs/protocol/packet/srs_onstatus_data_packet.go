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
package packet

import (
	"go_srs/srs/global"
	"go_srs/srs/protocol/amf0"
	"go_srs/srs/utils"
)

type SrsOnStatusDataPacket struct {
	CommandName amf0.SrsAmf0String
	Data        *amf0.SrsAmf0Object
}

func NewSrsOnStatusDataPacket() *SrsOnStatusDataPacket {
	return &SrsOnStatusDataPacket{
		CommandName: amf0.SrsAmf0String{Value: amf0.SrsAmf0Utf8{Value: amf0.RTMP_AMF0_COMMAND_ON_STATUS}},
		Data:        amf0.NewSrsAmf0Object(),
	}
}

func (this *SrsOnStatusDataPacket) GetMessageType() int8 {
	return global.RTMP_MSG_AMF0DataMessage
}

func (this *SrsOnStatusDataPacket) GetPreferCid() int32 {
	return global.RTMP_CID_OverStream
}

func (p *SrsOnStatusDataPacket) Decode(stream *utils.SrsStream) error {
	return nil
}

func (this *SrsOnStatusDataPacket) Encode(stream *utils.SrsStream) error {
	_ = this.CommandName.Encode(stream)
	_ = this.Data.Encode(stream)
	return nil
}
