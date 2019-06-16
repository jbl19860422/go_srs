package packet

import (
	"go_srs/srs/protocol/amf0"
	"go_srs/srs/utils"
	"go_srs/srs/global"
)

type SrsOnStatusCallPacket struct {
	CommandName   amf0.SrsAmf0String
	TransactionId amf0.SrsAmf0Number
	NullObj       amf0.SrsAmf0Object
	Data          *amf0.SrsAmf0Object
}

func NewSrsOnStatusCallPacket() *SrsOnStatusCallPacket {
	return &SrsOnStatusCallPacket{
		CommandName:   amf0.SrsAmf0String{Value: amf0.SrsAmf0Utf8{Value:amf0.RTMP_AMF0_COMMAND_ON_STATUS}},
		TransactionId: amf0.SrsAmf0Number{Value: 0},
		Data:          amf0.NewSrsAmf0Object(),
	}
}

func (this *SrsOnStatusCallPacket) GetMessageType() int8 {
	return global.RTMP_MSG_AMF0CommandMessage
}

func (this *SrsOnStatusCallPacket) GetPreferCid() int32 {
	return global.RTMP_CID_OverStream
}

func (p *SrsOnStatusCallPacket) Decode(stream *utils.SrsStream) error {
	return nil
}

func (this *SrsOnStatusCallPacket) Encode(stream *utils.SrsStream) error {
	_ = this.CommandName.Encode(stream)
	_ = this.TransactionId.Encode(stream)
	_ = this.NullObj.Encode(stream)
	_ = this.Data.Encode(stream)
	return nil
}
