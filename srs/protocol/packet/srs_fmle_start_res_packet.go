package packet

import(
	"go_srs/srs/utils"
	"go_srs/srs/protocol/amf0"
	"go_srs/srs/global"
)

type SrsFMLEStartResPacket struct {
	CommandName   	amf0.SrsAmf0String
	TransactionId 	amf0.SrsAmf0Number
	NullObj			amf0.SrsAmf0Null
	UndefinedObj	amf0.SrsAmf0Undefined
}

func NewSrsFMLEStartResPacket(trans_id float64) *SrsFMLEStartResPacket {
	return &SrsFMLEStartResPacket{
		CommandName:   	amf0.SrsAmf0String{Value:amf0.SrsAmf0Utf8{Value:amf0.RTMP_AMF0_COMMAND_RESULT}},
		TransactionId: 	amf0.SrsAmf0Number{Value:trans_id},
	}
}

func (this *SrsFMLEStartResPacket) GetMessageType() int8 {
	return global.RTMP_MSG_AMF0CommandMessage
}

func (this *SrsFMLEStartResPacket) GetPreferCid() int32 {
	return global.RTMP_CID_OverConnection
}

func (this *SrsFMLEStartResPacket) Decode(stream *utils.SrsStream) error {
	if err := this.TransactionId.Decode(stream); err != nil {
		return err
	}

	if err := this.NullObj.Decode(stream); err != nil {
		return err
	}

	if err := this.UndefinedObj.Decode(stream); err != nil {
		return err
	}
	return nil
}

func (this *SrsFMLEStartResPacket) Encode(stream *utils.SrsStream) error {
	_ = this.CommandName.Encode(stream)
	_ = this.TransactionId.Encode(stream)
	_ = this.NullObj.Encode(stream)
	_ = this.UndefinedObj.Encode(stream)
	return nil
}
