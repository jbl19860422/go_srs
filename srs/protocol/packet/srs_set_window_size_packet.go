package packet

import (
	"encoding/binary"
	"go_srs/srs/global"
	"go_srs/srs/utils"
)
type SrsSetWindowAckSizePacket struct {
	AckowledgementWindowSize int32
}

func NewSrsSetWindowAckSizePacket() *SrsSetWindowAckSizePacket {
	return &SrsSetWindowAckSizePacket{}
}

func (this *SrsSetWindowAckSizePacket) GetMessageType() int8 {
	return global.RTMP_MSG_WindowAcknowledgementSize
}

func (this *SrsSetWindowAckSizePacket) GetPreferCid() int32 {
	return global.RTMP_CID_ProtocolControl
}

func (this *SrsSetWindowAckSizePacket) Decode(stream *utils.SrsStream) error {
	var err error
	this.AckowledgementWindowSize, err = stream.ReadInt32(binary.LittleEndian)
	return err
}

func (this *SrsSetWindowAckSizePacket) Encode(stream *utils.SrsStream) error {
	stream.WriteInt32(this.AckowledgementWindowSize, binary.LittleEndian)
	return nil
}
