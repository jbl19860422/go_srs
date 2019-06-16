package packet

import (
	_ "errors"
	_ "log"
	"go_srs/srs/protocol/amf0"
	"go_srs/srs/utils"
	"go_srs/srs/global"
)

type SrsFMLEStartPacket struct {
	CommandName		amf0.SrsAmf0String
	TransactionId 	amf0.SrsAmf0Number
	StreamName    	amf0.SrsAmf0String
	NullObj			amf0.SrsAmf0Null
}

func NewSrsFMLEStartPacket(name string) *SrsFMLEStartPacket {
	return &SrsFMLEStartPacket{
		CommandName: amf0.SrsAmf0String{Value:amf0.SrsAmf0Utf8{Value:name}},
	}
}

func (s *SrsFMLEStartPacket) GetMessageType() int8 {
	return global.RTMP_MSG_AMF0CommandMessage
}

func (s *SrsFMLEStartPacket) GetPreferCid() int32 {
	return global.RTMP_CID_OverConnection
}

func (this *SrsFMLEStartPacket) Decode(stream *utils.SrsStream) error {
	if err := this.TransactionId.Decode(stream); err != nil {
		return err
	}

	if err := this.NullObj.Decode(stream); err != nil {
		return err
	}

	if err := this.StreamName.Decode(stream); err != nil {
		return err
	}
	return nil
}

func (this *SrsFMLEStartPacket) Encode(stream *utils.SrsStream) error {
	_ = this.CommandName.Encode(stream)
	_ = this.TransactionId.Encode(stream)
	_ = this.NullObj.Encode(stream)
	_ = this.StreamName.Encode(stream)
	return nil
}
