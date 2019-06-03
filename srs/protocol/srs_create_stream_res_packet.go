package protocol

type SrsCreateStreamResPacket struct {
	command_name   string
	transaction_id float64
	CommandObj     *SrsAmf0Object
	Stream_id      float64
}

func NewSrsCreateStreamResPacket(tid float64, sid float64) *SrsCreateStreamResPacket {
	return &SrsCreateStreamResPacket{
		command_name:   RTMP_AMF0_COMMAND_RESULT,
		transaction_id: tid,
		CommandObj:     NewSrsAmf0Object(),
		Stream_id:      sid,
	}
}
func (s *SrsCreateStreamResPacket) get_message_type() int8 {
	return RTMP_MSG_AMF0CommandMessage
}

func (s *SrsCreateStreamResPacket) get_prefer_cid() int32 {
	return RTMP_CID_OverConnection
}

func (this *SrsCreateStreamResPacket) decode(stream *SrsStream) error {
	var err error
	this.transaction_id, err = srs_amf0_read_number(stream)
	if err != nil {
		return err
	}

	err = srs_amf0_read_null(stream)
	if err != nil {
		return err
	}

	this.Stream_id, err = srs_amf0_read_number(stream)
	if err != nil {
		return err
	}

	return nil
}

func (s *SrsCreateStreamResPacket) encode() ([]byte, error) {
	stream := NewSrsStream([]byte{}, 0)
	err := srs_amf0_write_string(stream, s.command_name)
	if err != nil {
		return nil, err
	}

	err = srs_amf0_write_number(stream, s.transaction_id)
	if err != nil {
		return nil, err
	}

	err = srs_amf0_write_null(stream)
	if err != nil {
		return nil, err
	}

	err = srs_amf0_write_number(stream, s.Stream_id)
	if err != nil {
		return nil, err
	}

	return stream.p, nil
}
