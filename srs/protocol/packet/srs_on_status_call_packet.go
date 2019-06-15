package packet

type SrsOnStatusCallPacket struct {
	CommandName   	SrsAmf0String
	TransactionId 	SrsAmf0Number
	NullObj         SrsAmf0Object
	Data            *SrsAmf0Object
}

func NewSrsOnStatusCallPacket() *SrsOnStatusCallPacket {
	return &SrsOnStatusCallPacket{
		CommandName:	SrsAmf0String{Value:RTMP_AMF0_COMMAND_ON_STATUS},
		TransactionId: 	SrsAmf0Number{Value:0},
		Data:           NewSrsAmf0Object(),
	}
}

func (s *SrsOnStatusCallPacket) GetMessageType() int8 {
	return RTMP_MSG_AMF0CommandMessage
}

func (s *SrsOnStatusCallPacket) GetPreferCid() int32 {
	return RTMP_CID_OverStream
}

func (p *SrsOnStatusCallPacket) Decode(stream *SrsStream) error {
	return nil
}

func (this *SrsOnStatusCallPacket) Encode(stream *SrsStream) error {
	_ = this.CommandName.Encode(stream)
	_ = this.TransactionId.Encode(stream)
	_ = this.NullObj.Encode(stream)
	_ = this.Data.Encode(stream)
	return nil
}
