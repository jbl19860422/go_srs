package packet

import (
	"errors"
	"log"
)

type SrsFMLEStartPacket struct {
	CommandName		SrsAmf0String
	TransactionId 	SrsAmf0Number
	StreamName    	SrsAmf0String
	NullObj			SrsAmf0Null
}

func NewSrsFMLEStartPacket(name string) *SrsFMLEStartPacket {
	return &SrsFMLEStartPacket{
		CommandName: SrsAmf0String{Value:name},
	}
}

func (s *SrsFMLEStartPacket) GetMessageType() int8 {
	return RTMP_MSG_AMF0CommandMessage
}

func (s *SrsFMLEStartPacket) GetPreferCid() int32 {
	return RTMP_CID_OverConnection
}

func (this *SrsFMLEStartPacket) Decode(stream *SrsStream) error {
	if err := this.TransactionId.Decode(stream); err != nil {
		return err
	}

	if err := this.NullObj.Decode(stream); err != nil {
		return err
	}

	if err != this.StreamName.Decode(stream); err != nil {
		return err
	}
	return nil
}

func (s *SrsFMLEStartPacket) Encode(stream *SrsStream) error {
	_ = this.CommandName.Encode(stream)
	_ = this.TransactionId.Encode(stream)
	_ = this.NullObj.Encode(stream)
	_ = this.StreamName.Encode(stream)
	return nil
}
