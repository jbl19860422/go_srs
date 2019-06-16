package packet

import(
	"go_srs/srs/utils"
	"go_srs/srs/protocol/amf0"
	"go_srs/srs/global"
)
type SrsCreateStreamPacket struct {
	CommandName   amf0.SrsAmf0String
	TransactionId amf0.SrsAmf0Number
	CommandObj    *amf0.SrsAmf0Object
	NullObj		  amf0.SrsAmf0Null
}

func NewSrsCreateStreamPacket() *SrsCreateStreamPacket {
	return &SrsCreateStreamPacket{
		CommandName: amf0.SrsAmf0String{Value:amf0.SrsAmf0Utf8{Value:amf0.RTMP_AMF0_COMMAND_CREATE_STREAM}},
		CommandObj:  amf0.NewSrsAmf0Object(),
	}
}

func (s *SrsCreateStreamPacket) GetMessageType() int8 {
	return global.RTMP_MSG_AMF0CommandMessage
}

func (s *SrsCreateStreamPacket) GetPreferCid() int32 {
	return global.RTMP_CID_OverConnection
}

func (this *SrsCreateStreamPacket) Decode(stream *utils.SrsStream) error {
	err := this.TransactionId.Decode(stream)
	if err != nil {
		return err
	}

	err = this.NullObj.Decode(stream)
	if err != nil {
		return err
	}

	return nil
}

func (this *SrsCreateStreamPacket) Encode(stream *utils.SrsStream) error {
	_ = this.CommandName.Encode(stream)
	_ = this.TransactionId.Encode(stream)
	_ = this.NullObj.Encode(stream)
	return nil
}
