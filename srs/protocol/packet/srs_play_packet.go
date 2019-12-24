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
	"go_srs/srs/global"
	"go_srs/srs/protocol/amf0"
	"go_srs/srs/utils"
)

type SrsPlayPacket struct {
	CommandName   amf0.SrsAmf0String
	TransactionId amf0.SrsAmf0Number
	NullObj       amf0.SrsAmf0Null
	StreamName    amf0.SrsAmf0String
	Start         amf0.SrsAmf0Number
	Duration      amf0.SrsAmf0Number
	Reset         amf0.SrsAmf0Boolean
}

func NewSrsPlayPacket() *SrsPlayPacket {
	return &SrsPlayPacket{
		CommandName:   amf0.SrsAmf0String{Value: amf0.SrsAmf0Utf8{Value: amf0.RTMP_AMF0_COMMAND_PLAY}},
		TransactionId: amf0.SrsAmf0Number{Value: 0},
		Start:         amf0.SrsAmf0Number{Value: -2.0},
		Duration:      amf0.SrsAmf0Number{Value: -1.0},
		Reset:         amf0.SrsAmf0Boolean{Value: true},
	}
}

func (s *SrsPlayPacket) GetMessageType() int8 {
	return global.RTMP_MSG_AMF0CommandMessage
}

func (s *SrsPlayPacket) GetPreferCid() int32 {
	return global.RTMP_CID_OverStream
}

func (this *SrsPlayPacket) Decode(stream *utils.SrsStream) error {
	if err := this.TransactionId.Decode(stream); err != nil {
		return err
	}

	if err := this.NullObj.Decode(stream); err != nil {
		return err
	}

	if err := this.StreamName.Decode(stream); err != nil {
		return err
	}

	if len(this.StreamName.GetValue().(string)) > 0 {
		if err := this.Start.Decode(stream); err != nil {
			return err
		}
	}
	//todo fix this
	return nil
}

func (this *SrsPlayPacket) Encode(stream *utils.SrsStream) error {
	_ = this.CommandName.Encode(stream)
	_ = this.TransactionId.Encode(stream)
	_ = this.NullObj.Encode(stream)
	_ = this.StreamName.Encode(stream)
	//why
	if this.Start.GetValue().(float64) != -2 || this.Duration.GetValue().(float64) != -1 || !this.Reset.GetValue().(bool) {
		err := this.Start.Encode(stream)
		if err != nil {
			return err
		}
	}
	//why
	if this.Duration.GetValue().(float64) != -1 || !this.Reset.GetValue().(bool) {
		_ = this.Duration.Encode(stream)
	}

	if !this.Reset.GetValue().(bool) {
		_ = this.Reset.Encode(stream)
	}

	return nil
}
