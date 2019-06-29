package packet

import(
	"go_srs/srs/utils"
	"go_srs/srs/protocol/amf0"
	"go_srs/srs/global"
)

type SrsOnBwDonePacket struct {
	CommandName   	amf0.SrsAmf0String
	TransactionId 	amf0.SrsAmf0Number
	NullObj			amf0.SrsAmf0Null
}

func NewSrsOnBwDonePacket() *SrsOnBwDonePacket {
	return &SrsOnBwDonePacket{
		CommandName:   	amf0.SrsAmf0String{Value:amf0.SrsAmf0Utf8{Value:amf0.RTMP_AMF0_COMMAND_ON_BW_DONE}},
		TransactionId: 	amf0.SrsAmf0Number{Value:0},
	}
}

func (this *SrsOnBwDonePacket) GetMessageType() int8 {
	return global.RTMP_MSG_AMF0CommandMessage
}

func (this *SrsOnBwDonePacket) GetPreferCid() int32 {
	return global.RTMP_CID_OverConnection
}

func (this *SrsOnBwDonePacket) Decode(stream *utils.SrsStream) error {
	if err := this.TransactionId.Decode(stream); err != nil {
		return err
	}

	if err := this.NullObj.Decode(stream); err != nil {
		return err
	}

	return nil
}

func (this *SrsOnBwDonePacket) Encode(stream *utils.SrsStream) error {
	_ = this.CommandName.Encode(stream)
	_ = this.TransactionId.Encode(stream)
	_ = this.NullObj.Encode(stream)
	return nil
}