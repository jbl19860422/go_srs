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
	"encoding/binary"
	"go_srs/srs/utils"
	"go_srs/srs/global"
)

type SrsPeerBandwidthType int

// 5.6. Set Peer Bandwidth (6)
const (
	_                       SrsPeerBandwidthType = iota
	SrsPeerBandwidthHard                         = 0x00
	SrsPeerBandwidthSoft                         = 0x01
	SrsPeerBandwidthDynamic                      = 0x02
)

type SrsSetPeerBandwidthPacket struct {
	Bandwidth int32
	Type      int8
}

func NewSrsSetPeerBandwidthPacket() *SrsSetPeerBandwidthPacket {
	return &SrsSetPeerBandwidthPacket{
		Bandwidth: 0,
		Type:      SrsPeerBandwidthDynamic,
	}
}

func (this *SrsSetPeerBandwidthPacket) GetMessageType() int8 {
	return global.RTMP_MSG_SetPeerBandwidth
}

func (this *SrsSetPeerBandwidthPacket) GetPreferCid() int32 {
	return global.RTMP_CID_ProtocolControl
}

func (this *SrsSetPeerBandwidthPacket) Decode(stream *utils.SrsStream) error {
	return nil
}

func (this *SrsSetPeerBandwidthPacket) Encode(stream *utils.SrsStream) error {
	stream.WriteInt32(this.Bandwidth, binary.BigEndian)
	stream.WriteByte(byte(this.Type))
	return nil
}
