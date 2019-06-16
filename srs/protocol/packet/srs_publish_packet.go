package packet

import (
	"go_srs/srs/global"
	"go_srs/srs/protocol/amf0"
	"go_srs/srs/utils"
)

type SrsPublishPacket struct {
	CommandName   amf0.SrsAmf0String
	TransactionId amf0.SrsAmf0Number
	NullObj       amf0.SrsAmf0Null
	StreamName    amf0.SrsAmf0String
	Type          amf0.SrsAmf0String
}

func NewSrsPublishPacket() *SrsPublishPacket {
	return &SrsPublishPacket{
		CommandName:   amf0.SrsAmf0String{Value: amf0.SrsAmf0Utf8{Value: amf0.RTMP_AMF0_COMMAND_PUBLISH}},
		TransactionId: amf0.SrsAmf0Number{Value: 0},
		Type:          amf0.SrsAmf0String{Value: amf0.SrsAmf0Utf8{Value: "live"}},
	}
}

func (s *SrsPublishPacket) GetMessageType() int8 {
	return global.RTMP_MSG_AMF0CommandMessage
}

func (s *SrsPublishPacket) GetPreferCid() int32 {
	return global.RTMP_CID_OverStream
}

func (this *SrsPublishPacket) Decode(stream *utils.SrsStream) error {
	if err := this.TransactionId.Decode(stream); err != nil {
		return err
	}

	if err := this.NullObj.Decode(stream); err != nil {
		return err
	}

	if err := this.StreamName.Decode(stream); err != nil {
		return err
	}

	if !stream.Empty() {
		if err := this.Type.Decode(stream); err != nil {
			return err
		}
	}
	return nil
}

func (this *SrsPublishPacket) Encode(stream *utils.SrsStream) error {
	_ = this.CommandName.Encode(stream)
	_ = this.TransactionId.Encode(stream)
	_ = this.NullObj.Encode(stream)
	_ = this.StreamName.Encode(stream)
	_ = this.Type.Encode(stream)
	return nil
}
