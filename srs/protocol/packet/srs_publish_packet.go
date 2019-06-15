package packet

import (
	"errors"
)

type SrsPublishPacket struct {
	CommandName   	SrsAmf0String
	TransactionId 	SrsAmf0Number
	NullObj     	SrsAmf0Null
	StreamName    	SrsAmf0String
	Type            SrsAmf0String
}

func NewSrsPublishPacket() *SrsPublishPacket {
	return &SrsPublishPacket{
		CommandName:	SrsAmf0String{Value:RTMP_AMF0_COMMAND_PUBLISH},
		TransactionId: 	SrsAmf0Number{Value:0},
		CommandObj:     NewSrsAmf0Object(),
		Type:           SrsAmf0String{Value:"live"},
	}
}

func (s *SrsPublishPacket) GetMessageType() int8 {
	return RTMP_MSG_AMF0CommandMessage
}

func (s *SrsPublishPacket) GetPreferCid() int32 {
	return RTMP_CID_OverStream
}

func (this *SrsPublishPacket) Decode(stream *SrsStream) error {
	if err := this.TransactionId.Decode(stream); err != nil {
		return err
	}

	if err := this.NullObj.Decode(stream); err != nil {
		return err
	}
	
	if err := this.StreamName.Decode(stream); err != nil {
		return err
	}

	if !stream.empty() {
		if err := this.Type.Deocde(stream); err != nil {
			return err
		}
	}
	return nil
}

func (s *SrsPublishPacket) Encode(stream *SrsStream) error {
	_ = this.CommandName.Encode(stream)
	_ = this.TransactionId.Encode(stream)
	_ = this.NullObj.Encode(stream)
	_ = this.StreamName.Encode(stream)
	_ = this.Type.Encode(stream)
	return nil
}
