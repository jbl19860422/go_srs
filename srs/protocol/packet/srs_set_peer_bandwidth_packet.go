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
	stream.WriteInt32(this.Bandwidth, binary.LittleEndian)
	stream.WriteByte(byte(this.Type))
	return nil
}
