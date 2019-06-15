package packet

type SrsCreateStreamPacket struct {
	CommandName   SrsAmf0String
	TransactionId SrsAmf0Number
	CommandObj    *SrsAmf0Object
	NullObj		  SrsAmf0Null
}

func NewSrsCreateStreamPacket() *SrsCreateStreamPacket {
	return &SrsCreateStreamPacket{
		CommandName: SrsAmf0String{Value:RTMP_AMF0_COMMAND_CREATE_STREAM},
		CommandObj:   NewSrsAmf0Object(),
	}
}

func (s *SrsCreateStreamPacket) GetMessageType() int8 {
	return RTMP_MSG_AMF0CommandMessage
}

func (s *SrsCreateStreamPacket) GetPreferCid() int32 {
	return RTMP_CID_OverConnection
}

func (this *SrsCreateStreamPacket) Decode(stream *SrsStream) error {
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

func (s *SrsCreateStreamPacket) Encode(stream *SrsStream) error {
	_ = this.CommandName.Encode(stream)
	_ = this.TransactionId.Encode(stream)
	_ = this.NullObj.Encode(stream)
	return nil
}
