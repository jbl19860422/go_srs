package packet

import "errors"

type SrsFMLEStartResPacket struct {
	CommandName   	SrsAmf0String
	TransactionId 	SrsAmf0Number
	NullObj			SrsAmf0Null
	UndefinedObj	SrsAmf0Undefined
}

func NewSrsFMLEStartResPacket(trans_id float64) *SrsFMLEStartResPacket {
	return &SrsFMLEStartResPacket{
		CommandName:   	SrsAmf0String{Value:RTMP_AMF0_COMMAND_RESULT},
		TransactionId: 	SrsAmf0Number{Value:trans_id}
	}
}

func (s *SrsFMLEStartResPacket) GetMessageType() int8 {
	return RTMP_MSG_AMF0CommandMessage
}

func (s *SrsFMLEStartResPacket) GetPreferCid() int32 {
	return RTMP_CID_OverConnection
}

func (p *SrsFMLEStartResPacket) Decode(stream *SrsStream) error {
	if err := this.TransactionId.Decode(stream); err != nil {
		return err
	}

	if err := this.NullObj.Deocde(stream); err != nil {
		return err
	}

	if err := this.UndefinedObj.Decode(stream); err != nil {
		return err
	}
	return nil
}

func (this *SrsFMLEStartResPacket) Encode(stream *SrsStream) error {
	_ = this.CommandName.Encode(stream)
	_ = this.TransactionId.Encode(stream)
	_ = this.NullObj.Encode(stream)
	_ = this.UndefinedObj.Encode(stream)
	return nil
}
