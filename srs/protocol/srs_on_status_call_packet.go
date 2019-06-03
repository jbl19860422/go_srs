package protocol

type SrsOnStatusCallPacket struct {
	command_name   string
	Transaction_id float64
	Args           *SrsAmf0Object
	Data           *SrsAmf0Object
}

func NewSrsOnStatusCallPacket() *SrsOnStatusCallPacket {
	return &SrsOnStatusCallPacket{
		command_name:   RTMP_AMF0_COMMAND_ON_STATUS,
		Transaction_id: 0,
		Args:           NewSrsAmf0Object(),
		Data:           NewSrsAmf0Object(),
	}
}

func (s *SrsOnStatusCallPacket) get_message_type() int8 {
	return RTMP_MSG_AMF0CommandMessage
}

func (s *SrsOnStatusCallPacket) get_prefer_cid() int32 {
	return RTMP_CID_OverStream
}

func (p *SrsOnStatusCallPacket) decode(stream *SrsStream) error {
	return nil
}

func (this *SrsOnStatusCallPacket) encode() ([]byte, error) {
	stream := NewSrsStream([]byte{}, 0)

	err := srs_amf0_write_string(stream, this.command_name)
	if err != nil {
		return nil, err
	}

	err = srs_amf0_write_number(stream, this.Transaction_id)
	if err != nil {
		return nil, err
	}

	err = srs_amf0_write_null(stream)
	if err != nil {
		return nil, err
	}

	err = this.Data.write(stream)
	if err != nil {
		return nil, err
	}

	return stream.p, nil
}
