package packet

type SrsCreateStreamResPacket struct {
	CommandName   	SrsAmf0String
	TransactionId 	SrsAmf0Number
	NullObj			SrsAmf0Null
	CommandObj     	*SrsAmf0Object
	StreamId      	SrsAmf0Number
}

func NewSrsCreateStreamResPacket(tid float64, sid float64) *SrsCreateStreamResPacket {
	return &SrsCreateStreamResPacket{
		CommandName:   	SrsAmf0String{Value:RTMP_AMF0_COMMAND_RESULT},
		TransactionId: 	SrsAmf0Number{Value:tid},
		CommandObj:     NewSrsAmf0Object(),
		StreamId:      	SrsAmf0Number{Value:sid},
	}
}
func (s *SrsCreateStreamResPacket) GetMessageType() int8 {
	return RTMP_MSG_AMF0CommandMessage
}

func (s *SrsCreateStreamResPacket) GetPreferCid() int32 {
	return RTMP_CID_OverConnection
}

func (this *SrsCreateStreamResPacket) Decode(stream *SrsStream) error {
	if err := this.TransactionId.Decode(stream); err != nil {
		return err
	}
	
	if err := this.NullObj.Decode(stream); err != nil {
		return err
	}

	if err := this.StreamId.Decode(stream); err != nil {
		return err
	}

	return nil
}

func (s *SrsCreateStreamResPacket) Encode(stream *SrsStream) error {
	_ = this.CommandName.Encode(stream)
	_ = this.TransactionId.Encode(stream)
	_ = this.NullObj.Encode(stream)
	_ = this.StreamId.Encode(stream)
	return nil
}
