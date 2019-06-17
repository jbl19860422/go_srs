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
